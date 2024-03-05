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
		s.HelloWorldHandler(w, r)
	})
	r.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.healthHandler(w, r)
	})
	r.GET("/item/:name", s.getItem)
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

func (s *Server) getItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")

	resp := map[string]string{
		"title":            name,
		"Author":           "Jean Baudrillard",
		"Publication Year": "1981",
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error handling JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}
