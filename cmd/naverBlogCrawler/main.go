package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
)

const (
	NUM_WORKERS = 1
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}
	// sources := seed.ListSources(domain.NaverBlog)

	sources := []domain.CrawlingSource{
		{SourceName: "mardukas",
			SourceID: "mardukas",
			Type:     domain.NaverBlog,
			Constraint: []string{
				"9", "61", "62", "111", "65", "66", "1", "68", "70", "67",
			},
			SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
			SpreadSheetName: "mardukas",
		},
		{SourceName: "paperchan",
			SourceID: "paperchan",
			Type:     domain.NaverBlog,
			Constraint: []string{
				"5", "2", "3", "4", "15",
			},
			SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
			SpreadSheetName: "paperchan",
		},
	}

	for _, channel := range sources {
		posts, err := crawling.FetchAllBlogPosts(&channel, 1)
		if err != nil {
			log.Println("Err", err)
			panic(err)
		}

		updating.SaveToSheetAppendNaverBlog(sheetsService, &channel, posts)
	}

}
