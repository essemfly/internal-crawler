package main

import (
	"log"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/registering"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
)

func main() {
	ctx, cancel := registering.OpenChrome()
	registering.Login(ctx)

	channels := seed.Seed()
	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}
	for _, channel := range channels {
		videos, err := updating.ListUnProcessedVideos(sheetsService, channel)
		if err != nil {
			log.Fatalf("Error fetching unprocessed videos: %v", err)
		}
		registering.AddStoreToList(ctx, channel, videos)
	}

	defer cancel()
}
