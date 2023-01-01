package main

import (
	"regexp"
	"sort"
)

type BookValue struct {
	Book  map[string]interface{}
	value int
}

type BookValues []BookValue

func (u BookValues) Len() int {
	return len(u)
}

func (u BookValues) Less(i, j int) bool {
	return u[i].value > u[j].value
}

func (u BookValues) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func searchOr(bookvalues BookValues, keywords []interface{}, searchAttribute []string, offset int, limit int) map[string]interface{} {
	for i, book := range bookvalues {
		for _, word := range keywords {
			// 大文字小文字を区別しない検索
			r := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word.(string)))
			for v, att := range searchAttribute {
				if book.Book[att] == nil {
					continue
				}
				if book.Book[att] == "Genre" || book.Book[att] == "SubGenre" {
					continue
				}
				if r.MatchString(book.Book[att].(string)) {
					bookvalues[i].value += (v + 1)
					break
				}
			}
		}
	}

	return createResult(bookvalues, offset, limit)
}

func searchAnd(bookvalues BookValues, keywords []interface{}, searchAttribute []string, offset int, limit int) map[string]interface{} {
	for i, book := range bookvalues {
		for _, word := range keywords {
			// 大文字小文字を区別しない検索
			r := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word.(string)))
			isMatch := false
			for v, att := range searchAttribute {
				if book.Book[att] == nil {
					continue
				}
				if book.Book[att] == "Genre" || book.Book[att] == "SubGenre" {
					continue
				}
				if r.MatchString(book.Book[att].(string)) {
					bookvalues[i].value += (v + 1)
					isMatch = true
					break
				}
			}
			if !isMatch {
				bookvalues[i].value = 0
				break
			}
		}
	}

	return createResult(bookvalues, offset, limit)
}

func createResult(bookvalues BookValues, offset int, limit int) map[string]interface{} {
	sort.Sort(bookvalues)
	var books []map[string]interface{}
	for _, book := range bookvalues {
		if book.value != 0 {
			books = append(books, book.Book)
		}
	}
	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":     empty_list,
			"max_books": 0,
		}
		return data
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
		return data
	}
}
