package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type DbContext interface {
	Dispose()
	GetProductsCollection() *mongo.Collection
}

type connection struct {
	client           *mongo.Client
	database         *mongo.Database
	connectionConfig DbContextConfig
}

type DbContextConfig struct {
	ContextTimeout time.Duration
}

func BuildDbContext(connectionString string, database string, configSetup func(config *DbContextConfig)) (DbContext, error) {

	var err error
	var connection connection

	//init default values
	connection.connectionConfig.ContextTimeout = 30 * time.Second

	//option to override default values
	configSetup(&connection.connectionConfig)

	ctx, cancel := context.WithTimeout(context.Background(), connection.connectionConfig.ContextTimeout)
	defer cancel()

	connection.client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = connection.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	connection.database = connection.client.Database(database)

	return connection, nil
}

func (connection connection) GetProductsCollection() *mongo.Collection {
	collection := connection.database.Collection("products")
	return collection
}

func (connection connection) Dispose() {
	ctx, cancel := context.WithTimeout(context.Background(), connection.connectionConfig.ContextTimeout)
	defer cancel()
	connection.client.Disconnect(ctx)
}
