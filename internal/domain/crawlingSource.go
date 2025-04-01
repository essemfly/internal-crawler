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
	ID              uint `gorm:"primaryKey"` // ID is the primary key
	Type            CrawlingSourceType
	IsActive        bool `gorm:"default:true"`
	SourceName      string
	SourceID        string // Channel ID for Youtube, Blog ID for NaverBlog, User ID for Instagram, Cafe ID for DaumCafe
	SpreadSheetID   string // Sheet ID not to be used in DB
	SpreadSheetName string // Sheet name not to be used in DB
	NaverListID     string
	NaverListName   string
	WebhookURL      string // webhook URL for slack notification
	Constraint      string // Constraint for filters
}

func (CrawlingSource) TableName() string {
	return "crawling_sources"
}
