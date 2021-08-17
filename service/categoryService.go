package service

import (
	"log"

	"github.com/sing3demons/go-fiber-mongo/models"
	"github.com/sing3demons/go-fiber-mongo/repository"
)

type CategoryService interface {
	FindAll() ([]models.Category, error)
	Create(category models.Category) (*models.Category, error)
}

type categoryService struct {
	Repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{Repository: repository}
}

func (service *categoryService) FindAll() ([]models.Category, error) {
	category, err := service.Repository.FindAll()
	if err != nil {
		log.Printf("error :%v", err.Error())
	}
	return category, nil
}

func (service *categoryService) Create(category models.Category) (*models.Category, error) {
	result, err := service.Repository.Create(category)
	if err != nil {
		log.Printf("error :%v", err.Error())
	}
	return result, nil
}
