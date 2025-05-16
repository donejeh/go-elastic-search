package api

import (
	"context"
	"encoding/json"
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

	// Get embedding vector from OpenAI
	embeddingVec := embedding.GetEmbedding(query)

	// Build base semantic search body
	searchBody := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "embedding",
			"query_vector":   embeddingVec,
			"k":              10,
			"num_candidates": 100,
		},
	}

	// Optional filter
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

	// Optional sorting
	if sortBy == "popularity" {
		searchBody["sort"] = []interface{}{
			map[string]interface{}{
				"popularity": map[string]string{"order": "desc"},
			},
		}
	}

	// Encode search body
	var b strings.Builder
	if err := json.NewEncoder(&b).Encode(searchBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute search
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

	// Decode response
	var rBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&rBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rBody)
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
