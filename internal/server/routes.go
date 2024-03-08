package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.GET("/book/id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.getBook(w, r, nil)
	})
	r.GET("/list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.listBooks(w, r, nil)
	})
	r.GET("/authors", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.listAuthors(w, r, nil)
	})
	r.GET("/authors/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.getAuthor(w, r, ps)
	})
	r.GET("/books/author/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.getBooksByAuthorID(w, r, ps)
	})
	return r
}

type Book = struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	AuthorId    int      `json:"authorId"`
	Tags        []string `json:"tags"`
}

type Author struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
	Description string `json:"description"`
}

func (s *Server) getBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookID := ps.ByName("id")

	data, _, err := s.sb.From("books").
		Select("*", "1", false).
		Eq("id", bookID).
		Single().
		Execute()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var book Book
	if err := json.Unmarshal(data, &book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
func (s *Server) getBooksByAuthorID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authorID := ps.ByName("id")
	log.Printf("Author ID: %v", authorID)

	if authorID == "" {
		http.Error(w, "Missing author ID", http.StatusBadRequest)
		return
	}

	data, _, err := s.sb.From("books").
		Select("*", "exact", false).
		Eq("authorId", authorID).
		Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Books: %v", books)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
func (s *Server) listAuthors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, _, err := s.sb.From("authors").Select("*", "exact", false).Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var authors []Author
	if err := json.Unmarshal(data, &authors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

func (s *Server) getAuthor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authorID := ps.ByName("id")
	log.Printf("Author ID: %v", authorID)
	data, _, err := s.sb.From("authors").
		Select("*", "exact", false).
		Eq("id", authorID).
		Single().
		Execute()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var author Author
	if err := json.Unmarshal(data, &author); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}
func (s *Server) listBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, _, err := s.sb.From("books").Select("*", "exact", false).Execute()

	if err != nil {
		log.Printf("Error querying books: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		log.Printf("Error unmarshaling books: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Books: %+v", books)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
