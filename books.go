package books

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Copies int
}


func (catalog *Catalog) SetCopies(ID string, copies int) error {
    if copies < 0 {
		return fmt.Errorf("negative number of copies: %d", copies) // this is the error we can return => !nil
	}
	book, ok := (*catalog)[ID] // need to verify  that the book is present
	if !ok {
		return fmt.Errorf("Error finding book")
	}
	book.Copies = copies
	(*catalog)[ID] = book
	fmt.Println("Catalog after 'update' => ", catalog)

	return nil // this is an error we can return => nil
}

func (catalog *Catalog) Sync(file string) error {
	payload, err := json.Marshal(catalog)
	if err != nil {
		fmt.Printf("Error marshalling data: %v", err)
	}
	err = os.WriteFile(file, payload, 0o644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
	return nil
}

func (book Book) String() string {
	return fmt.Sprintf("%s by %s (copies: %v)", book.Title, book.Author, book.Copies)
}

func OpenCatalog(path string) (Catalog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // defer means when the function ends
	catalog := Catalog{}
	err = json.NewDecoder(file).Decode(&catalog)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

type Catalog map[string]Book

func (catalog Catalog) GetAllBooks() []Book {
	return slices.Collect(maps.Values(catalog)) // turns it into a slice
}

func (catalog Catalog) GetBook(ID string) (Book, bool) {
	book, ok := catalog[ID]
	return book, ok
}

func (catalog Catalog) AddBook(book Book) {
	catalog[book.ID] = book // works because book.ID matches the key in the map.
}
