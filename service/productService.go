package service

import (
	"log"

	"github.com/sing3demons/go-fiber-mongo/models"
	"github.com/sing3demons/go-fiber-mongo/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	FindAll() ([]models.Product, error)
	FindOne(filter primitive.M) (*models.Product, error)
	Create(product models.Product) (*models.Product, error)
	Update(filter primitive.M, update primitive.D) error
	Delete(filter primitive.M) error
}

type productService struct {
	Repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{Repository: repository}
}

func (service *productService) FindAll() ([]models.Product, error) {
	product, err := service.Repository.FindAll()
	if err != nil {
		log.Printf("error :%v", err.Error())
	}

	return product, nil
}

func (service *productService) FindOne(filter primitive.M) (*models.Product, error) {
	product, err := service.Repository.FindOne(filter)
	if err != nil {
		log.Printf("error :%v", err.Error())
	}

	return product, nil
}

func (service *productService) Create(product models.Product) (*models.Product, error) {
	result, err := service.Repository.Create(product)
	if err != nil {
		log.Printf("error :%v", err.Error())
	}

	return result, nil
}

func (service *productService) Update(filter primitive.M, update primitive.D) error {
	err := service.Repository.Update(filter, update)
	if err != nil {
		log.Printf("error :%v", err.Error())
	}

	return nil
}

func (service *productService) Delete(filter primitive.M) error {
	err := service.Repository.Delete(filter)
	if err != nil {
		log.Printf("error :%v", err.Error())
	}

	return nil
}
