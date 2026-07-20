package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sameer2006-s/document-processing-engine/internal/chat"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	myDB "github.com/sameer2006-s/document-processing-engine/internal/db"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/ocr"
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

	temporalClient, err := temporal.NewTemporalClient()
	if err != nil {
		log.Fatal("Failed to create temporal client: ", err)
	}
	defer temporalClient.Close()
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	minioClient, err := document.NewMinIOClient(cfg.Minio)
	if err != nil {
		log.Fatal("Failed to create minio client: ", err)
	}
	provider := chat.NewChatProvider(cfg.Chat.GitHubToken, cfg.Chat.Model)
	ocrService := ocr.NewOCRService(db, *minioClient)

	documentRepository, err := document.NewDocumentRepository(db, cfg.Minio)
	if err != nil {
		log.Fatal("Failed to create document repository: ", err)
	}
	documentService := document.NewDocumentService(documentRepository, nil)

	chatRepository := chat.NewChatRepository(db)
	chatService := chat.NewChatService(provider, documentService, chatRepository)
	
	documentActivity := temporal.NewDocumentActivity(ocrService, chatService, documentService)

	w := temporal.NewTemporalWorker(temporalClient, documentActivity)
	if err = w.RegisterWorkflows(); err != nil {
		log.Fatal("Failed to register workflows: ", err)
	}
	if err = w.RegisterActivities(); err != nil {
		log.Fatal("Failed to register activities: ", err)
	}
	if err = w.Run(); err != nil {
		log.Fatal("Failed to run worker: ", err)
	}
	log.Println("Worker run successfully.")
}
