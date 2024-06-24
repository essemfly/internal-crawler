package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/essemfly/internal-crawler/internal/registering"
)

const (
	CAFE     = "dongarry"
	BOARD    = "Dilr"
	PAGE_NUM = 1
)

type Article struct {
	URL          string
	TxtDetail    string
	Username     string
	CreatedAt    string
	ViewCount    string
	CommentCount string
}

func main() {
	ctx, cancel := registering.OpenChrome()
	defer cancel()
	log.Println("1")

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// log.Println("11")
	// if err := login(ctx); err != nil {
	// 	log.Fatalf("Failed to login: %v", err)
	// }

	if err := scrapeArticles(ctx); err != nil {
		log.Fatalf("Failed to scrape articles: %v", err)
	}
}

// func login(ctx context.Context) error {
// 	var loginURL = fmt.Sprintf("https://m.cafe.daum.net/%s", CAFE)

// 	err := chromedp.Run(ctx,
// 		chromedp.Navigate(loginURL),
// 		chromedp.Click(`#daumMinidaum > a`),
// 		chromedp.WaitVisible(`#loginId--1`),
// 		chromedp.SendKeys(`#loginId--1`, USERNAME),
// 		chromedp.SendKeys(`#password--2`, PASSWORD),
// 		chromedp.Click(`#mainContent > div > div > form > div.confirm_btn > button.btn_g.highlight.submit`),
// 		chromedp.Sleep(15*time.Second),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func scrapeArticles(ctx context.Context) error {
	err := navigateToBoard(ctx)
	if err != nil {
		return err
	}
	for i := 1; i <= PAGE_NUM; i++ {
		log.Printf("Navigating to board page (page %d)...", i)
		_, err = extractArticles(ctx)
		if err != nil {
			return err
		}
		if i+1 > PAGE_NUM {
			break
		}
		goToNextPage(ctx, i+1)
	}

	return nil
}

func navigateToBoard(ctx context.Context) error {
	boardURL := fmt.Sprintf("https://m.cafe.daum.net/%s/%s", CAFE, BOARD)

	return chromedp.Run(ctx,
		chromedp.Navigate(boardURL),
		chromedp.Sleep(time.Second*2),
	)
}

func goToNextPage(ctx context.Context, currentPage int) error {
	nextPageSelector := fmt.Sprintf("#pagingNav > span:nth-child(%d)", currentPage%5+2)
	if currentPage%5 == 0 {
		nextPageSelector = "#mArticle > div.paging_board > a.btn_page.btn_next"
	}

	return chromedp.Run(ctx,
		chromedp.Click(nextPageSelector),
		chromedp.Sleep(2*time.Second),
	)
}

func extractArticles(ctx context.Context) ([]Article, error) {
	var items []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.WaitReady(`#slideArticleList li`),
		chromedp.Nodes(`#slideArticleList li`, &items),
	)
	if err != nil {
		return nil, err
	}

	var articles []Article

	for _, item := range items {
		var txtDetail, username, createdAt, viewCount, commentCount, articleURL string

		err := chromedp.Run(ctx,
			chromedp.Text(`.txt_detail`, &txtDetail, chromedp.ByQuery, chromedp.FromNode(item)),
			chromedp.Text(`.username`, &username, chromedp.ByQuery, chromedp.FromNode(item)),
			chromedp.Text(`.created_at`, &createdAt, chromedp.ByQuery, chromedp.FromNode(item)),
			chromedp.Text(`.view_count`, &viewCount, chromedp.ByQuery, chromedp.FromNode(item)),
			chromedp.Text(`.num_cmt`, &commentCount, chromedp.ByQuery, chromedp.FromNode(item)),
			chromedp.AttributeValue(`a.link_cafe`, "href", &articleURL, nil, chromedp.ByQuery, chromedp.FromNode(item)),
		)
		if err != nil {
			return nil, err
		}

		// Append the extracted details to the articles list
		articles = append(articles, Article{
			TxtDetail:    txtDetail,
			Username:     username,
			CreatedAt:    createdAt,
			ViewCount:    viewCount,
			CommentCount: commentCount,
			URL:          articleURL,
		})
	}

	log.Println("Articles:", articles)
	return articles, nil
}
