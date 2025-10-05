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
// It's a different, and complimentary mechanism

func TestOpenCatalog_ReadsSameDataWrittenBySync(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	path := t.TempDir() + "/catalog"
	err := catalog.Sync(path)
	if err != nil {
		t.Fatal(err)
	}
	newCatalog, err := books.OpenCatalog(path)
	if err != nil {
		t.Fatal(err)
	}
	bookList := newCatalog.GetAllBooks()
	assertTestBooks(t, bookList)
}

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
	bookList := catalog.GetAllBooks()
	assertTestBooks(t, bookList)
}

// “a failing test that calls Sync on a catalog and then OpenCatalog on the resulting file.”
func TestSyncWritesCatalogDataToFile(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	err := catalog.Sync("testdata/catalog.new")
	if err != nil {
		t.Fatal(err)
	}
	newCatalog, err := books.OpenCatalog("testdata/catalog.new")
	if err != nil {
		t.Fatal(err)
	}
	bookList := newCatalog.GetAllBooks()
	assertTestBooks(t, bookList)
}

func TestOpenCatalog_LoadsCatalogDataFromFile(t *testing.T) {
	t.Parallel()
	catalog, err := books.OpenCatalog("testdata/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	bookList := catalog.GetAllBooks()
	assertTestBooks(t, bookList)
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

func TestSetCopies_OnCatalogModifiesSpecificBook(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	book, ok := catalog.GetBook("abc") // check if our book exists
	if !ok {
		t.Fatal("book not found")
	}
	if book.Copies != 1 { // we want it to be 1
		t.Fatalf("want 1 copy before change, got %d", book.Copies)
	}
	err := catalog.SetCopies("abc", 2) // failure point
	if err != nil {
		t.Fatal("error: ", err)
	}
	book, ok = catalog.GetBook("abc") // get the new book values...and error point
	if !ok {
		t.Fatal("book not found")
	}
	if book.Copies != 2 { // compare our values
		t.Fatalf("want 2 copies after change, got %d", book.Copies)
	}

}

func TestSetCopies_SetsNumberOfCopiesToGivenValue(t *testing.T) {
	t.Parallel()
	catalog := books.Catalog{"test": {
		Copies: 5,
	}}
	err := catalog.SetCopies("test", 12)
	if err != nil {
		t.Fatal(err)
	}
	if catalog["test"].Copies != 12 {
		t.Errorf("want 12 copies, got %d", catalog["test"].Copies)
	}
}

func TestSetCopies_ReturnsErrorIfCopiesNegative(t *testing.T) {
	t.Parallel()
	catalog := books.Catalog{"test": {}}
	err := catalog.SetCopies("test", -1)
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

func assertTestBooks(t *testing.T, got []books.Book) {
	t.Helper()
	want := []books.Book{
		{
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
			ID:     "abc",
		},
		{
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
			ID:     "xyz",
		},
	}
	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.Author, b.Author)
	})
	if !slices.Equal(want, got) {
		t.Fatalf("want %#v, got %#v", want, got)
	}
}
