package registering

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/repository"
)

func AddBlogStoreToList(ctx context.Context, naverBlogSrvc *repository.NaverBlogService, channel *domain.CrawlingSource, links []string) {
	bookmarkSelector := `a[href="#bookmark"]`
	saveButtonSelector := `button.swt-save-btn`
	storeNameSelector := `h1#_header`
	bodySelector := `body`

	saveButtonName := "저장"
	emptyStoreName := "플레이스"
	alreadySaved := "선택됨"

	for idx, link := range links {
		if idx%100 == 0 {
			log.Println("IDX: ", idx)
		}
		var storeName string
		var isSelected string

		log.Println("X0 ", link)
		if link == "" {
			continue
		}
		// 1. StoreName 찾기
		err := chromedp.Run(ctx,
			chromedp.Navigate(link),
			chromedp.WaitVisible(bodySelector, chromedp.ByQuery),
			chromedp.Text(storeNameSelector, &storeName, chromedp.NodeVisible),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("X1 - find store name", storeNameSelector)
		if storeName == emptyStoreName {
			log.Println("EMPTY OUT")
			continue
		}

		log.Println("X2 - GO to bookmark")
		// 2. 저장하기 버튼 누르고, Naver List의 이름과 같은것 찾기 -> 선택됨일 경우 pass
		err = chromedp.Run(ctx,
			chromedp.Click(bookmarkSelector, chromedp.NodeVisible),
			chromedp.Sleep(1*time.Second),
			chromedp.Text(fmt.Sprintf(`//button[.//strong[contains(text(), "%s")]]//span[contains(@class, "swt-save-group-check-area")]//span[contains(@class, "swt-blind")]`, channel.NaverListName), &isSelected, chromedp.NodeVisible),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("X3 - isSelected", isSelected)
		if isSelected == alreadySaved {
			log.Println("X3 - isSelected")
			continue
		}

		log.Println("X4 - Find name of list,", channel.NaverListName)
		// 3. 저장하기에 List click하기
		err = chromedp.Run(ctx,
			chromedp.Click(fmt.Sprintf(`//button[.//strong[contains(text(), "%s")]]`, channel.NaverListName), chromedp.NodeVisible),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("X15 - click save button")
		// 4. 저장버튼찾기
		var buttonText string
		err = chromedp.Run(ctx,
			chromedp.Text(saveButtonSelector, &buttonText, chromedp.NodeVisible),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("X6")
		// 5. 저장버튼 누르기 -> (저장삭제면 안하고 저장일때만 함)
		if buttonText == saveButtonName {
			err = chromedp.Run(ctx,
				chromedp.Click(saveButtonSelector, chromedp.NodeVisible),
				chromedp.Sleep(3*time.Second))
			log.Println("X8")
			if err != nil {
				log.Fatal(err)
			}
		}

		sameArticles, _ := naverBlogSrvc.ListByNaverPlaces(link)
		for _, sa := range sameArticles {
			naverBlogSrvc.UpdateArticleToProcessed(sa.ID)
		}
	}
}

func AddStoreToList(ctx context.Context, channel *domain.CrawlingSource, videos []*domain.YoutubeVideoStruct) {
	storeNameSelector := `h1#_header`
	bookmarkSelector := `a[href="#bookmark"]`
	saveButtonSelector := `button.swt-save-btn`
	bodySelector := `body`

	saveButtonName := "저장"
	emptyStoreName := "플레이스"
	alreadySaved := "선택됨"

	for _, video := range videos {
		if video.NaverLink != "" {
			links := strings.Split(video.NaverLink, ",")
			for _, link := range links {
				var storeName string
				var isSelected string

				// 1. StoreName Select
				err := chromedp.Run(ctx,
					chromedp.Navigate(link),
					chromedp.WaitVisible(bodySelector, chromedp.ByQuery),
					chromedp.Text(storeNameSelector, &storeName, chromedp.NodeVisible),
				)
				if err != nil {
					log.Fatal(err)
				}

				if storeName == emptyStoreName {
					continue
				}

				// 2. 저장하기 버튼 누르고, Naver List의 이름과 같은것 찾기 -> 선택됨일 경우 pass
				err = chromedp.Run(ctx,
					chromedp.Click(bookmarkSelector, chromedp.NodeVisible),
					chromedp.Sleep(1*time.Second),
					chromedp.Text(fmt.Sprintf(`//button[.//strong[contains(text(), "%s")]]//span[contains(@class, "swt-save-group-check-area")]//span[contains(@class, "swt-blind")]`, channel.NaverListName), &isSelected, chromedp.NodeVisible),
				)
				if err != nil {
					log.Fatal(err)
				}

				if isSelected == alreadySaved {
					continue
				}

				// 3. 저장하기에 List click하기
				err = chromedp.Run(ctx,
					chromedp.Click(fmt.Sprintf(`//button[.//strong[contains(text(), "%s")]]`, channel.NaverListName), chromedp.NodeVisible),
					chromedp.Sleep(1*time.Second),
				)
				if err != nil {
					log.Fatal(err)
				}

				// 4. 저장버튼찾기
				var buttonText string
				err = chromedp.Run(ctx,
					chromedp.Text(saveButtonSelector, &buttonText, chromedp.NodeVisible),
				)
				if err != nil {
					log.Fatal(err)
				}

				// 5. 저장버튼 누르기 -> (저장삭제면 안하고 저장일때만 함)
				if buttonText == saveButtonName {
					err = chromedp.Run(ctx,
						chromedp.Click(saveButtonSelector, chromedp.NodeVisible),
						chromedp.Sleep(3*time.Second))
					if err != nil {
						log.Fatal(err)
					}

				}
			}
		}
	}

}
