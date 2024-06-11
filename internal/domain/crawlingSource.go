package domain

type CrawlingSourceType string

const (
	Youtube   CrawlingSourceType = "Youtube"
	NaverBlog CrawlingSourceType = "NaverBlog"
	Instagram CrawlingSourceType = "Instagram"
)

type CrawlingSource struct {
	Id              int
	SourceName      string
	Type            CrawlingSourceType
	SourceID        string
	NaverListID     string
	NaverListName   string
	SpreadSheetID   *string
	SpreadSheetName *string
	Constraint      *string
}
