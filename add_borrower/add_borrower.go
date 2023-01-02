package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	uri := "mongodb://" + user + ":" + password + "@" + "localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	insertBorrower(client)
}

// Borrowerの配列が追加されていない本にBorrowerの空配列を追加する
func insertBorrower(client *mongo.Client) {
	col := client.Database("library-app").Collection("books")

	var books []map[string]interface{}
	cursor, err := col.Find(context.TODO(), bson.D{{Key: "borrower", Value: bson.D{{Key: "$exists", Value: false}}}})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		books = append(books, book)

		filter := bson.D{{Key: "id", Value: book["id"]}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "borrower", Value: bson.A{}}}}}

		_, err := col.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			panic(err)
		}

	}
	fmt.Println("Documents updated.")
}
