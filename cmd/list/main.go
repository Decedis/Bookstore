package main

import (
	books "github.com/Decedis/bookstore"

	"fmt"
)

func main() {
	fmt.Println("Books in stock: ...")
	allBooks := books.GetAllBooks()

	for _, book := range allBooks {
		fmt.Println(books.BookToString(book))
	}
}
