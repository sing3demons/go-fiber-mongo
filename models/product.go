package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Desc       string             `bson:"desc"`
	Price      int                `bson:"price"`
	Image      string             `bson:"image"`
	CategoryID primitive.ObjectID `bson:"categoryId,omitempty"`
	Category   Category
}
