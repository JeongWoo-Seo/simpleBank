package gapi

import (
	"fmt"

	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"github.com/JeongWoo-Seo/simpleBank/token"
	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/JeongWoo-Seo/simpleBank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
