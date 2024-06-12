package domain

type CrawlingSourceType string

const (
	Youtube   CrawlingSourceType = "Youtube"
	NaverBlog CrawlingSourceType = "NaverBlog"
	Instagram CrawlingSourceType = "Instagram"
	DaumCafe  CrawlingSourceType = "DaumCafe"
	Wishket   CrawlingSourceType = "Wishket"
	Daangn    CrawlingSourceType = "Daangn"
)

type CrawlingSource struct {
	SourceName      string
	SourceID        string
	SpreadSheetID   string
	SpreadSheetName string
	Type            CrawlingSourceType
	NaverListID     string
	NaverListName   string
	WebhookURL      string
	Constraint      []string
}
