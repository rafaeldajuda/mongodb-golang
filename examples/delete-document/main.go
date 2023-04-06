package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	Collection string
}

var mongoConfig MongoConfig
var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoConfig.User = os.Getenv("MONGO_USER")
	mongoConfig.Password = os.Getenv("MONGO_PASSWORD")
	mongoConfig.Host = os.Getenv("MONGO_HOST")
	mongoConfig.Port = os.Getenv("MONGO_PORT")
	mongoConfig.Database = os.Getenv("MONGO_DATABASE")
	mongoConfig.Collection = os.Getenv("MONGO_COLLECTION")
}

func main() {
	// mongodb uri format
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port)
	if uri == "" {
		log.Fatal("You must set your 'uri' variable.")
	}

	// connection
	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("MongoDB connection error: " + err.Error())
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("MongoDB disconnection error: " + err.Error())
		}
	}()

	// check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// set database and collection
	collection = client.Database(mongoConfig.Database).Collection(mongoConfig.Collection)

	// filter - get the document ID in your MongoDB
	id, err := primitive.ObjectIDFromHex("642f238625b383fba5ca34f0")
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{
		"_id": id,
	}

	// delete object
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	// print result
	fmt.Println(result.DeletedCount)

}
