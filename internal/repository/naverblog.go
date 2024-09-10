package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type NaverBlogService struct {
	db *gorm.DB
}

func NewNaverBlogService() *NaverBlogService {
	db := config.DB()
	db.AutoMigrate(&domain.NaverBlogArticle{})
	return &NaverBlogService{
		db,
	}
}
