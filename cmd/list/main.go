package main

import (
	"fmt"

	books "github.com/Decedis/bookstore"
)

func main() {
	fmt.Println("Books in stock: ...")
	catalog, err := books.OpenCatalog("testdata/catalog.json")

	if err != nil {
		fmt.Errorf("Could not open catalog: %v\n", err)
		return
	}

	for _, book := range catalog.GetAllBooks() {
		fmt.Println(book)
	}
}
