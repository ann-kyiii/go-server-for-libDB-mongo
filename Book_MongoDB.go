package main

type BookMongoDB struct {
	Id         int64       `json:"id", bson:"id"`
	BookName   string      `json:"bookName", bson:"bookName"`
	Genre      string      `json:"genre", bson:"genre"`
	SubGenre   string      `json:"subGenre", bson:"subGenre"`
	ISBN       string      `json:"iSBN", bson:"ISBN"`
	Find       int64       `json:"find", bson:"find"`
	Sum        int64       `json:"sum", bson:"sum"`
	Author     string      `json:"author", bson:"author"`
	Publisher  string      `json:"publisher", bson:"publisher"`
	Pubdate    string      `json:"pubdate", bson:"pubdate"`
	Exist      string      `json:"exist", bson:"exist"`
	LocateAt4F bool        `json:"locateAt4F", bson:"locateAt4F"`
	WithDisc   string      `json:"withDisc", bson:"withDisc"`
	Other      string      `json:"other", bson:"other"`
	Borrower   interface{} `json:"borrower", bson:"borrower"`
	ImgURL     string      `json:"imgURL", bson:"imgURL"`
}
