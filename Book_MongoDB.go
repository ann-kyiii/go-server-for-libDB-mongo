package main

type BookMongoDB struct {
	Id         int64       `bson:"id"`
	BookName   string      `bson:"bookName"`
	Genre      string      `bson:"genre"`
	SubGenre   string      `bson:"subGenre"`
	ISBN       string      `bson:"ISBN"`
	Find       int64       `bson:"find"`
	Sum        int64       `bson:"sum"`
	Author     string      `bson:"author"`
	Publisher  string      `bson:"publisher"`
	Pubdate    string      `bson:"pubdate"`
	Exist      string      `bson:"exist"`
	LocateAt4F bool        `bson:"locateAt4F"`
	WithDisc   string      `bson:"withDisc"`
	Other      string      `bson:"other"`
	Borrower   interface{} `bson:"borrower"`
	ImgURL     string      `bson:"imgURL"`
}
