package main

import (
	"log"
	"math/rand"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"

	"go.uber.org/zap"
)

const (
	chunkSize        = 100
	numWorkers       = 5
	GlobalStartIndex = 783940000 // 2023-02-14 17:00:00
	// 788540500 : 2024-06-20 10:00:00
	// 783940000 : 2024-06-11 14:00:00
	// 533830000 : 2023-02-13 04:00:00
	// 533500000 : 2023-02-12 16:00:00
	// 531250000 : 2023-02-08 15:00:00
)

func DanggnCrawler() {

	sources := seed.ListSources(domain.Wishket)
	channel := sources[0]

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}

	workers := make(chan bool, numWorkers)
	done := make(chan bool, numWorkers)

	for c := 0; c < numWorkers; c++ {
		done <- true
	}

	lastIndex, err := updating.ReadLastIndex(sheetsService, channel)
	if err != nil {
		log.Fatalf("Error reading last index: %v", err)
	}

	log.Println("Last Index! ", lastIndex)
	for isIndexExists(lastIndex + chunkSize) {
		startIndex := lastIndex + 1
		lastIndex = startIndex + chunkSize - 1

		keywords, err := updating.ReadKeywords(sheetsService, channel)
		if err != nil {
			// config.Logger.Error("failed to find live keywords", zap.Error(err))
			log.Fatalln("failed to find live keywords", zap.Error(err))
			return
		}

		workers <- true
		<-done
		go func() {
			crawling.CrawlDanggnIndex(workers, done, keywords, startIndex, lastIndex)
		}()
	}

	for c := 0; c < numWorkers; c++ {
		<-done
	}
}

func isIndexExists(index int) bool {
	_, err := crawling.CrawlPage(index)
	errCounts := 0

	for err != nil {
		errCounts += 1
		if errCounts > 5 {
			return false
		}

		n := rand.Intn(11)
		_, err = crawling.CrawlPage(index + n)
	}

	return true
}
