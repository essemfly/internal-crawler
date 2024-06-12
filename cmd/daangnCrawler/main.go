package main

import (
	"log"
	"math/rand"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"go.uber.org/zap"
)

const (
	chunkSize        = 100
	numWorkers       = 5
	GlobalStartIndex = 783940000 // 2023-02-14 17:00:00
	// 783940000 : 2024-06-11 14:00:00
	// 533830000 : 2023-02-13 04:00:00
	// 533500000 : 2023-02-12 16:00:00
	// 531250000 : 2023-02-08 15:00:00
)

func DanggnCrawler() {

	workers := make(chan bool, numWorkers)
	done := make(chan bool, numWorkers)

	for c := 0; c < numWorkers; c++ {
		done <- true
	}

	lastIndex := getLastIndex()
	log.Println("Last Index! ", lastIndex)
	for isIndexExists(lastIndex + chunkSize) {
		startIndex := lastIndex + 1
		lastIndex = startIndex + chunkSize - 1

		keywords, err := config.Repo.DaangnKeywords.FindLiveKeywords()
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

func getLastIndex() int {
	lastIndex, err := config.Repo.CrawlThreads.FindLastIndex()
	if err != nil || lastIndex == 0 {
		return GlobalStartIndex
	}
	return lastIndex
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
