package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(cfg MongoConfig) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.Host)

	// Connect to MongoDB
	mongoDb, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = mongoDb.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return mongoDb.Database(cfg.Name), nil
}
