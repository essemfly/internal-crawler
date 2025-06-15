package seed

import (
	"log"
	"os"

	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/internal/repository"
)

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

	meatKing := domain.CrawlingSource{
		Type:          domain.Youtube,
		SourceName:    "정육왕",
		SourceID:      "UC1oXmhvYHVI2bApphh3IzuQ",
		NaverListName: "정육왕 모음",
		WebhookURL:    os.Getenv("YOUTUBE_WEBHOOK"),
	}

	blueRibbon := domain.CrawlingSource{
		Type:          domain.Public,
		SourceName:    "블루리본",
		SourceID:      "UC1oXmhvYHVI2bApphh3IzuQ",
		NaverListName: "블루리본 모음",
		WebhookURL:    os.Getenv("YOUTUBE_WEBHOOK"),
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

	travelYoutubes := []domain.CrawlingSource{
		{
			Type:            domain.Youtube,
			SourceName:      "빠니보틀",
			SourceID:        "UCNhofiqfw5nl-NeDJkXtPvw",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "아무거나보틀",
			SourceID:        "UCYbxBWWLTBZTWXjtMUjk_eA",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "곽튜브",
			SourceID:        "UClRNDVO8093rmRTtLe4GEPw",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "뜨랑낄로",
			SourceID:        "UCWqWR1sFAz9UsjwkFHEkBsw",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "여행가제이",
			SourceID:        "UCxU8QX7IRRIW0VLuoWWoxbw",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "체코제",
			SourceID:        "UCaoqDZPllYXLAH_5OBRLLrw",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "원지의하루",
			SourceID:        "UC9gxOp_-R78phMHmv2bW_sg",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "천천히 세계여행 앤젤리나",
			SourceID:        "UCGwUygnNcB1kcbmWl_jfzlQ",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "캡틴따거",
			SourceID:        "UCt_7uH4Igz0T_K3Qzbs1Wig",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "노마드션",
			SourceID:        "UCfCOEG2kjX_x4KAdWX-YUcA",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "모칠레로",
			SourceID:        "UCUcy82tGagXlj4t1p2todeQ",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "버드모이",
			SourceID:        "UCOoM7iaJkVWAmITWHKTPhnQ",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
		},
		{
			Type:            domain.Youtube,
			SourceName:      "잰잰바리",
			SourceID:        "UCSxhYq6K0mxF24SmMGeNNQA",
			SpreadSheetID:   "",
			SpreadSheetName: "",
			NaverListID:     "",
			NaverListName:   "",
			WebhookURL:      os.Getenv("TRAVELER_WEBHOOK"),
			Constraint:      "",
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
	crwlSrvc.CreateCrawlingSource(&meatKing)
	crwlSrvc.CreateCrawlingSource(&blueRibbon)
	crwlSrvc.CreateCrawlingSource(&wishket)
	crwlSrvc.CreateCrawlingSource(&basketGuest)
	crwlSrvc.CreateCrawlingSource(&daangn)
	for _, blogSource := range blogSources {
		crwlSrvc.CreateCrawlingSource(&blogSource)
	}
	for _, travelYoutube := range travelYoutubes {
		crwlSrvc.CreateCrawlingSource(&travelYoutube)
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
