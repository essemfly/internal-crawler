package seed

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Depreciated
func getRegisteredSourcesFromSheet() []*domain.CrawlingSource {
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	sheetName := "channels"
	readRange := sheetName + "!A2:H"

	ctx := context.Background()
	creds := option.WithCredentialsFile(config.JsonKeyFilePath) // Replace with your credentials file path
	srv, err := sheets.NewService(ctx, creds)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var sources []*domain.CrawlingSource
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			typeInString := row[1].(string)
			typeInStruct := domain.Youtube
			switch typeInString {
			case "wishket":
				typeInStruct = domain.Wishket
			case "daangn":
				typeInStruct = domain.Daangn
			case "daumcafe":
				typeInStruct = domain.DaumCafe
			case "naverblog":
				typeInStruct = domain.NaverBlog
			}

			constraints := row[4].(string)

			source := &domain.CrawlingSource{
				SourceName:      row[0].(string),
				SourceID:        row[3].(string),
				SpreadSheetID:   spreadsheetID,
				SpreadSheetName: row[5].(string),
				Type:            typeInStruct,
				NaverListID:     "",
				NaverListName:   "",
				WebhookURL:      row[2].(string),
				Constraint:      constraints,
			}

			// Check if row[6] and row[7] exist and set them appropriately
			if len(row) > 6 {
				source.NaverListID = row[6].(string)
			}
			if len(row) > 7 {
				source.NaverListName = row[7].(string)
			}
			sources = append(sources, source)
		}
	}
	return sources
}

func CrawlingSeeds() {

	kimSawon := domain.CrawlingSource{
		Type:            domain.Youtube,
		SourceName:      "김사원세끼",
		SourceID:        "UC-x55HF1-IilhxZOzwJm7JA",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "김사원세끼",
		NaverListID:     "df3adacda4e34ecf8a457bffded5dd95",
		NaverListName:   "김사원세끼 모음",
		WebhookURL:      os.Getenv("YOUTUBE_WEBHOOK"),
	}

	bakMarket := domain.CrawlingSource{
		Type:            domain.Youtube,
		SourceName:      "그시장에가오",
		SourceID:        "UCyn-K7rZLXjGl7VXGweIlcA",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "그시장에가오",
		NaverListID:     "f35771875a2a4bd39b0665c573a330ff",
		NaverListName:   "님아 그 시장을 가오",
		WebhookURL:      os.Getenv("YOUTUBE_WEBHOOK"),
		Constraint:      "그 시장을 가오_EP",
	}

	sikyung := domain.CrawlingSource{
		Type:            domain.Youtube,
		SourceName:      "먹을텐데",
		SourceID:        "UCl23-Cci_SMqyGXE1T_LYUg",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "먹을텐데",
		NaverListID:     "7a6f094ce74a4617af423bf2fb0c4582",
		NaverListName:   "성시경의 먹을텐데",
		WebhookURL:      os.Getenv("YOUTUBE_WEBHOOK"),
		Constraint:      "먹을텐데",
	}

	wishket := domain.CrawlingSource{
		Type:            domain.Wishket,
		SourceName:      "위시켓",
		SourceID:        "https://www.wishket.com/project/?d=M4JwLgvAdgpg7gMhgYwCYQCogK4yA%3D%3D%3D",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "wishket",
		WebhookURL:      os.Getenv("WISHKET_WEBHOOK"),
	}

	basketGuest := domain.CrawlingSource{
		Type:            domain.DaumCafe,
		SourceName:      "동아리농구방",
		SourceID:        "https://m.cafe.daum.net/dongarry/Dilr",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "농구게스트",
		WebhookURL:      os.Getenv("BASKET_WEBHOOK"),
	}

	blogSources := []domain.CrawlingSource{
		{SourceName: "mardukas",
			Type:            domain.NaverBlog,
			SourceID:        "mardukas",
			Constraint:      "9,61,62,111,65,66,1,68,70,67",
			SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
			SpreadSheetName: "mardukas",
			WebhookURL:      os.Getenv("YOUTUBE_WEBHOOK"),
		},
		{SourceName: "paperchan",
			Type:            domain.NaverBlog,
			SourceID:        "paperchan",
			Constraint:      "5,2,3,4,15",
			SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
			SpreadSheetName: "paperchan",
			WebhookURL:      os.Getenv("YOUTUBE_WEBHOOK"),
		},
	}

	daangn := domain.CrawlingSource{
		Type:            domain.Daangn,
		SourceName:      "당근마켓",
		SpreadSheetID:   "1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY",
		SpreadSheetName: "당근마켓",
		WebhookURL:      os.Getenv("DAANGN_WEBHOOK"),
	}

	daangnConfig := domain.DaangnConfig{
		CurrentIdx: 830096400,
	}

	daangnKeywords := domain.DaangnKeyword{
		Keyword: "농구",
	}

	crwlSrvc := repository.NewCrawlingService()

	crwlSrvc.CreateCrawlingSource(&kimSawon)
	crwlSrvc.CreateCrawlingSource(&bakMarket)
	crwlSrvc.CreateCrawlingSource(&sikyung)
	crwlSrvc.CreateCrawlingSource(&wishket)
	crwlSrvc.CreateCrawlingSource(&basketGuest)
	crwlSrvc.CreateCrawlingSource(&daangn)
	for _, blogSource := range blogSources {
		crwlSrvc.CreateCrawlingSource(&blogSource)
	}

	daangnSrvc := repository.NewDaangnService()
	daangnSrvc.CreateDaangnConfig(&daangnConfig)
	daangnSrvc.CreateDaangnKeyword(&daangnKeywords)

}

func ListSources(sourceType domain.CrawlingSourceType) []*domain.CrawlingSource {

	crwlSrvc := repository.NewCrawlingService()

	sources, err := crwlSrvc.ListAllCrawlingSources()
	if err != nil {
		log.Println("Error getting sources: ", err)
		return nil
	}
	var filteredSources []*domain.CrawlingSource
	for _, source := range sources {
		if source.Type == sourceType {
			filteredSources = append(filteredSources, source)
		}
	}
	return filteredSources
}
