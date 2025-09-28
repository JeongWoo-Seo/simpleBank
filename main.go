package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/JeongWoo-Seo/simpleBank/api"
	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	_ "github.com/JeongWoo-Seo/simpleBank/doc/statik"
	"github.com/JeongWoo-Seo/simpleBank/gapi"
	"github.com/JeongWoo-Seo/simpleBank/mail"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/JeongWoo-Seo/simpleBank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Enviroment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ cannot connect to db:")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGrop, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGrop, config, redisOpt, store)
	runGatewayServer(ctx, waitGrop, config, store, taskDistributor)
	runGrpcServer(ctx, waitGrop, config, store, taskDistributor)

	err = waitGrop.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("err from wait group")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migrate instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up:")
	}

	log.Info().Msg("db migrate successfully")
}

func runTaskProcessor(
	ctx context.Context,
	waitGrop *errgroup.Group,
	config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processer")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGrop.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task proccessor")

		taskProcessor.Shutdown()
		log.Info().Msg("task proccessor is stoped")

		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGrop *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("can not create server")
	}

	grpcMux := runtime.NewServeMux()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // 허용할 Origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // 쿠키/인증정보 허용 여부
	})
	handler := c.Handler(gapi.HttpLogger(mux))
	httpServer := &http.Server{
		Handler: handler,
		Addr:    config.HTTPServerAddress,
	}

	waitGrop.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("can not start HTTP gateway server")
			return err
		}
		return nil
	})

	waitGrop.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown http gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdwon http gateway server")
			return err
		}

		log.Info().Msg("http gateway server is stoped")
		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	waitGrop *errgroup.Group,
	config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("can not create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("can not create listener")
	}

	waitGrop.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Err(err).Msg("can not start gRPC server")
			return err
		}
		return nil
	})

	waitGrop.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown grpc server")

		grpcServer.GracefulStop()
		log.Info().Msg("grpc server stoped")

		return nil
	})
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
	}

	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("fail start server")
	}
}
