package main

import (
	"fmt"
	"os"

	books "github.com/Decedis/bookstore"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: find <BOOK ID>")
		return
	}
	catalog, err := books.OpenCatalog("testdata/catalog.json")
	if err != nil {
		fmt.Errorf("Opening catalog: %v\n", err)
		return
	}
	ID := os.Args[1]
	book, ok := catalog.GetBook(ID)
	if !ok {
		fmt.Println("Sorry, I couldn't find that book in the catalog.")
		return
	}
	fmt.Println(book)
}
