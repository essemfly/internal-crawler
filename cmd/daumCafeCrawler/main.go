package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/repository"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/internal/updating"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
)

const (
	PAGE_NUM = 1
	BASEURL  = "https://cafe.daum.net"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	sources := seed.ListSources(domain.DaumCafe)

	daumCafeSvc := repository.NewDaumCafeService()

	for _, channel := range sources {
		ctx, cancel := pkg.OpenChrome()
		defer cancel()

		articles, err := crawling.ScrapeArticlesFromBoard(ctx, PAGE_NUM, channel.SourceID)
		if err != nil || len(articles) == 0 {
			log.Fatalf("Failed to scrape articles: %v", err)
		}

		lastArticle := daumCafeSvc.GetLastArticle(channel.SourceName)

		newArticles := []*domain.GuestArticle{}
		for _, a := range articles {
			if lastArticle != nil && BASEURL+a.URL == lastArticle.URL {
				break
			}
			a.CafeName = channel.SourceName
			a.URL = BASEURL + a.URL
			newArticles = append([]*domain.GuestArticle{a}, newArticles...)
		}

		for _, na := range newArticles {
			updating.SendGuestArticleToSlack(channel, na)
		}

		daumCafeSvc.CreateArticles(newArticles)
	}
}

// func main() {
// 	err := godotenv.Load()
// 	if err != nil && !os.IsNotExist(err) {
// 		fmt.Println("Error loading .env file:", err)
// 		return
// 	}

// 	sources := seed.ListSources(domain.DaumCafe)

// 	for _, channel := range sources {
// 		sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
// 		if err != nil {
// 			log.Fatalf("Error creating Sheets service: %v", err)
// 		}

// 		lastGuestArticle, err := updating.GetLastGuestArticle(sheetsService, channel)
// 		if err != nil {
// 			log.Fatalf("Error fetching last guest article: %v", err)
// 		}

// 		ctx, cancel := pkg.OpenChrome()
// 		defer cancel()

// 		articles, err := crawling.ScrapeArticlesFromBoard(ctx, PAGE_NUM, channel.SourceID)
// 		if err != nil || len(articles) == 0 {
// 			log.Fatalf("Failed to scrape articles: %v", err)
// 		}

// 		newArticles := []*domain.GuestArticle{}
// 		for _, a := range articles {
// 			if lastGuestArticle != nil && BASEURL+a.URL == lastGuestArticle.URL {
// 				break
// 			}
// 			newArticles = append([]*domain.GuestArticle{a}, newArticles...)
// 		}

// 		for _, na := range newArticles {
// 			updating.SendGuestArticleToSlack(channel, na)
// 		}
// 		updating.UpdateGuestArticleCheckpoint(sheetsService, channel, newArticles)
// 	}
// }
