package probes

import (
	"fmt"
)

type crossLanguageTask struct{}

func (t *crossLanguageTask) Name() string {
	return "Cross-Language Capability Task"
}

func (t *crossLanguageTask) MetricName() string {
	return "Cross-Language Similarity"
}

func (t *crossLanguageTask) Run(models []string) (map[string]TaskResult, error) {
	pairs := []struct {
		russian string
		french  string
	}{
		{
			russian: "В лесу родилась ёлочка",
			french:  "Un sapin est né dans la forêt",
		},
		{
			russian: "В лесу она росла",
			french:  "Dans la forêt, il a grandi",
		},
		{
			russian: "Зимой и летом стройная",
			french:  "En hiver et en été, élancé",
		},
		{
			russian: "Зеленая была",
			french:  "Il était vert",
		},
	}

	results := make(map[string]TaskResult)

	for _, model := range models {
		fmt.Printf("Model: %s\n", model)
		var totalSimilarity float64
		count := 0

		for i, pair := range pairs {
			russianEmb, err := getEmbedding(model, pair.russian)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for Russian phrase %d: %v", i+1, err)
			}

			frenchEmb, err := getEmbedding(model, pair.french)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for French phrase %d: %v", i+1, err)
			}

			sim, err := cosineSimilarity(russianEmb, frenchEmb)
			if err != nil {
				return nil, fmt.Errorf("error computing similarity for pair %d: %v", i+1, err)
			}

			fmt.Printf("Cosine similarity for pair %d (%s vs. %s): %.4f\n", i+1, pair.russian, pair.french, sim)
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
	fmt.Println("Registering Cross-Language Capability Task")
	RegisterTask(&crossLanguageTask{})
}
