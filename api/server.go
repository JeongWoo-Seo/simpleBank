package api

import (
	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router.POST("/account", server.createAccount)
	server.router.GET("/account/:id", server.getAccount)
	server.router.GET("/account", server.listAccount)

	server.router.POST("/transfers", server.createTransfer)

	return server
}

func (s *Server) StartServer(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
