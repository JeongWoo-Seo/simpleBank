package api

import (
	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
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

	server.router.POST("/account", server.createAccount)
	server.router.GET("/account/:id", server.getAccount)
	server.router.GET("/account", server.listAccount)
	return server
}

func (s *Server) StartServer(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
