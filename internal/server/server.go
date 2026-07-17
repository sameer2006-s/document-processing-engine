package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
)

func RunServer(authHandler *auth.AuthHandler) error {
	router := gin.Default()
	router.POST("/register", authHandler.RegisterUser)
	router.POST("/login", authHandler.LoginUser)
	err := router.Run(":9000")
	if err != nil {
		return err
	}
	log.Println("Server started successfully on port 9000")
	return nil
}