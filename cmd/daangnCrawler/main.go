package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"

	"go.uber.org/zap"
)

const (
	chunkSize        = 500
	numWorkers       = 4
	waitTime         = 2
	numRetries       = 8
	randomRange      = 30
	GlobalStartIndex = 783940000 // 2023-02-14 17:00:00
	// 788540500 : 2024-06-20 10:00:00
	// 783940000 : 2024-06-11 14:00:00
	// 533830000 : 2023-02-13 04:00:00
	// 533500000 : 2023-02-12 16:00:00
	// 531250000 : 2023-02-08 15:00:00
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func initialize() (*domain.CrawlingSource, *sheets.Service, error) {
	sources := seed.ListSources(domain.Daangn)
	channel := sources[0]

	log.Println("Initialize with channel: ", channel.SourceName)
	sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
	if err != nil {
		return channel, nil, err
	}

	return channel, sheetsService, nil
}

func crawlAndProcess(channel *domain.CrawlingSource, sheetsService *sheets.Service, keywords []string, startIndex, lastIndex int) error {
	var wg sync.WaitGroup
	errChan := make(chan error, numWorkers)
	rangeSize := (lastIndex - startIndex + 1) / numWorkers

	for workerIndex := 0; workerIndex < numWorkers; workerIndex++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			localStartIndex := startIndex + workerIndex*rangeSize
			localLastIndex := localStartIndex + rangeSize - 1
			if workerIndex == numWorkers-1 {
				localLastIndex = lastIndex
			}

			if err := crawlAndSave(channel, sheetsService, keywords, localStartIndex, localLastIndex); err != nil {
				errChan <- err
			}
		}(workerIndex)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func crawlAndSave(channel *domain.CrawlingSource, sheetsService *sheets.Service, keywords []string, startIndex, lastIndex int) error {
	newPds, err := crawling.CrawlDanggnIndex(channel, keywords, startIndex, lastIndex)
	if err != nil {
		return fmt.Errorf("forbidden error encountered in CrawlDanggnIndex: %v", err)
	}

	for _, pd := range newPds {
		if err := updating.SendDaangnProductToSlack(channel, pd); err != nil {
			log.Println("Failed to send daangn product to slack:", err)
			continue
		}
		if err := updating.SaveToSheetAppend(sheetsService, channel, newPds); err != nil {
			log.Println("Failed to save to sheet:", err)
			continue
		}
	}
	return nil
}

func main() {
	if err := loadEnv(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	channel, sheetsService, err := initialize()
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	lastIndex, err := updating.ReadLastIndex(sheetsService, channel)
	if err != nil {
		log.Fatalf("Error reading last index: %v", err)
	}

	keywords, err := updating.ReadKeywords(sheetsService, channel)
	if err != nil {
		log.Fatalln("failed to find live keywords", zap.Error(err))
		return
	}

	log.Println("Last Index with keywords ", lastIndex, keywords)

	for {
		if !isIndexExists(lastIndex + chunkSize) {
			break
		}

		time.Sleep(waitTime * time.Second)
		startIndex := lastIndex + 1
		lastIndex = startIndex + chunkSize - 1

		if err := crawlAndProcess(channel, sheetsService, keywords, startIndex, lastIndex); err != nil {
			log.Fatalf("Processing error: %v", err)
		}
		lastIndex += chunkSize
		if err := updating.UpdateLastIndex(sheetsService, channel, lastIndex); err != nil {
			log.Fatalf("Error updating last index: %v", err)
		}
	}
}

func isIndexExists(index int) bool {
	_, err := crawling.CrawlPage(index)
	errCounts := 0

	for err != nil {
		errCounts += 1
		if errCounts > numRetries {
			return false
		}

		n := rand.Intn(randomRange)
		_, err = crawling.CrawlPage(index + n)
	}

	return true
}
