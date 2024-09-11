package repository

import (
	"time"

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

func (ns *NaverBlogService) GetLatestArticleDate(blogId string) (time.Time, error) {
	var article domain.NaverBlogArticle
	result := ns.db.Where("channel = ?", blogId).Order("created_at desc").First(&article)
	return article.CreatedAt, result.Error
}
