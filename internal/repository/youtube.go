package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type YoutubeService struct {
	db *gorm.DB
}

func NewYoutubeService() *YoutubeService {
	db := config.DB()
	db.AutoMigrate(&domain.YoutubeVideoStruct{})
	return &YoutubeService{
		db,
	}
}
