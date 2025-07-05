package main

import (
	"fmt"

	books "github.com/Decedis/bookstore"
)

func main() {
	book, ok := books.GetBook("xyz")
	if !ok {
		fmt.Println("Sorry, I couldn't find that book in the catalog.")
		return
	}
	fmt.Println(books.BookToString(book))
}
