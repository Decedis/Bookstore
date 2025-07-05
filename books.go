package books

import (
	"fmt"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Copies int
}

func BookToString(book Book) string {
	return fmt.Sprintf("%s by %s (copies: %v)", book.Title, book.Author, book.Copies)
}

var catalog = []Book{
	{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	},
	{
		ID:     "def",
		Title:  "White Heat",
		Author: "Dominic Sandbrook",
		Copies: 2,
	},
}

func GetAllBooks() []Book {
	return catalog
}

func GetBook(id string) (Book, bool) {
	for _, book := range catalog {
		if book.ID == id {
			return book, true
		}
	}
	return Book{}, false
}
