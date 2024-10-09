package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/repository"
	"github.com/essemfly/internal-crawler/internal/seed"
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

	naverBlogSrvc := repository.NewNaverBlogService()

	sources := seed.ListSources(domain.NaverBlog)
	for _, channel := range sources {
		categories := strings.Split(channel.Constraint, ",")
		latestDate, err := naverBlogSrvc.GetLatestArticleDate(channel.SourceName)
		if err != nil {
			log.Println("Error getting latest article date:", err)
			latestDate = time.Date(2008, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		for _, categoryNo := range categories {
			posts, err := crawling.FetchAllBlogPostsByCategory(categoryNo, channel, NUM_WORKERS, latestDate)
			if err != nil {
				log.Println("Err", err)
				panic(err)
			}

			if len(posts) > 0 {
				naverBlogSrvc.CreateNaverBlogArticles(posts)
				log.Printf("Added %d new posts for category %s\n", len(posts), crawling.GetContentKeyValues(channel.SourceName, categoryNo))
			} else {
				log.Printf("No new posts for category %s\n", crawling.GetContentKeyValues(channel.SourceName, categoryNo))
			}
		}
	}
}
