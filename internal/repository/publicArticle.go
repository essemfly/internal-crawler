package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type PublicArticleService struct {
	db *gorm.DB
}

func NewPublicArticleService() *PublicArticleService {
	db := config.DB()
	db.AutoMigrate(&domain.PublicArticle{})
	return &PublicArticleService{
		db,
	}
}

func (s *PublicArticleService) CreatePublicArticle(article *domain.PublicArticle) error {
	return s.db.Create(article).Error
}

func (s *PublicArticleService) CreatePublicArticles(articles []*domain.PublicArticle) error {
	return s.db.CreateInBatches(articles, 100).Error
}

func (s *PublicArticleService) GetPublicArticle(id uint) (*domain.PublicArticle, error) {
	var article domain.PublicArticle
	return &article, s.db.First(&article, id).Error
}

func (s *PublicArticleService) GetAllPublicArticles() ([]*domain.PublicArticle, error) {
	var articles []*domain.PublicArticle
	return articles, s.db.Find(&articles).Error
}

func (s *PublicArticleService) UpdatePublicArticle(article *domain.PublicArticle) error {
	return s.db.Save(article).Error
}
