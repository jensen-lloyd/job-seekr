package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	mongoURI        = "mongodb://mongodb:27017"
	mongoDatabase   = "job_hunting"
	mongoCollection = "jobs"
)

var mongoClient *mongo.Client
var jobsCollection *mongo.Collection

func connectMongo() error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI),
	)

	if err != nil {
		return err
	}

	// Actually test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	mongoClient = client
	jobsCollection = client.Database(mongoDatabase).Collection(mongoCollection)

	log.Println("Connected to MongoDB")

	return nil
}








func initialiseMongo() error {

	if jobsCollection == nil {
		return fmt.Errorf("MongoDB is not connected")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Make Job.ID unique
	_, err := jobsCollection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "id", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	)

	if err != nil {
		return err
	}

	log.Println("MongoDB initialised")

	return nil
}





func jobExists(jobID string) (bool, error) {

	if jobsCollection == nil {
		return false, fmt.Errorf("MongoDB is not connected")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	filter := bson.M{
		"id": jobID,
	}

	count, err := jobsCollection.CountDocuments(
		ctx,
		filter,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
