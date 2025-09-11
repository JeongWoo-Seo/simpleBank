package api

import (
	"fmt"

	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/token"
	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)

	authRouters := s.router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRouters.POST("/accounts", s.createAccount)
	authRouters.GET("/accounts/:id", s.getAccount)
	authRouters.GET("/accounts", s.listAccount)

	authRouters.POST("/transfers", s.createTransfer)
}

func (s *Server) StartServer(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
