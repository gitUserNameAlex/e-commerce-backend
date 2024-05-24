package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Database

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://admin:password@mongo:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB1")

	collection = client.Database("shop")
}

func GetProductsByQuery(params map[string]string) ([]bson.M, error) {
	var results []bson.M

	filter := bson.M{}
	for key, value := range params {
		filter[key] = bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}
	}

	cur, err := collection.Collection("products").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}

func GetAllProducts() ([]bson.M, error) {
	var results []bson.M

	cur, err := collection.Collection("products").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}

func GetProductsByCategory(categoryID int, limit int) ([]bson.M, error) {
	var results []bson.M

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	cur, err := collection.Collection("products").Find(context.Background(), bson.M{"category.id": categoryID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}

func GetProductById(id string) (bson.M, error) {
	var result bson.M

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.Collection("products").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
