package repository

import (
	"context"
	"time"

	"github.com/sing3demons/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	Create(category models.Category) (*models.Category, error)
}

type categoryRepository struct {
	DB *mongo.Database
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepository{DB: db}
}

func (tx *categoryRepository) collection() *mongo.Collection {
	return tx.DB.Collection("category")
}

func (tx *categoryRepository) FindAll() ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tx.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var categories []models.Category

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (tx *categoryRepository) Create(category models.Category) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := tx.collection().InsertOne(ctx, category)
	if err != nil || result == nil {
		return nil, err
	}

	return &category, nil
}
