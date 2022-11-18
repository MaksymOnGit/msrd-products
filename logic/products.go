package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"msrd-products/db"
	"msrd-products/models"
	"time"
)

type ProductsRepository interface {
	Insert(product models.CreateProductRequest) (*models.Product, error)
	FindById(id string) (*models.Product, error)
	Update(product models.UpdateProductRequest) (*models.Product, error)
	UpdateByEvent(product models.UpdateProductEvent) (*models.Product, error)
	SoftDeleteById(id string) error
	SoftBatchDeleteById(ids []string) error
	QueryProducts(request models.QueryRequest) (error, models.QueryResponse[models.Product])
}

type productRepository struct {
	collection *mongo.Collection
	context    context.Context
}

func NewProductsRepository(context context.Context, dbContext db.DbContext) ProductsRepository {
	return &productRepository{dbContext.GetProductsCollection(), context}
}

func (r productRepository) Insert(product models.CreateProductRequest) (newProduct *models.Product, err error) {
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	res, err := r.collection.InsertOne(r.context, product)

	if err != nil {
		log.Println(err)
		return
	}

	query := bson.M{"_id": res.InsertedID}

	err = r.collection.FindOne(r.context, query).Decode(&newProduct)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r productRepository) Update(product models.UpdateProductRequest) (newProduct *models.Product, err error) {
	product.UpdatedAt = time.Now()

	_, err = r.collection.UpdateOne(r.context, bson.M{"_id": product.Id}, bson.M{"$set": product})

	if err != nil {
		log.Println(err)
		return
	}

	err = r.collection.FindOne(r.context, bson.M{"_id": product.Id}).Decode(&newProduct)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r productRepository) UpdateByEvent(product models.UpdateProductEvent) (newProduct *models.Product, err error) {
	product.UpdatedAt = time.Now()

	_, err = r.collection.UpdateOne(r.context, bson.M{"_id": product.Id}, bson.M{"$set": product})

	if err != nil {
		log.Println(err)
		return
	}

	err = r.collection.FindOne(r.context, bson.M{"_id": product.Id}).Decode(&newProduct)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r productRepository) FindById(id string) (product *models.Product, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		return
	}

	err = r.collection.FindOne(r.context, bson.M{"_id": oid}).Decode(&product)

	if err != nil {
		log.Println(err)
		return
	}

	return product, nil
}

func (r productRepository) SoftDeleteById(id string) (err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err = r.collection.UpdateOne(r.context, bson.M{"_id": oid}, bson.M{"$set": bson.M{"updated_at": time.Now(), "deleted": true}})

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r productRepository) SoftBatchDeleteById(ids []string) (err error) {

	oids := make([]primitive.ObjectID, len(ids))
	for i := range ids {
		oids[i], err = primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			log.Println(err)
			return
		}
	}

	_, err = r.collection.UpdateMany(r.context, bson.M{"_id": bson.M{"$in": oids}}, bson.M{"$set": bson.M{"updated_at": time.Now(), "deleted": true}})

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (r productRepository) QueryProducts(request models.QueryRequest) (err error, response models.QueryResponse[models.Product]) {

	var opts options.FindOptions
	opts.
		SetSkip(request.Offset).
		SetLimit(request.Rows)
	if request.SortField != "" && request.SortOrder != 0 {
		opts.SetSort(bson.D{{request.SortField, request.SortOrder}})
	}

	curs, err := r.collection.Find(r.context, bson.M{"deleted": nil}, &opts)

	var product models.Product
	for curs.Next(r.context) {
		err = curs.Decode(&product)
		if err != nil {
			log.Println(err)
			return
		}
		response.Result = append(response.Result, product)
	}

	totalRecCount, err := r.collection.EstimatedDocumentCount(r.context)
	if err != nil {
		log.Println(err)
		return
	}

	response.TotalRecordsCount = totalRecCount
	response.TotalPagesCount = totalRecCount / request.Rows
	currentRequestedPage := request.Offset / request.Rows
	if currentRequestedPage < response.TotalPagesCount {
		response.Page = currentRequestedPage + 1
	} else {
		response.Page = response.TotalPagesCount + 1
	}
	response.RecordsPerPageCount = request.Rows
	response.IsPrev = response.Page > 1
	response.IsNext = response.Page < response.TotalPagesCount

	return
}
