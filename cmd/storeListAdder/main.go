package main

import (
	"log"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/registering"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
)

// Used when crawls whole sheets of an youtube channel and add to list
func main() {
	ctx, cancel := registering.OpenChrome()
	registering.Login(ctx)

	sources := seed.ListSources(domain.Youtube)
	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}
	for _, channel := range sources {
		videos, err := updating.ListUnProcessedVideos(sheetsService, channel)
		if err != nil {
			log.Fatalf("Error fetching unprocessed videos: %v", err)
		}
		registering.AddStoreToList(ctx, channel, videos)
	}

	defer cancel()
}
