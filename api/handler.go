package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/donejeh/go-elastic-search/elastic"
	"github.com/donejeh/go-elastic-search/embedding"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	filterTag := r.URL.Query().Get("tag")
	sortBy := r.URL.Query().Get("sort")

	embeddingVec, err := embedding.GetEmbedding(query)
	useSemanticSearch := err == nil && len(embeddingVec) > 0

	if !useSemanticSearch {
		log.Printf("Embedding failed or empty: %v — falling back to keyword search.", err)
	}

	var searchBody map[string]interface{}

	if useSemanticSearch {
		// Primary semantic search
		searchBody = map[string]interface{}{
			"knn": map[string]interface{}{
				"field":          "embedding",
				"query_vector":   embeddingVec,
				"k":              10,
				"num_candidates": 100,
			},
		}

		if filterTag != "" {
			searchBody["query"] = map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": []interface{}{
						map[string]interface{}{
							"term": map[string]interface{}{
								"tags.keyword": filterTag,
							},
						},
					},
				},
			}
		}

		if sortBy == "popularity" {
			searchBody["sort"] = []interface{}{
				map[string]interface{}{
					"popularity": map[string]string{"order": "desc"},
				},
			}
		}
	} else {
		// Fallback directly to keyword search
		searchBody = buildKeywordFallbackBody(query, filterTag, sortBy)
	}

	// Encode and execute initial search
	var b strings.Builder
	if err := json.NewEncoder(&b).Encode(searchBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := elastic.ES.Search(
		elastic.ES.Search.WithContext(context.Background()),
		elastic.ES.Search.WithIndex("products"),
		elastic.ES.Search.WithBody(strings.NewReader(b.String())),
		elastic.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var rBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&rBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fallback to keyword match if no hits
	if useSemanticSearch {
		hits := rBody["hits"].(map[string]interface{})["hits"].([]interface{})
		if len(hits) == 0 {
			searchBody = buildKeywordFallbackBody(query, filterTag, sortBy)

			var fallback strings.Builder
			if err := json.NewEncoder(&fallback).Encode(searchBody); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			res, err = elastic.ES.Search(
				elastic.ES.Search.WithContext(context.Background()),
				elastic.ES.Search.WithIndex("products"),
				elastic.ES.Search.WithBody(strings.NewReader(fallback.String())),
				elastic.ES.Search.WithTrackTotalHits(true),
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&rBody); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rBody)
}

// buildKeywordFallbackBody returns a traditional keyword search query
func buildKeywordFallbackBody(query string, tag string, sortBy string) map[string]interface{} {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"name", "description", "tags"},
					},
				},
				"filter": []interface{}{},
			},
		},
	}

	if tag != "" {
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					"tags.keyword": tag,
				},
			},
		)
	}

	if sortBy == "popularity" {
		body["sort"] = []interface{}{
			map[string]interface{}{
				"popularity": map[string]string{"order": "desc"},
			},
		}
	}

	return body
}

// func SearchHandler(w http.ResponseWriter, r *http.Request) {
// 	q := r.URL.Query().Get("q")
// 	results, err := elastic.SearchProducts(q)
// 	if err != nil {
// 		http.Error(w, "Search error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(results)
// }
