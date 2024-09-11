package registering

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/essemfly/internal-crawler/internal/domain"
)

func AddStoreToList(ctx context.Context, channel *domain.CrawlingSource, videos []*domain.YoutubeVideoStruct) {
	bookmarkSelector := `a[href="#bookmark"]`
	saveButtonSelector := `button.swt-save-btn`
	storeNameSelector := `h1#_header`
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

				err = chromedp.Run(ctx,
					chromedp.Click(fmt.Sprintf(`//button[.//strong[contains(text(), "%s")]]`, channel.NaverListName), chromedp.NodeVisible),
					chromedp.Sleep(1*time.Second),
				)
				if err != nil {
					log.Fatal(err)
				}

				var buttonText string
				err = chromedp.Run(ctx,
					chromedp.Text(saveButtonSelector, &buttonText, chromedp.NodeVisible),
				)
				if err != nil {
					log.Fatal(err)
				}

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
