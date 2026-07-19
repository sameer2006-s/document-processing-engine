package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	myDB "github.com/sameer2006-s/document-processing-engine/internal/db"
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

	temporalClient, err := temporal.NewTemporalClient("document-processing-workflow", "document-processing-task-queue")
	if err != nil {
		log.Fatal("Failed to create temporal client: ", err)
	}
	defer temporalClient.Client.Close()
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	worker := temporal.NewTemporalWorker(temporalClient)
	err = worker.RegisterWorkflows()
	if err != nil {
		log.Fatal("Failed to register workflows: ", err)
	}
	err = worker.RegisterActivities()
	if err != nil {
		log.Fatal("Failed to register activities: ", err)
	}
	err = worker.Run()
	if err != nil {
		log.Fatal("Failed to run worker: ", err)
	}
	log.Println("Worker run successfully.")
}