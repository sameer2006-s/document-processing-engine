package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	myDB "github.com/sameer2006-s/document-processing-engine/internal/db"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"github.com/sameer2006-s/document-processing-engine/internal/server"
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
	authService := auth.NewAuthService(db)
	authHandler := auth.NewAuthHandler(authService)

	err = server.RunServer(authHandler)
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}
	

	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()
}