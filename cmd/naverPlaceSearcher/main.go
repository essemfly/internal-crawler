package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/seed"
	"github.com/essemfly/internal-crawler/pkg"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
)

type Place struct {
	Name    string
	Address string
	URL     string
}

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	sources := seed.ListSources(domain.Youtube)
	for _, channel := range sources {
		if channel.SourceName != "먹을텐데" {
			continue
		}

		sheetsService, err := pkg.CreateSheetsService(config.JsonKeyFilePath)
		if err != nil {
			log.Fatalf("Error creating Sheets service: %v", err)
		}

		readRange := fmt.Sprintf("%s!F:F", channel.SpreadSheetName)
		places, err := readPlacesFromSheet(sheetsService, channel.SpreadSheetID, readRange)
		if err != nil {
			log.Fatalf("Error reading from sheet: %v", err)
		}

		log.Println("Found places:", places[30])

		ctx, cancel := pkg.OpenChrome()
		defer cancel()

		for i, place := range places {
			err := searchPlace(ctx, &places[i])
			if err != nil {
				fmt.Printf("Error processing %s at %s: %v\n", place.Name, place.Address, err)
				places[i].URL = ""
			}
			time.Sleep(1 * time.Second) // 페이지 로딩 대기 및 IP 차단 방지
		}

	}

}

func readPlacesFromSheet(srv *sheets.Service, spreadsheetId, readRange string) ([]Place, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	var places []Place
	for _, row := range resp.Values {
		if len(row) > 0 {
			// F열의 각 셀 안에서 상호명과 주소를 줄바꿈으로 분리
			lines := strings.Split(fmt.Sprintf("%v", row[0]), "\n")
			if len(lines) >= 2 {
				place := Place{
					Name:    strings.TrimSpace(lines[0]),
					Address: strings.TrimSpace(lines[1]),
				}
				places = append(places, place)
			}
		}
	}

	return places, nil
}

func searchPlace(ctx context.Context, place *Place) error {
	var naverPlaceURL string
	place.Name = strings.Trim(place.Name, "[]")

	// 첫번쨰 클릭 -> 공유버튼 클릭 -> spi_copyurl가져오기
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://map.naver.com/p/search/`+place.Name),
		chromedp.Sleep(10*time.Second),
	)

	if err != nil {
		return err
	}

	place.URL = naverPlaceURL
	return nil
}
