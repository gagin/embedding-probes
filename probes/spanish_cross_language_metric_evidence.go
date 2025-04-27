package probes

import (
	"fmt"
)

type spanishCrossLanguageMetricEvidenceTask struct{}

func (t *spanishCrossLanguageMetricEvidenceTask) Name() string {
	return "Spanish Cross-Language Metric Evidence Task"
}

func (t *spanishCrossLanguageMetricEvidenceTask) MetricName() string {
	return "Accuracy"
}

func (t *spanishCrossLanguageMetricEvidenceTask) Run(models []string) (map[string]TaskResult, error) {
	type Evidence struct {
		text     string
		lang     string
		relevant bool
	}
	metric := "How comprehensive are the organization’s employee wellness programs?"
	evidenceChunks := []Evidence{
		{
			text:     "At Horizon Inc., we prioritize employee well-being with a comprehensive program including gym memberships, mental health counseling, and financial planning workshops. We offer mindfulness sessions and a subsidized healthy meal plan to support all aspects of employee health.",
			lang:     "English",
			relevant: true, // Class A: comprehensive
		},
		{
			text:     "Nuestra empresa valora la salud de los empleados y ofrece un centro de fitness en la sede con una feria de salud anual. Los empleados disfrutan de tarifas de gimnasio con descuento y un desafío de pasos en primavera. Nos enfocamos en la aptitud física, pero estamos explorando más opciones según los comentarios de la encuesta anual.",
			lang:     "Spanish",
			relevant: true, // Class B: moderate
		},
		{
			text:     "Horizon Inc. se dedica a entregar productos de alta calidad, con equipos que trabajan arduamente para cumplir plazos. Recientemente actualizamos nuestra oficina con muebles ergonómicos y una sala de descanso moderna para mejorar la comodidad durante largas horas de trabajo.",
			lang:     "Spanish",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "Our company focuses on innovation and productivity. We recently introduced a new project management tool to streamline workflows and ensure timely delivery of client projects. Team meetings are held weekly to align on goals and address challenges, fostering a collaborative environment.",
			lang:     "English",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "En Horizon Inc., nuestro equipo de ingenieros trabaja en proyectos de vanguardia. Hemos invertido en estaciones de trabajo de última generación y ofrecemos formación continua para mantener las habilidades técnicas al día. La cafetería de la empresa se renovó para incluir opciones de comida rápida.",
			lang:     "Spanish",
			relevant: false, // Class C: minimal/no wellness programs
		},
	}

	results := make(map[string]TaskResult)

	for _, model := range models {
		fmt.Printf("Model: %s\n", model)
		correct := 0
		metricEmb, err := getEmbedding(model, metric)
		if err != nil {
			return nil, fmt.Errorf("error getting embedding for metric: %v", err)
		}

		for i, evidence := range evidenceChunks {
			evidenceEmb, err := getEmbedding(model, evidence.text)
			if err != nil {
				return nil, fmt.Errorf("error getting embedding for evidence %d (%s): %v", i+1, evidence.lang, err)
			}

			sim, err := cosineSimilarity(metricEmb, evidenceEmb)
			if err != nil {
				return nil, fmt.Errorf("error computing similarity for evidence %d (%s): %v", i+1, evidence.lang, err)
			}

			predictedRelevant := sim > 0.5
			fmt.Printf("Evidence %d (%s): Similarity=%.4f, Predicted=%v, Actual=%v\n", i+1, evidence.lang, sim, predictedRelevant, evidence.relevant)
			if predictedRelevant == evidence.relevant {
				correct++
			}
		}

		accuracy := float64(correct) / float64(len(evidenceChunks))
		fmt.Printf("Accuracy: %.4f (%d/%d correct)\n", accuracy, correct, len(evidenceChunks))
		results[model] = TaskResult{Metric: accuracy}
	}

	// Determine winner
	maxAccuracy := -1.0
	var winner string
	for model, result := range results {
		if result.Metric > maxAccuracy {
			maxAccuracy = result.Metric
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
	fmt.Println("Registering Spanish Cross-Language Metric Evidence Task")
	RegisterTask(&spanishCrossLanguageMetricEvidenceTask{})
}
