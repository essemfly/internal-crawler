package main

import (
	"context"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	spreadSheet := "17kH954utV2XAvUJpEU59R2pFdSzz53j9odej6LcuzG0"
	sheetName1 := "MASONAR"
	sheetName2 := "감탄쟁이"

	channels := []domain.CrawlingSource{
		{
			SourceName:      "MASONAR",
			Type:            domain.Youtube,
			SourceID:        "PLCli9_-EPRzXkSedQWHH_gZlTki-bslYO",
			NaverListID:     "",
			SpreadSheetID:   spreadSheet,
			SpreadSheetName: sheetName1,
		},
		{
			SourceName:      "감탄쟁이",
			Type:            domain.Youtube,
			SourceID:        "PLkf5HhN3cnyFW0umotcYoWcjk0kCxVUWP",
			NaverListID:     "",
			SpreadSheetID:   spreadSheet,
			SpreadSheetName: sheetName2,
		},
	}

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	for _, channel := range channels {
		videos, err := crawling.GetChannelVideos(youtubeService, channel.SourceID, nil)
		if err != nil {
			log.Fatalf("Error fetching channel videos: %v", err)
		}

		log.Println("Starting YouTube Crawler named: " + channel.SourceName)
		log.Println("Number of videos: ", len(videos))
		videoSheets := domain.ConvertToYoutubeVideoStruct(videos)

		updating.SaveToSheetAtTop(sheetsService, &channel, videoSheets)
	}

}
