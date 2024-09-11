package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
		for _, categoryNo := range categories {
			posts, err := crawling.FetchAllBlogPostsByCategory(categoryNo, channel, NUM_WORKERS)
			if err != nil {
				log.Println("Err", err)
				panic(err)
			}
			naverBlogSrvc.CreateNaverBlogArticles(posts)
		}
	}
}
