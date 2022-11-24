package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"msrd-products/db"
	"msrd-products/models"
)

type DocumentsRepository interface {
	UpdateByEvent(document models.UpdateDocumentEvent) (*models.Document, error)
	FindById(id string) (*models.Document, error)
}

type documentsRepository struct {
	collection *mongo.Collection
	context    context.Context
}

func NewDocumentsRepository(context context.Context, dbContext db.DbContext) DocumentsRepository {
	return &documentsRepository{dbContext.GetDocumentsCollection(), context}
}

func (r documentsRepository) UpdateByEvent(document models.UpdateDocumentEvent) (newDocument *models.Document, err error) {

	_, err = r.collection.UpdateOne(r.context, bson.M{"_id": document.Id}, bson.M{"$set": bson.M{"status": document.Status}})

	if err != nil {
		log.Println(err)
		return
	}

	err = r.collection.FindOne(r.context, bson.M{"_id": document.Id}).Decode(&newDocument)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r documentsRepository) FindById(id string) (document *models.Document, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		return
	}

	err = r.collection.FindOne(r.context, bson.M{"_id": oid}).Decode(&document)

	if err != nil {
		log.Println(err)
		return
	}

	return document, nil
}
