package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Desc       string             `json:"desc" bson:"desc"`
	Price      int                `json:"price" bson:"price"`
	Image      string             `json:"image" bson:"image"`
	CategoryID primitive.ObjectID `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
}
