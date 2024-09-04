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
	SourceName      string // Name for identify
	SourceID        string // channel ID, blog ID or URL
	SpreadSheetID   string // Google SpreadSheet ID: 1ufLv1glLILVXP0ZZ5xue9f5JEUp1gauZfQzx9hmBoXY
	SpreadSheetName string // Sheet name in SpreadSheet
	Type            CrawlingSourceType
	NaverListID     string   // Naver Blog List ID
	NaverListName   string   // Naver Map List Name
	WebhookURL      string   // Webhook URL for slack
	Constraint      []string // filterings of source
}
