package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	"github.com/sameer2006-s/document-processing-engine/internal/db"
)

func main() {
	_ = godotenv.Load() // Ignore error in production

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()
}