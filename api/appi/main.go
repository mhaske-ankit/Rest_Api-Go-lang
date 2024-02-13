package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type Book struct {
	Id        string `json:"id"`
	Book_Name string `json:"book_name"`
	Author    string `json:"author"`
}

var books []Book

func main() {
	books = []Book{
		{Id: "1", Book_Name: "The India Story", Author: "Bimal Jalal"},
		{Id: "2", Book_Name: "Wealth of Nations", Author: "Adam Smith"},
		{Id: "3", Book_Name: "Malgudi day", Author: "R K Narayan"},
		{Id: "5", Book_Name: "jungle_book", Author: "akshay savant"},
		{Id: "6", Book_Name: "animal farm", Author: "George Orwell"},
		{Id: "7", Book_Name: "war and peace", Author: "leo toistoy"},
		{Id: "8", Book_Name: "politics", Author: "aristole"},
	}

	repo := &bookRepository{}
	h := &NoteHandler{
		Repository: repo,
	}
	router := initializeRoutes(h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe() // Run the HTTP server
}

func initializeRoutes(h *NoteHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/books", h.getAllBooks)
	mux.HandleFunc("/books/add", h.addBook)
	mux.HandleFunc("/books/update", h.updateBook)
	mux.HandleFunc("/books/delete", h.deleteBook)
	mux.HandleFunc("/books/", h.getBookByID)

	return mux
}

type NoteHandler struct {
	Repository BookRepository
}

func (h *NoteHandler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *NoteHandler) addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	books = append(books, newBook)
	w.WriteHeader(http.StatusCreated)
}

func (h *NoteHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, b := range books {
		if b.Id == id {
			books[i] = updatedBook
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func (h *NoteHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	for i, b := range books {
		if b.Id == id {
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func (h *NoteHandler) getBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")
	for _, b := range books {
		if b.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(b)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

type BookRepository interface {
	GetAll() ([]Book, error)
	GetById(id string) (Book, error)
	Add(Book) error
	Update(id string, updatedBook Book) error
	Delete(id string) error
}

type bookRepository struct{}

func (r *bookRepository) GetAll() ([]Book, error) {
	return books, nil
}

func (r *bookRepository) GetById(id string) (Book, error) {
	for _, b := range books {
		if b.Id == id {
			return b, nil
		}
	}
	return Book{}, errors.New("Book not found")
}

func (r *bookRepository) Add(book Book) error {
	books = append(books, book)
	return nil
}

func (r *bookRepository) Update(id string, updatedBook Book) error {
	for i, b := range books {
		if b.Id == id {
			books[i] = updatedBook
			return nil
		}
	}
	return errors.New("Book not found")
}

func (r *bookRepository) Delete(id string) error {
	for i, b := range books {
		if b.Id == id {
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}
	return errors.New("Book not found")
}
