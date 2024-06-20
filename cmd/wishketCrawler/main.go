package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	sources := seed.ListSources(domain.Wishket)
	channel := sources[0]

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}

	lastestProjectURL := updating.GetLastProjectUrl(sheetsService, channel)

	projects := crawling.CrawlWishket()

	if len(projects) == 0 {
		log.Println("No projects found.")
		return
	}

	if projects[0].URL == lastestProjectURL {
		log.Println("No new projects found.")
		return
	}

	for _, project := range projects {
		if project.URL == lastestProjectURL {
			break
		}
		err = updating.SendWishketProjectToSlack(channel, project)
		if err != nil {
			log.Fatalf("Error sending to Slack: %v", err)
		}
	}

	updating.UpdateCheckpoint(projects[0].URL)
}
