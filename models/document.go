package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id     primitive.ObjectID `bson:"_id" validate:"required"`
	Status string             `bson:"status"`
}

type UpdateDocumentEvent struct {
	Id     primitive.ObjectID `bson:"_id" validate:"required"`
	Status string             `bson:"status"`
}
