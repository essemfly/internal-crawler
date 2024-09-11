package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/registering"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
)

// Used when crawls whole sheets of an youtube channel and add to list
func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ctx, cancel := pkg.OpenChrome()
	registering.NaverLogin(ctx)

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

func test() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ctx, cancel := pkg.OpenChrome()
	registering.NaverLogin(ctx)
	defer cancel()

	testSource := &domain.CrawlingSource{
		NaverListName: "테스트",
		NaverListID:   "1123eae3eca246a78d0b846683bf1a5c",
	}

	video := &domain.YoutubeVideoStruct{
		NaverLink: "https://map.naver.com/p/entry/place/37163963",
	}
	registering.AddStoreToList(ctx, testSource, []*domain.YoutubeVideoStruct{video})
	registering.AddStoreToList(ctx, testSource, []*domain.YoutubeVideoStruct{video})
}
