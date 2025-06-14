package mylists

import (
	"time"

	"github.com/essemfly/internal-crawler/internal/domain"
)

type Article struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	Title        string
	Description  string
	Link         string
	ThumbnailUrl string

	crawlingSource   domain.CrawlingSource
	crawlingSourceID uint
	PublishedAt      time.Time

	CreatedAt  time.Time
	UpdatedAt  time.Time
	Restaurant []Restaurant
}
