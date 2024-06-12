package main

import (
	"context"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	sources := seed.ListSources(domain.Youtube)

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}

	// Youtube Video crawling
	for _, channel := range sources {
		latestVideo, err := updating.GetCurrentTopVideo(sheetsService, channel.SpreadSheetID, channel.SpreadSheetName)
		if err != nil {
			log.Fatalf("Error getting current top video: %v", err)
		}

		log.Println("Starting YouTube Crawler named: " + channel.SourceName)
		youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
		if err != nil {
			log.Fatalf("Error creating YouTube client: %v", err)
		}
		videos, err := crawling.GetChannelVideos(youtubeService, channel.SourceID, latestVideo)
		if err != nil {
			log.Fatalf("Error fetching channel videos: %v", err)
		}
		videoSheets := domain.ConvertToYoutubeVideoStruct(videos)
		videoSheets = domain.FilterWithChannelConstraints(videoSheets, channel)
		updating.SaveToSheetAtTop(sheetsService, channel, videoSheets)
		updating.SendVideosToSlack(channel, videoSheets)
	}

}
