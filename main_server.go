package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	session "github.com/ipfans/echo-session"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}

	e := echo.New()

	store := session.NewCookieStore([]byte("secret-key"))
	store.MaxAge(2)
	e.Use(session.Sessions("ESESSION", store))
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("Request: %v\n", string(reqBody))
	}))
	initRouting(e)

	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/api/v1/bookId/:id", getBookWithID)
	e.POST("/api/v1/search", searchBooks)
	e.POST("/api/v1/searchGenre", searchGenre)
	e.POST("/api/v1/searchSubGenre", searchSubGenre)
	e.POST("/api/v1/borrow", borrowBook)
	e.POST("/api/v1/return", returnBook)
}

func hello(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}

func getBookWithID(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return errors.Wrapf(err, "errors when book id convert to int: %s", bookId)
	}

	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	var book BookMongoDB
	col.FindOne(context.Background(), bson.M{"id": bookId}).Decode(&book)

	return c.JSON(http.StatusOK, book)
}

func searchBooks(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	keywords := m["keywords"].([]interface{})
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	isAndSearch := m["isAndSearch"].(bool)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Println("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Println("【Error】", err)
		panic(err)
	}

	// get DB data
	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	var bookvalues1 BookValues
	filter := bson.D{{Key: "exist", Value: "〇"}}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// end find

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		cursor.Decode(&result)

		var book BookValue
		delete(result, "_id")
		book.Book = result
		bookvalues1 = append(bookvalues1, book)
	}

	var bookvalues2 BookValues
	filter = bson.D{{Key: "exist", Value: "一部発見"}}
	cursor, err = col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// end find

	var results2 []map[string]interface{}
	if err = cursor.All(context.TODO(), &results2); err != nil {
		panic(err)
	}

	for _, result := range results2 {
		var book BookValue
		cursor.Decode(&result)
		delete(result, "_id")
		book.Book = result
		bookvalues2 = append(bookvalues2, book)
	}

	bookvalues := append(bookvalues1, bookvalues2...)

	// search
	searchAttribute := []string{"publisher", "author", "bookName", "pubdate", "ISBN"}
	var data map[string]interface{}
	if isAndSearch {
		data = searchAnd(bookvalues, keywords, searchAttribute, offset, limit)
	} else {
		data = searchOr(bookvalues, keywords, searchAttribute, offset, limit)
	}
	return c.JSON(http.StatusOK, data)
}

func searchGenre(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	genre := m["genre"].(string)
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Println("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Println("【Error】", err)
		panic(err)
	}

	// get DB data
	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	var books []map[string]interface{}
	filter := bson.D{{Key: "exist", Value: "〇"}, {Key: "genre", Value: genre}}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		delete(book, "_id")
		books = append(books, book)
	}

	filter = bson.D{{Key: "exist", Value: "一部発見"}, {Key: "genre", Value: genre}}
	cursor, err = col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		delete(book, "_id")
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":     empty_list,
			"max_books": 0,
		}
		return c.JSON(http.StatusOK, data)
	} else {
		first := offset
		var last int
		if len(books) < offset+limit {
			last = len(books)
		} else {
			last = offset + limit
		}
		data := map[string]interface{}{
			"books":     books[first:last],
			"max_books": len(books),
		}
		return c.JSON(http.StatusOK, data)
	}
}

func searchSubGenre(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	subGenre := m["subGenre"].(string)
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Println("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Println("【Error】", err)
		panic(err)
	}

	// get DB data
	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	var books []map[string]interface{}
	filter := bson.D{{Key: "exist", Value: "〇"}, {Key: "subGenre", Value: subGenre}}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		delete(book, "_id")
		books = append(books, book)
	}

	filter = bson.D{{Key: "exist", Value: "一部発見"}, {Key: "subGenre", Value: subGenre}}
	cursor, err = col.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &books); err != nil {
		panic(err)
	}

	for _, book := range books {
		cursor.Decode(&book)
		delete(book, "_id")
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":     empty_list,
			"max_books": 0,
		}
		return c.JSON(http.StatusOK, data)
	} else {
		first := offset
		var last int
		if len(books) < offset+limit {
			last = len(books)
		} else {
			last = offset + limit
		}
		data := map[string]interface{}{
			"books":     books[first:last],
			"max_books": len(books),
		}
		return c.JSON(http.StatusOK, data)
	}
}

func borrowBook(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	t_id := m["id"].(string)
	name := m["name"].(string)
	id, err := strconv.Atoi(t_id)
	if err != nil {
		log.Println("【Error】", err)
		panic(err)
	}

	// get DB data
	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	// 既に登録されている名前では登録しない
	_, err = col.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$addToSet": bson.M{"borrower": name}})
	if err != nil {
		panic(err)
	}

	var book BookMongoDB
	col.FindOne(context.Background(), bson.M{"id": id}).Decode(&book)

	return c.JSON(http.StatusOK, book)
}

func returnBook(c echo.Context) error {
	session := session.Default(c)
	session.Set("AccessServer", "completed")
	session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	t_id := m["id"].(string)
	name := m["name"].(string)
	id, err := strconv.Atoi(t_id)
	if err != nil {
		log.Println("【Error】", err)
		panic(err)
	}

	// get DB data
	client := ConnectDB()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	COLLECTION_NAME := os.Getenv("COLLECTION_NAME")
	col := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	_, err = col.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$pull": bson.M{"borrower": name}})
	if err != nil {
		panic(err)
	}

	var book BookMongoDB
	col.FindOne(context.Background(), bson.M{"id": id}).Decode(&book)

	return c.JSON(http.StatusOK, book)
}
