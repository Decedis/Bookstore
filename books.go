// Package books provides types and methods to generate and manipulate Book and
// Catalog data.
package books

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"sync"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Copies int
}

func (catalog *Catalog) SetCopies(ID string, copies int) error {
	book, ok := catalog.data[ID]
	if !ok {
		return fmt.Errorf("ID %q not found", ID)
	}
	err := book.SetCopies(copies)
	if err != nil {
		return err
	}
	catalog.data[ID] = book
	return nil
}

func (catalog *Catalog) GetCopies(ID string) (int, error) {
	book, ok := catalog.data[ID]
	if !ok {
		return 0, fmt.Errorf("ID %q not found", ID)
	}
	return book.Copies, nil
}

func (catalog *Catalog) Sync() error {
	catalog.mu.RLock() // TODO what does this do?
	defer catalog.mu.RUnlock()
	file, err := os.Create(catalog.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(catalog.data)
	if err != nil {
		return err
	}
	return nil
}

func (book Book) String() string {
	return fmt.Sprintf("%s by %s (copies: %v)", book.Title, book.Author, book.Copies)
}

func NewCatalog() *Catalog {
	return &Catalog{
		mu:   &sync.RWMutex{},
		data: map[string]Book{},
	}
}

func OpenCatalog(path string) (*Catalog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // defer means when the function ends
	catalog := Catalog{
		mu:   &sync.RWMutex{},
		data: map[string]Book{},
	}
	err = json.NewDecoder(file).Decode(&catalog.data)
	if err != nil {
		return nil, err
	}
	catalog.Path = path // remember where you came from
	return &catalog, nil
}

// SetCopies - A Methods on the Book type to set individual copy values on
// instances of the Book struct.
func (book *Book) SetCopies(copies int) error {
	if copies < 0 {
		return fmt.Errorf("cannot set copies to negative number: %d", copies)
	}
	book.Copies = copies
	return nil
}

// Catalog is a composite type of a map of Books, with strings to directly reference the Book instances within the map.
type Catalog struct {
	mu   *sync.RWMutex
	data map[string]Book
	Path string
}

func (catalog *Catalog) GetAllBooks() []Book {
	return slices.Collect(maps.Values(catalog.data)) // turns it into a slice
}

// GetBook method to fetch specific Book value via string ID. Returns either
// the Book value and/or boolean value for success.
func (catalog *Catalog) GetBook(ID string) (Book, bool) {
	// mu := catalog.mu
	// mu.Lock()
	// defer mu.Unlock()
	book, ok := catalog.data[ID]
	return book, ok
}

// AddBook method to add a value of type "Book" to composite type of Catalog.
func (catalog *Catalog) AddBook(book Book) error {
	_, ok := catalog.data[book.ID]
	if ok {
		return fmt.Errorf("Book already exists: %q", book.ID)
	}
	catalog.data[book.ID] = book // works because book.ID matches the key in the map.
	return nil
}
