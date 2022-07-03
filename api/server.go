package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/zohaibAsif/simple_bank_management_system/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/account", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)

	router.POST("/transfer", server.createTransfer)
	router.GET("/transfer/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	router.GET("/entry/:id", server.getEntry)
	router.GET("/entries", server.listEntries)

	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
