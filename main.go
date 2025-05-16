package main

import (
	"log"
	"net/http"

	"github.com/donejeh/go-elastic-search/api"
	"github.com/donejeh/go-elastic-search/elastic"

	"github.com/gorilla/mux"
)

func main() {
	elastic.Init()
	elastic.CreateProductIndex()
	elastic.BulkInsertProducts()

	r := mux.NewRouter()
	r.HandleFunc("/search", api.SearchHandler).Methods("GET")

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
