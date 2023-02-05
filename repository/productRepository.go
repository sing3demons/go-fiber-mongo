package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sing3demons/go-fiber-mongo/database"
	"github.com/sing3demons/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindOne(filter primitive.M) (*models.Product, error)
	Create(product models.Product) (*models.Product, error)
	Update(filter primitive.M, docs []interface{}) error
	Delete(filter primitive.M) error
}

type productRepository struct {
	DB    *mongo.Database
	Cache database.RedisCache
}

func (tx *productRepository) collection() *mongo.Collection {
	return tx.DB.Collection("products")
}

func NewProductRepository(db *mongo.Database, cache database.RedisCache) ProductRepository {
	return &productRepository{DB: db, Cache: cache}
}

func (tx *productRepository) FindAll() ([]models.Product, error) {

	cacheProducts, _ := tx.Cache.GetProducts("products")

	if cacheProducts != nil {
		fmt.Println("Get...Redis")
		return cacheProducts, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// id, _ := primitive.ObjectIDFromHex("61a30848e3201fba74a4a8bb")
	// conditions = append(conditions, bson.M{"$match": bson.M{"_id": id}})

	lookup := bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "category"},
				{Key: "localField", Value: "categoryId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "categoryId"},
			},
		},
	}

	unwind := bson.D{{Key: "$unwind", Value: "$categoryId"}}

	project := bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "name", Value: 1},
				{Key: "desc", Value: 1},
				{Key: "image", Value: 1},
				{Key: "categoryId", Value: "$categoryId._id"},
				{Key: "category", Value: "$categoryId"},
			},
		},
	}

	pipeline := bson.A{
		lookup,
		unwind,
		project,
	}

	var products []models.Product

	cursor, err := tx.collection().Aggregate(ctx, pipeline)
	// cursor, err := tx.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	tx.Cache.Set("products", products)

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

func (tx *productRepository) Update(filter primitive.M, docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := tx.collection().UpdateMany(ctx, filter, docs)
	if err != nil {
		return err
	}

	// if err := tx.collection().FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
	// 	return err
	// }

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
