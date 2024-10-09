package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/registering"
	"github.com/essemfly/internal-crawler/internal/repository"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
)

// Used when crawls whole sheets of an youtube channel and add to list
func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ctx, cancel := pkg.OpenChrome()
	registering.NaverLogin(ctx)

	naverBlogSrvc := repository.NewNaverBlogService()

	sources := seed.ListSources(domain.NaverBlog)
	for _, channel := range sources {
		articles, err := naverBlogSrvc.ListUnprocessedArticles()
		if err != nil {
			fmt.Println("Error in listing", err)
		}

		// registering.AddStoreToList(ctx, channel, articles)
		registering.AddBlogStoreToList(ctx, naverBlogSrvc, channel, articles)
	}

	defer cancel()
}

func test() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ctx, cancel := pkg.OpenChrome()
	registering.NaverLogin(ctx)
	defer cancel()

	testSource := &domain.CrawlingSource{
		NaverListName: "테스트",
		NaverListID:   "1123eae3eca246a78d0b846683bf1a5c",
	}

	video := &domain.YoutubeVideoStruct{
		NaverLink: "https://map.naver.com/p/entry/place/37163963",
	}
	registering.AddStoreToList(ctx, testSource, []*domain.YoutubeVideoStruct{video})
	registering.AddStoreToList(ctx, testSource, []*domain.YoutubeVideoStruct{video})
}
