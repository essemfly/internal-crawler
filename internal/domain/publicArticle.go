package domain

import "time"

type PublicType string

const (
	BlueRibbon PublicType = "BlueRibbon"
	Michelin   PublicType = "Michelin"
)

type PublicArticle struct {
	ID          uint `gorm:"primaryKey"` // Primary key field
	Title       string
	Description string
	Address     string
	Notes       string
	Type        PublicType

	Latitude     *float64
	Longitude    *float64
	Link         string
	ThumbnailUrl string

	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (PublicArticle) TableName() string {
	return "public_articles"
}
