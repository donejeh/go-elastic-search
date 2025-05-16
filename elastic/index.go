package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func CreateProductIndex() {
	mapping := `
    {
        "mappings": {
            "properties": {
                "name": { "type": "text" },
                "description": { "type": "text" },
                "tags": { "type": "keyword" },
				"popularity": { "type": "integer" }
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
		panic(err) // Or use a better error-handling approach
	}
	var products []map[string]interface{}
	err = json.Unmarshal(data, &products)
	if err != nil {
		panic(err)
	}

	for _, doc := range products {
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
