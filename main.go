package main

import (
	"context"
	"log"
	"os"

	"scraper/config"
	"scraper/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := db.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(cfg.Database.Name).Collection(cfg.Database.Collection)

	websites, err := LoadSource()
	if err != nil {
		log.Fatalf("Failed to load source: %v", err)
	}
	for _, website := range websites {
		log.Println(website)
		website.Scrape(collection, cfg)
	}
}
