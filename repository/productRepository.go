package repository

import (
	"context"
	"time"

	"github.com/sing3demons/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindOne(filter primitive.M) (*models.Product, error)
	Create(product models.Product) (*models.Product, error)
	Update(filter primitive.M, update primitive.D) error
	Delete(filter primitive.M) error
}

type productRepository struct {
	DB *mongo.Database
}

func (tx *productRepository) collection() *mongo.Collection {
	return tx.DB.Collection("products")
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{DB: db}
}

func (tx *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tx.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (tx *productRepository) FindOne(filter primitive.M) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product models.Product

	if err := tx.collection().FindOne(ctx, filter).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (tx *productRepository) Create(product models.Product) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := tx.collection().InsertOne(ctx, product)
	if err != nil || result == nil {
		return nil, err
	}

	return &product, nil
}

func (tx *productRepository) Update(filter primitive.M, update primitive.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := tx.collection().FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return err
	}
	return nil
}

func (tx *productRepository) Delete(filter primitive.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := tx.collection().FindOneAndDelete(ctx, filter).Err(); err != nil {
		return err
	}
	return nil
}
