package books_test

import (
	"cmp"
	"slices"
	"testing"

	books "github.com/Decedis/bookstore"
)

// Remember: testing files occur outside the main function
// We can test within main, but that is when the function returns
// it's values, and nil.
// It's a different, and complimentary mechanism.
func TestBookToString_FormatsBookInfoAsString(t *testing.T) {
	t.Parallel()

	input := books.Book{
		Title:  "Sea Room",
		Author: "Adam Nicolson",
		Copies: 2,
	}
	want := "Sea Room by Adam Nicolson (copies: 2)"
	got := input.String()
	if want != got {
		// t.Fatal("BookToString: unexpected result")
		t.Fatalf("BookToString TEST HAS FAILED:: ==> want %q, got %q", want, got)
	}
}

func TestGetAllBooks_ReturnsAllBooks(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	want := []books.Book{
		{
			ID:     "abc",
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
		},
		{
			ID:     "xyz",
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
		},
	}
	got := catalog.GetAllBooks()
	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.Author, b.Author)
	})
	if !slices.Equal(want, got) {
		t.Fatalf("GetAllBooks TEST HAS FAILED:: ==> want %#v got %#v", want, got)
	}
}

func TestOpenCatalog_LoadsCatalogDataFromFile(t *testing.T) {
	t.Parallel()
	catalog, err := books.OpenCatalog("testdata/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	want := []books.Book{
		{
			ID:     "abc",
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
		},
		{
			ID:     "xyz",
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
		},
	}
	got := catalog.GetAllBooks()
	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.Author, b.Author)
	})
	if !slices.Equal(want, got) {
		t.Fatalf("GetAllBooks TEST HAS FAILED:: ==> want %#v got %#v", want, got)
	}
}

func TestGetBook_FindsBookInCatalogByID(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	want := books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	}

	got, ok := catalog.GetBook("abc")
	if !ok {
		t.Fatal("book not found...")
	}
	if want != got {
		t.Fatalf("GetBook TEST HAS FAILED:: ==> want %v does not equal got %v", want, got)
	}
}

func TestGetBook_ReturnsFalseWhenBookNotFound(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	_, ok := catalog.GetBook("nonexistent ID")
	if ok {
		t.Fatal("want false for nonexistent ID, got true")
	}
}

// // Pretty sure this test is redundant with the next test below it.
// func TestAddBook(t *testing.T) {
// 	t.Parallel()
// 	catalog := getTestCatalog()
// 	catalog.AddBook(books.Book{
// 		ID:     "123",
// 		Title:  "The Prize of all the Oceans",
// 		Author: "Glyn Williams",
// 		Copies: 2,
// 	})
// 	_, ok := catalog.GetBook("123")
// 	if !ok {
// 		t.Fatal("Added book not found")
// 	}
// }
//

func TestAddBook_AddsGivenBookToCatalog(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	_, ok := catalog.GetBook("123")
	if ok {
		t.Fatal("book already present")
	}
	catalog.AddBook(books.Book{
		ID:     "123",
		Title:  "The Prize of all the Oceans",
		Author: "Glyn Williams",
		Copies: 2,
	})
	_, ok = catalog.GetBook("123")
	if !ok {
		t.Fatal("added book not found")
	}
}

func TestSetCopies_SetsNumberOfCopiesToGivenValue(t *testing.T) {
	t.Parallel()
	book := books.Book{
		Copies: 5,
	}
	err := book.SetCopies(12)
	if err != nil {
		t.Fatal(err)
	}
	if book.Copies != 12 {
		t.Errorf("want 12 copies, got %d", book.Copies)
	}
}

func TestSetCopies_ReturnsErrorIfCopiesNegative(t *testing.T){
	t.Parallel()
	book := books.Book{}
	err := book.SetCopies(-1)
	if err == nil {
		t.Error("want error for negative copies, got nil")
		// We want the error to exist
		// If this block triggers, it means it doesn't exist.
		// If the error doesn't exist, despite the invalid input,
		// then that means the validation is failing,
		// thus, this test fails.
		// If the error exists, then this test passes, the method works.
	}
}

func getTestCatalog() books.Catalog {
	return books.Catalog{
		"abc": {
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
			ID:     "abc",
		},
		"xyz": {
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
			ID:     "xyz",
		},
	}
}
