package seed

// import (
// 	"os"

// 	"github.com/essemfly/internal-crawler/internal/domain"
// )

// // DEPRECIATED: Get seed crawling sources from excel sheet
// func Seed() []*domain.CrawlingSource {
// 	spreadsheetID := os.Getenv("YOUTUBE_SPREADSHEET_ID")

// 	sheetName1 := "김사원세끼"
// 	kimSawon := domain.CrawlingSource{
// 		Id:              1,
// 		SourceName:      "김사원세끼",
// 		Type:            domain.Youtube,
// 		SourceID:        "UC-x55HF1-IilhxZOzwJm7JA",
// 		NaverListID:     "df3adacda4e34ecf8a457bffded5dd95",
// 		NaverListName:   "김사원세끼 모음",
// 		SpreadSheetID:   spreadsheetID,
// 		SpreadSheetName: sheetName1,
// 	}

// 	sheetName2 := "그시장에가오"
// 	constraint2 := "그 시장을 가오_EP"
// 	bakMarket := domain.CrawlingSource{
// 		Id:              2,
// 		SourceName:      sheetName2,
// 		Type:            domain.Youtube,
// 		SourceID:        "UCyn-K7rZLXjGl7VXGweIlcA",
// 		NaverListID:     "f35771875a2a4bd39b0665c573a330ff",
// 		NaverListName:   "님아 그 시장을 가오",
// 		SpreadSheetID:   spreadsheetID,
// 		SpreadSheetName: sheetName2,
// 		Constraint:      &constraint2,
// 	}

// 	// sheetName3 := "정육왕"
// 	// meatKing := domain.CrawlingSource{
// 	// 	Id:              3,
// 	// 	SourceName:      sheetName3,
// 	// 	Type:            domain.Youtube,
// 	// 	SourceID:        "UC1oXmhvYHVI2bApphh3IzuQ",
// 	// 	NaverListID:     "46ddcb5d43404213aa820b4138e5cdaf",
// 	// 	NaverListName:   "정육왕",
// 	// 	SpreadSheetID:   &spreadsheetID,
// 	// 	SpreadSheetName: &sheetName3,
// 	// }

// 	return []*domain.CrawlingSource{&kimSawon, &bakMarket}
// }
