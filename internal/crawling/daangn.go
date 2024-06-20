package crawling

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const (
	ProductURL = "https://www.daangn.com/articles/"
)

func CrawlDanggnIndex(channel *domain.CrawlingSource, keywords []string, startIndex, lastIndex int) []*domain.DaangnProduct {
	log.Println("start crawling danggn index", zap.Int("startIndex", startIndex), zap.Int("lastIndex", lastIndex))
	pds := []*domain.DaangnProduct{}
	for i := startIndex; i <= lastIndex; i++ {
		newProduct, err := CrawlPage(i)
		if err != nil {
			if err.Error() == "Not Found" {
				continue
			}

			// config.Logger.Error("failed to crawl page", zap.Error(err))
			log.Fatalln("failed to crawl page", zap.Error(err))
			continue
		}

		pds = addProductForKeywords(pds, newProduct, keywords)
	}

	return pds
}

func CrawlPage(index int) (*domain.DaangnProduct, error) {
	if index%500 == 0 {
		// config.Logger.Info("start crawling danggn page", zap.Int("index", index))
		log.Println("start crawling danggn page", zap.Int("index", index))
	}
	c := colly.NewCollector(
		colly.AllowedDomains("www.daangn.com"),
		colly.UserAgent("Mozilla/5.2 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)
	c.SetRequestTimeout(10 * time.Second)

	indexStr := strconv.Itoa(index)
	url := ProductURL + indexStr

	newProduct := domain.DaangnProduct{
		DanggnIndex: indexStr,
		Url:         url,
	}

	c.OnHTML("head", func(e *colly.HTMLElement) {
		// Find the value of the product:availability meta header
		availability := e.ChildAttr("meta[name='product:availability']", "content")
		if availability == "oos" {
			newProduct.Status = domain.DANGGN_STATUS_SOLDOUT
		} else {
			newProduct.Status = domain.DANGGN_STATUS_SALE
		}

		regionName := e.ChildAttr("meta[name='twitter:data2']", "content")
		newProduct.SellerRegionName = regionName
	})

	c.OnHTML("#article-images", func(e *colly.HTMLElement) {
		e.ForEach("img", func(_ int, e *colly.HTMLElement) {
			imageUrl := e.Attr("data-lazy")
			imageUrl = strings.Replace(imageUrl, "s=1440x1440", "s=480x480", 1)
			newProduct.Images = append(newProduct.Images, imageUrl)
		})
	})

	c.OnHTML("#article-profile", func(e *colly.HTMLElement) {
		nickName := e.ChildText("#nickname")
		temperature := e.ChildText("#temperature-wrap dd")

		newProduct.SellerNickName = nickName
		newProduct.SellerTemperature = temperature
	})

	c.OnHTML("#article-description", func(e *colly.HTMLElement) {
		categoryAndWrittenDate := e.ChildText("#article-category")
		spliter := "∙"
		categoryParsers := strings.Split(categoryAndWrittenDate, spliter)
		if len(categoryParsers) > 1 {
			if strings.Contains(categoryParsers[1], "분") {
				numMinutes, _ := pkg.ExtractIntFromString(categoryParsers[1])
				now := time.Now()
				newProduct.WrittenAt = now.Add(-1 * time.Duration(numMinutes) * time.Minute)
			} else if strings.Contains(categoryParsers[1], "시간") {
				numHours, _ := pkg.ExtractIntFromString(categoryParsers[1])
				now := time.Now()
				newProduct.WrittenAt = now.Add(-1 * time.Duration(numHours) * time.Hour)
			} else if strings.Contains(categoryParsers[1], "일") {
				numDays, _ := pkg.ExtractIntFromString(categoryParsers[1])
				now := time.Now()
				newProduct.WrittenAt = now.Add(-1 * time.Duration(numDays) * 24 * time.Hour)
			}
		}

		title := e.ChildText("#article-title")
		price := e.ChildText("#article-price")
		description := e.ChildText("#article-detail")
		articleCounts := e.ChildText("#article-counts")

		newProduct.Name = title
		newProduct.Price = pkg.ParsePriceString(price)
		newProduct.Description = description

		likeCount, viewCount, chatCount := pkg.ParseViewCountsString(articleCounts)
		newProduct.CrawlCategory = categoryParsers[0]
		newProduct.LikeCounts = likeCount
		newProduct.ViewCounts = viewCount
		newProduct.ChatCounts = chatCount
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func addProductForKeywords(pds []*domain.DaangnProduct, product *domain.DaangnProduct, keywords []string) []*domain.DaangnProduct {
	for _, keyword := range keywords {
		if strings.Contains(product.Name, keyword) {
			product.Keyword = keyword
			pds = append(pds, product)
		}
	}
	return pds
}
