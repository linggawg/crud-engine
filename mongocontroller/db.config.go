package mongocontroller

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDb *mongo.Client

func InitMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var err error
	mongoDb, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = mongoDb.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return mongoDb, nil
}

func DbMongo() *mongo.Database {
	return mongoDb.Database("dbexample")
}
