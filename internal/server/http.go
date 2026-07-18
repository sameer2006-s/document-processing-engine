package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
)

func RunHttpServer(authHandler *auth.AuthHandler, authService *auth.AuthService, documentHandler *document.DocumentHandler, cfg config.AppConfig) error {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/register", authHandler.RegisterUser)
	router.POST("/login", authHandler.LoginUser)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(authService))
	protected.POST("/upload", documentHandler.UploadFile)
	protected.GET("/get-file/:id", documentHandler.GetFile)

	err := router.Run(":" + cfg.Port)
	if err != nil {
		return err
	}
	log.Println("Server started successfully on port " + cfg.Port)
	return nil
}
