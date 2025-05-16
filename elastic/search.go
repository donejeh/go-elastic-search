package elastic

import (
	"bytes"
	"context"
	"encoding/json"
)

func SearchProducts(query string) ([]map[string]interface{}, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
			},
		},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(searchQuery)

	res, err := ES.Search(
		ES.Search.WithContext(context.Background()),
		ES.Search.WithIndex("products"),
		ES.Search.WithBody(&buf),
		ES.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	json.NewDecoder(res.Body).Decode(&r)

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []map[string]interface{}
	for _, hit := range hits {
		results = append(results, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	return results, nil
}

func SemanticSearch(vector []float32) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "embedding",
			"query_vector":   vector,
			"k":              10,
			"num_candidates": 100,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := ES.Search(
		ES.Search.WithContext(context.Background()),
		ES.Search.WithIndex("products"),
		ES.Search.WithBody(&buf),
		ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []map[string]interface{}
	for _, hit := range hits {
		results = append(results, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	return results, nil
}
