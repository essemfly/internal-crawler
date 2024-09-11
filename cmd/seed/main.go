package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	log.Println("Starting seed...")
	seed.CrawlingSeeds()
	log.Println("Seed complete.")
}
