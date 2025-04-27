package probes

import (
	"fmt"
)

type semanticSimilarityTask struct{}

func (t *semanticSimilarityTask) Name() string {
	return "Semantic Similarity Task"
}

func (t *semanticSimilarityTask) MetricName() string {
	return "Semantic Similarity"
}

func (t *semanticSimilarityTask) Run(models []string) (map[string]TaskResult, error) {
	pairs := []struct {
		original  string
		modified  string
	}{
		{
			original: "В лесу родилась ёлочка",
			modified: "В лесу выросла ёлочка",
		},
		{
			original: "В лесу она росла",
			modified: "В лесу она подрастала",
		},
		{
			original: "Зимой и летом стройная",
			modified: "Зимой и летом изящная",
		},
		{
			original: "Зеленая была",
			modified: "Изумрудная была",
		},
	}

	results := make(map[string]TaskResult)

	for _, model := range models {
		fmt.Printf("Model: %s\n", model)
		var totalSimilarity float64
		count := 0

		for i, pair := range pairs {
			originalEmb, err := getEmbedding(model, pair.original)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for original phrase %d: %v", i+1, err)
			}

			modifiedEmb, err := getEmbedding(model, pair.modified)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for modified phrase %d: %v", i+1, err)
			}

			sim, err := cosineSimilarity(originalEmb, modifiedEmb)
			if err != nil {
				return nil, fmt.Errorf("error computing similarity for pair %d: %v", i+1, err)
			}

			fmt.Printf("Cosine similarity for pair %d (%s vs. %s): %.4f\n", i+1, pair.original, pair.modified, sim)
			totalSimilarity += sim
			count++
		}

		avgSimilarity := totalSimilarity / float64(count)
		fmt.Printf("Average similarity: %.4f\n", avgSimilarity)
		results[model] = TaskResult{Metric: avgSimilarity}
	}

	maxSimilarity := -1.0
	var winner string
	for model, result := range results {
		if result.Metric > maxSimilarity {
			maxSimilarity = result.Metric
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
	fmt.Println("Registering Semantic Similarity Task")
	RegisterTask(&semanticSimilarityTask{})
}
