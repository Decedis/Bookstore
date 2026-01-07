// package utilizing the Books package to find copies
package main

import (
	"fmt"
	"os"
	"strconv"

	books "github.com/Decedis/bookstore"
)

func main() {
	// 1 get the current catalog
	// 2 accept an input
	// 3 that input must be valid
	// 4 that input must correspond to a key in the catalog
	// 5 the changes made must be valid changes (cannot be letters)
	// 6 return
	if len(os.Args) != 3 {
		fmt.Println("Usage: update <BOOK ID> <NUMBER OF COPIES>")
		return
	}
	catalog, err := books.OpenCatalog("testdata/catalog.json")
	if err != nil {
		fmt.Printf("Trouble getting catalog: %v\n", err)
		return
	}
	ID := os.Args[1]
	copies, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = catalog.SetCopies(ID, copies)
	if err != nil {
		fmt.Printf("updating book: %v\n", err)
		return
	}
	err = catalog.Sync()
	if err != nil {
		fmt.Printf("writing catalog: %v\n", err)
		return
	}
	fmt.Printf("Updated book %v to %d copies\n", ID, copies)
}
