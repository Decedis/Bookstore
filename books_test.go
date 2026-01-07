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
	catalog.Path = t.TempDir() + "/catalog"
	err := catalog.Sync()
	if err != nil {
		t.Fatal(err)
	}
	newCatalog, err := books.OpenCatalog(catalog.Path)
	if err != nil {
		t.Fatal(err)
	}
	bookList := newCatalog.GetAllBooks()
	assertTestBooks(t, bookList)
}

func TestSetCopies_IsRaceFree(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	go func() {
		for range 100 {
			err := catalog.SetCopies("abc", 0)
			if err != nil {
				panic(err)
			}
		}
	}()
	for range 100 {
		_, err := catalog.GetCopies("abc")
		if err != nil {
			t.Fatal(err)
		}
	}
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
	catalog.Path = t.TempDir() + "/catalog"
	// err := catalog.Sync("testdata/catalog.new")
	err := catalog.Sync()
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
	err := catalog.AddBook(books.Book{
		ID:     "123",
		Title:  "The Prize of all the Oceans",
		Author: "Glyn Williams",
		Copies: 2,
	})
	if err != nil {
		t.Fatal("Added book not found")
	}
	_, ok = catalog.GetBook("123")
	if !ok {
		t.Fatal("added book not found")
	}
}

func TestAddBook_ReturnsErrorIfIDExists(t *testing.T) {
	t.Parallel()
	catalog := getTestCatalog()
	_, ok := catalog.GetBook("abc")
	if !ok {
		t.Fatal("Book doesn't exists within catalog")
	}
	err := catalog.AddBook(books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	})
	if err == nil {
		t.Fatal("Want error for duplicate ID, got nil")
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
	catalog := books.NewCatalog()
	err := catalog.AddBook(books.Book{
		ID:     "test",
		Copies: 5,
	})
	err = catalog.SetCopies("test", 12)
	if err != nil {
		t.Fatal(err)
	}
	book, ok := catalog.GetBook("test")
	if !ok {
		t.Fatal("test book not found in catalog")
	}
	if book.Copies != 12 {
		t.Errorf("want 12 copies, got %d", book.Copies)
	}
}

func TestSetCopies_ReturnsErrorIfCopiesNegative(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	err := catalog.AddBook(
		books.Book{
			ID:     "test",
			Copies: 1,
		})
	if err != nil {
		t.Error("Could not AddBook: ", err)
	}
	err = catalog.SetCopies("test", -1)
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

func TestNewCatalog_CreatesEmptyCatalog(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	books := catalog.GetAllBooks()
	if len(books) > 0 {
		t.Errorf("want empty catalog, got %#v", books)
	}
}

func getTestCatalog() *books.Catalog {
	catalog := books.NewCatalog()
	err := catalog.AddBook(books.Book{
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
		ID:     "abc",
	})
	if err != nil {
		panic(err)
	}
	err = catalog.AddBook(books.Book{
		Title:  "White Heat",
		Author: "Dominic Sandbrook",
		Copies: 2,
		ID:     "xyz",
	})
	if err != nil {
		panic(err)
	}
	return catalog
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
