package repository

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type DaumCafeService struct {
	db *gorm.DB
}

func NewDaumCafeService() *DaumCafeService {
	db := config.DB()
	db.AutoMigrate(&domain.GuestArticle{})
	return &DaumCafeService{
		db,
	}
}

// CreateArticle adds a new GuestArticle to the database
func (s *DaumCafeService) CreateArticle(article *domain.GuestArticle) error {
	result := s.db.Create(article)
	return result.Error
}

// GetArticleByID retrieves a GuestArticle from the database by ID
func (s *DaumCafeService) GetArticleByID(id uint) (*domain.GuestArticle, error) {
	var article domain.GuestArticle
	result := s.db.First(&article, id)
	return &article, result.Error
}

// GetAllArticles retrieves all GuestArticles from the database
func (s *DaumCafeService) GetAllArticles() ([]domain.GuestArticle, error) {
	var articles []domain.GuestArticle
	result := s.db.Find(&articles)
	return articles, result.Error
}

// UpdateArticle updates an existing GuestArticle in the database
func (s *DaumCafeService) UpdateArticle(article *domain.GuestArticle) error {
	result := s.db.Save(article)
	return result.Error
}