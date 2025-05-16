package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/donejeh/go-elastic-search/embedding"
)

func CreateProductIndex() {
	mapping := `
	{
	"mappings": {
		"properties": {
		"name":        { "type": "text" },
		"description": { "type": "text" },
		"tags":        { "type": "keyword" },
		"popularity":  { "type": "integer" },
		"embedding":   { "type": "dense_vector", "dims": 1536, "index": true, "similarity": "cosine" }
		}
	}
	}`

	res, err := ES.Indices.Create("products", ES.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("Index created:", res.Status())
}

func BulkInsertProducts() {
	data, err := os.ReadFile("data/products.json")

	if err != nil {
		panic(err)
	}
	var products []map[string]interface{}
	err = json.Unmarshal(data, &products)
	if err != nil {
		panic(err)
	}

	for _, doc := range products {
		text := doc["name"].(string) + " " + doc["description"].(string)
		doc["embedding"], err = embedding.GetEmbedding(text)

		docJSON, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		res, err := ES.Index("products", bytes.NewReader(docJSON))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
	}

	fmt.Println("Indexed sample products.")
}
