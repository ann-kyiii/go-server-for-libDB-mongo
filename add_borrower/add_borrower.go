package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
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

// すべての本にBorrowerの空配列を追加する
func insertBorrower(client *mongo.Client) {
	col := client.Database("library-app").Collection("books")

	var books []map[string]interface{}
	cursor, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		books = append(books, book)

		filter := bson.D{{"id", book["id"]}}
		update := bson.D{{"$set", bson.D{{"borrower", bson.A{}}}}}

		_, err := col.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			panic(err)
		}

	}
	fmt.Printf("Documents updated.")
}
