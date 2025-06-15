package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	log.Println("Starting Blue Ribbon crawler...")

	ribbonTypes := []crawling.RibbonType{crawling.RIBBON_ONE, crawling.RIBBON_TWO, crawling.RIBBON_THREE}

	repository := repository.NewPublicArticleService()

	for _, ribbonType := range ribbonTypes {
		log.Printf("Crawling %s", ribbonType)
		articles, err := crawling.CrawlBlueRibbon(ribbonType)
		if err != nil {
			log.Printf("Error crawling %s: %v", ribbonType, err)
			continue
		}

		log.Printf("Found %d restaurants for %s", len(articles), ribbonType)
		repository.CreatePublicArticles(articles)

		time.Sleep(3 * time.Second)
	}

}
