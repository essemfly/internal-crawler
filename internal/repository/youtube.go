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

func (y *YoutubeService) SaveVideo(video *domain.YoutubeVideoStruct) {
	y.db.Create(video)
}

func (y *YoutubeService) SaveVideoList(videos []*domain.YoutubeVideoStruct) {
	y.db.Create(videos)
}

func (y *YoutubeService) GetLastVideo(channel string) (*domain.YoutubeVideoStruct, error) {
	var video domain.YoutubeVideoStruct
	result := y.db.Where("channel = ?", channel).
		Order("published_at DESC").
		First(&video)

	if result.Error != nil {
		return nil, result.Error
	}
	return &video, nil
}
