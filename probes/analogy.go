package probes

import (
	"fmt"
	"math"
)

type analogyTask struct{}

func (t *analogyTask) Name() string {
	return "Analogy Task"
}

func (t *analogyTask) MetricName() string {
	return "Euclidean Distance"
}

func (t *analogyTask) Run(models []string) (map[string]TaskResult, error) {
	terms := map[string]string{
		"p": "Paris",
		"f": "France",
		"e": "England",
		"l": "London",
	}

	results := make(map[string]TaskResult)

	for _, model := range models {
		fmt.Printf("Model: %s\n", model)
		embeddings := make(map[string][]float64)

		for key, term := range terms {
			emb, err := getEmbedding(model, term)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for %s: %v", term, err)
			}
			embeddings[key] = emb
		}

		if len(embeddings["p"]) != len(embeddings["f"]) || len(embeddings["f"]) != len(embeddings["e"]) || len(embeddings["e"]) != len(embeddings["l"]) {
			return nil, fmt.Errorf("embedding dimensions do not match")
		}

		result := make([]float64, len(embeddings["p"]))
		for i := 0; i < len(embeddings["p"]); i++ {
			result[i] = embeddings["p"][i] - embeddings["f"][i] + embeddings["e"][i]
		}

		distance := 0.0
		for i := 0; i < len(result); i++ {
			diff := result[i] - embeddings["l"][i]
			distance += diff * diff
		}
		distance = math.Sqrt(distance)

		fmt.Printf("Euclidean distance between (p - f + e) and l: %.4f\n", distance)
		results[model] = TaskResult{Metric: distance}
	}

	minDistance := math.MaxFloat64
	var winner string
	for model, result := range results {
		if result.Metric < minDistance {
			minDistance = result.Metric
			winner = model
		}
	}
	for model := range results {
		if model == winner {
			result := results[model]
			result.Winner = model
			results[model] = result
		}
	}

	return results, nil
}

func init() {
	fmt.Println("Registering Analogy Task")
	RegisterTask(&analogyTask{})
}
