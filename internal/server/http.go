package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/chat"
)

func RunHttpServer(authHandler *auth.AuthHandler, authService *auth.AuthService, documentHandler *document.DocumentHandler, chatHandler *chat.ChatHandler, cfg config.AppConfig) error {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/register", authHandler.RegisterUser)
	router.POST("/login", authHandler.LoginUser)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(authService))
	protected.POST("/upload", documentHandler.UploadFile)
	protected.GET("/documents", documentHandler.ListDocuments)
	protected.DELETE("/documents/:id", documentHandler.DeleteDocument)
	protected.GET("/get-file/:id", documentHandler.GetFile)
	protected.GET("/search-my-files", documentHandler.SearchMyFiles)
	protected.POST("/chat", chatHandler.Chat)
	protected.POST("/chat/stream", chatHandler.ChatStream)
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/favicon.svg", "./web/dist/favicon.svg")
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.File("./web/dist/index.html")
		}
	})

	err := router.Run(":" + cfg.Port)
	if err != nil {
		return err
	}
	log.Println("Server started successfully on port " + cfg.Port)
	return nil
}
