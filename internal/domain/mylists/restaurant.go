package mylists

import "time"

type Restaurant struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Name          string
	Address       *string
	Latitude      *float64
	Longitude     *float64
	GooglePlaceID *string
	NaverPlaceID  *string

	ArticleID    uint
	Article      Article   `gorm:"foreignKey:ArticleID;references:ID"`
	CreatedAt    time.Time `gorm:"default:now()"`
	UpdatedAt    time.Time `gorm:"default:now();autoUpdateTime:milli"`
	CollectionID *uint
	Collection   Collection `gorm:"foreignKey:CollectionID;references:ID"`
}
