package main

import (
	"log"
	"net/http"

	"github.com/donejeh/go-elastic-search/api"
	"github.com/donejeh/go-elastic-search/elastic"

	"github.com/gorilla/mux"

	"github.com/donejeh/go-elastic-search/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	elastic.Init()
	metrics.Init()
	elastic.CreateProductIndex()
	elastic.BulkInsertProducts()

	r := mux.NewRouter()
	r.HandleFunc("/search", api.SearchHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
