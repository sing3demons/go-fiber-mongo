package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Desc       string             `bson:"desc" json:"desc"`
	Price      int                `bson:"price" json:"price"`
	Image      string             `bson:"image" json:"image"`
	CategoryID primitive.ObjectID `bson:"categoryId,omitempty" json:"categoryId,omitempty"`
	Category   Category           `bson:"category,omitempty" json:"category,omitempty"`
}
