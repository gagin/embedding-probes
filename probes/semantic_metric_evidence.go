package probes

import (
	"fmt"
)

type semanticMetricEvidenceTask struct{}

func (t *semanticMetricEvidenceTask) Name() string {
	return "Semantic Metric Evidence Task"
}

func (t *semanticMetricEvidenceTask) MetricName() string {
	return "Weighted Similarity"
}

func (t *semanticMetricEvidenceTask) Run(models []string) (map[string]TaskResult, error) {
	// Define the metric and evidence chunks with ground truth relevance
	type Evidence struct {
		text     string
		relevant bool // True for Class A or B, False for Class C
	}
	metric := "How comprehensive are the organization’s employee wellness programs?"
	evidenceChunks := []Evidence{
		{
			text:     "At Horizon Inc., we prioritize employee well-being with a holistic wellness program. This includes gym memberships, mental health counseling through an on-site therapist, and financial planning workshops held quarterly. Last year, we expanded with mindfulness sessions and a subsidized healthy meal plan, ensuring staff thrive in all aspects of life.",
			relevant: true, // Class A: comprehensive
		},
		{
			text:     "Our company values health and provides a fitness center at headquarters with an annual health fair. Employees enjoy discounted gym rates and a spring step challenge. While we focus on physical fitness, we’re exploring more offerings based on feedback from our yearly employee survey.",
			relevant: true, // Class B: moderate
		},
		{
			text:     "Horizon Inc. delivers top-quality products, with teams working diligently to meet deadlines. We recently upgraded our office with ergonomic furniture and a modern break room to enhance comfort during long hours, reflecting our commitment to a productive environment.",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "Our company focuses on innovation and productivity. We introduced a new project management tool to streamline workflows and ensure timely delivery of client projects. Team meetings are held weekly to align on goals and address challenges, fostering a collaborative environment.",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "At Horizon Inc., our engineering team works on cutting-edge projects. We’ve invested in state-of-the-art workstations and offer continuous training to keep technical skills up to date. The company cafeteria was renovated to include fast-food options.",
			relevant: false, // Class C: minimal/no wellness programs
		},
	}

	results := make(map[string]TaskResult)

	for _, model := range models {
		fmt.Printf("Model: %s\n", model)
		metricEmb, err := getEmbedding(model, metric)
		if err != nil {
			return nil, fmt.Errorf("error getting embedding for metric: %v", err)
		}

		var totalSimilarity float64
		var relevantCount, irrelevantCount int
		for i, evidence := range evidenceChunks {
			evidenceEmb, err := getEmbedding(model, evidence.text)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for evidence %d: %v", i+1, err)
			}

			sim, err := cosineSimilarity(metricEmb, evidenceEmb)
			if err != nil {
				return nil, fmt.Errorf("error computing similarity for evidence %d: %v", i+1, err)
			}

			fmt.Printf("Evidence %d: Similarity=%.4f, Relevant=%v\n", i+1, sim, evidence.relevant)
			if evidence.relevant {
				totalSimilarity += sim
				relevantCount++
			} else {
				totalSimilarity += (1.0 - sim)
				irrelevantCount++
			}
		}

		weightedSimilarity := totalSimilarity / float64(len(evidenceChunks))
		fmt.Printf("Weighted Average Similarity: %.4f\n", weightedSimilarity)
		results[model] = TaskResult{Metric: weightedSimilarity}
	}

	// Determine winner
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
	fmt.Println("Registering Semantic Metric Evidence Task")
	RegisterTask(&semanticMetricEvidenceTask{})
}
