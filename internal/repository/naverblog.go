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

func (ns *NaverBlogService) CreateNaverBlogArticle(article *domain.NaverBlogArticle) error {
	result := ns.db.Create(article)
	return result.Error
}

func (ns *NaverBlogService) CreateNaverBlogArticles(articles []*domain.NaverBlogArticle) error {
	result := ns.db.Create(articles)
	return result.Error
}
