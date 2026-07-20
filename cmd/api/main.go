package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	myDB "github.com/sameer2006-s/document-processing-engine/internal/db"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/chat"
	"github.com/sameer2006-s/document-processing-engine/internal/server"
	"github.com/sameer2006-s/document-processing-engine/internal/temporal"
)

func main() {
	_ = godotenv.Load() // Ignore error in production

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := myDB.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	err = myDB.MyAutoMigrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migrated successfully.")
	authRepository := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepository)
	authHandler := auth.NewAuthHandler(authService)

	documentRepository, err := document.NewDocumentRepository(db, cfg.Minio)
	if err != nil {
		log.Fatal("Failed to create document repository: ", err)
	}

	temporalClient, err := temporal.NewTemporalClient()
	if err != nil {
		log.Fatal("Failed to create temporal client: ", err)
	}
	defer temporalClient.Close()

	documentService := document.NewDocumentService(documentRepository, temporalClient.StartDocumentProcessing)
	documentHandler := document.NewDocumentHandler(documentService)

	chatProvider := chat.NewChatProvider(cfg.Chat.GitHubToken, cfg.Chat.Model)
	chatRepository := chat.NewChatRepository(db)
	chatService := chat.NewChatService(chatProvider, documentService, chatRepository)
	chatHandler := chat.NewChatHandler(chatService)

	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	err = server.RunHttpServer(authHandler, authService, documentHandler, chatHandler, cfg.App)
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
