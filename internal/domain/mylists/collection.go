package mylists

import "time"

type Collection struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string

	Restaurants []Restaurant `gorm:"foreignKey:CollectionID;references:ID"`

	NaverListID      string
	NaverListName    string
	NaverViewCount   int `gorm:"default:0"`
	NaverSavingCount int `gorm:"default:0"`

	GoogleListID   *string
	GoogleListName *string

	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now();autoUpdateTime:milli"`
}
