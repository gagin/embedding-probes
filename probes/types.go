package probes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type Task interface {
	Name() string
	MetricName() string
	Run(models []string) (map[string]TaskResult, error)
}

type TaskResult struct {
	Metric float64
	Winner string
}

var TaskRegistry []Task

func RegisterTask(task Task) {
	TaskRegistry = append(TaskRegistry, task)
}

func getEmbedding(model, text string) ([]float64, error) {
	url := "http://localhost:11434/api/embeddings"
	payload := map[string]string{
		"model":  model,
		"prompt": text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var result struct {
		Embedding []float64 `json:"embedding"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result.Embedding, nil
}

func cosineSimilarity(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors have different lengths: %d vs %d", len(a), len(b))
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	normA = math.Sqrt(normA)
	normB = math.Sqrt(normB)
	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("zero magnitude vector")
	}

	return dotProduct / (normA * normB), nil
}
