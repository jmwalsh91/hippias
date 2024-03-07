package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.HelloWorldHandler(w, r)
	})
	r.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.healthHandler(w, r)
	})
	r.GET("/book/id", s.getBook)
	r.GET("/list", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.list(w, r, nil)
	})
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

type Book = struct {
	id          int
	title       string
	author      string
	description string
	authorId    int
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

func (s *Server) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, _, err := s.sb.From("Books").Select("*", "exact", false).Execute()

	if err != nil {
		log.Printf("Error querying books: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(data)
	log.Printf("Data: %v", jsonData, err)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
