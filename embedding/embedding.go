package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func GetEmbedding(text string) []float32 {
	_ = godotenv.Load(".env")
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.AdaEmbeddingV2,
		},
	)
	if err != nil {
		log.Fatalf("embedding error: %v", err)
	}

	return resp.Data[0].Embedding
}

func GetEmbeddingV2(query string) []float32 {
	payload := map[string]interface{}{
		"input": query,
		"model": "text-embedding-3-small",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	return result.Data[0].Embedding
}
