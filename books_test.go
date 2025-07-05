package books_test

import (
	"testing"

	"slices"

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
	got := books.BookToString(input)
	if want != got {
		//t.Fatal("BookToString: unexpected result")
		t.Fatalf("BookToString TEST HAS FAILED:: ==> want %q, got %q", want, got)
	}
}

func TestGetAllBooks_ReturnsAllBooks(t *testing.T) {
	t.Parallel()
	want := []books.Book{
		{
			ID:     "abc",
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
		},
		{
			ID:     "def",
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
		},
	}
	got := books.GetAllBooks()
	if !slices.Equal(want, got) {
		t.Fatalf("GetAllBooks TEST HAS FAILED:: ==> want %#v got %#v", want, got)
	}
}

func TestGetBook_FindsBookInCatalogByID(t *testing.T) {
	t.Parallel()
	want := books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	}

	got, ok := books.GetBook("abc")
	if !ok {
		t.Fatal("book not found...")
	}
	if want != got {
		t.Fatalf("GetBook TEST HAS FAILED:: ==> want %v does not equal got %v", want, got)
	}
}

func TestGetBook_ReturnsFalseWHenBookNotFound(t *testing.T) {
	t.Parallel()
	_, ok := books.GetBook("nonexistent ID")
	if ok {
		t.Fatal("want false for nonexistent ID, got true")
	}
}
