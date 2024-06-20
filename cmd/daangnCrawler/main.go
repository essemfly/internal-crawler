package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"

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

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	sources := seed.ListSources(domain.Daangn)
	channel := sources[0]
	log.Println("3Daangn Start")

	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		log.Fatalf("Error creating Sheets service: %v", err)
	}

	lastIndex, err := updating.ReadLastIndex(sheetsService, channel)
	if err != nil {
		log.Fatalf("Error reading last index: %v", err)
	}

	log.Println("Last Index! ", lastIndex)
	if !isIndexExists(lastIndex + chunkSize) {
		return
	}

	startIndex := lastIndex + 1
	lastIndex = startIndex + chunkSize - 1

	keywords, err := updating.ReadKeywords(sheetsService, channel)
	if err != nil {
		log.Fatalln("failed to find live keywords", zap.Error(err))
		return
	}

	log.Println("Keywords!", keywords)
	// workers <- true
	// <-done
	// go func() {
	// 	newPds := crawling.CrawlDanggnIndex(channel, keywords, startIndex, lastIndex)

	// 	for _, pd := range newPds {
	// 		err := updating.SendDaangnProductToSlack(channel, pd)
	// 		if err != nil {
	// 			// config.Logger.Error("failed to send daangn product to slack", zap.Error(err))
	// 			log.Println("failed to send daangn product to slack", zap.Error(err))
	// 		}
	// 		err = updating.SaveToSheetAtTopFromThirdRow(sheetsService, channel, newPds)
	// 		if err != nil {
	// 			// config.Logger.Error("failed to save to sheet", zap.Error(err))
	// 			log.Println("failed to save to sheet", zap.Error(err))
	// 		}
	// 	}
	// 	<-workers
	// 	done <- true
	// }()

	workers := make(chan bool, numWorkers)
	done := make(chan bool, numWorkers)

	for c := 0; c < numWorkers; c++ {
		done <- true
	}

	rangeSize := (lastIndex - startIndex + 1) / numWorkers

	for c := 0; c < numWorkers; c++ {
		workers <- true
		go func(workerIndex int) {
			defer func() {
				<-workers
				done <- true
			}()

			localStartIndex := startIndex + workerIndex*rangeSize
			localLastIndex := localStartIndex + rangeSize - 1
			if workerIndex == numWorkers-1 {
				localLastIndex = lastIndex
			}

			newPds := crawling.CrawlDanggnIndex(channel, keywords, localStartIndex, localLastIndex)

			for _, pd := range newPds {
				err := updating.SendDaangnProductToSlack(channel, pd)
				if err != nil {
					log.Println("failed to send daangn product to slack", err)
					continue
				}
				err = updating.SaveToSheetAtTopFromThirdRow(sheetsService, channel, newPds)
				if err != nil {
					log.Println("failed to save to sheet", err)
					continue
				}
			}
		}(c)
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
