package crawling

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/essemfly/internal-crawler/internal/domain"
)

func ScrapeArticlesFromBoard(ctx context.Context, pageNum int, cafeBoardURL string) ([]*domain.GuestArticle, error) {
	err := navigateToBoard(ctx, cafeBoardURL)
	if err != nil {
		return nil, err
	}
	articles := []*domain.GuestArticle{}
	for i := 1; i <= pageNum; i++ {
		log.Printf("Navigating to board page (page %d)...", i)
		newArticles, err := extractArticles(ctx)
		if err != nil {
			return nil, err
		}
		articles = append(articles, newArticles...)
		if i+1 > pageNum {
			break
		}
		goToNextPage(ctx, i+1)
	}

	return articles, nil
}

func navigateToBoard(ctx context.Context, cafeBoardURL string) error {
	return chromedp.Run(ctx,
		chromedp.Navigate(cafeBoardURL),
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

func extractArticles(ctx context.Context) ([]*domain.GuestArticle, error) {
	var items []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.WaitReady(`#slideArticleList li`),
		chromedp.Nodes(`#slideArticleList li`, &items),
	)
	if err != nil {
		return nil, err
	}

	var articles []*domain.GuestArticle

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
		articles = append(articles, &domain.GuestArticle{
			TxtDetail:    txtDetail,
			Username:     username,
			CreatedAt:    createdAt,
			ViewCount:    viewCount,
			CommentCount: commentCount,
			URL:          articleURL,
		})
	}

	return articles, nil
}
