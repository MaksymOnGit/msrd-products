package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Quantity    *float32           `json:"quantity" bson:"quantity"`
}

type CreateProductRequest struct {
	Name        string    `json:"name" bson:"name" validate:"required"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"-" bson:"created_at"`
	UpdatedAt   time.Time `json:"-" bson:"updated_at"`
}

type UpdateProductRequest struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description"`
	UpdatedAt   time.Time          `json:"-" bson:"updated_at"`
	Quantity    *float32           `json:"-" bson:"quantity"`
}
