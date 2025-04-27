package probes

import (
	"fmt"
)

type frenchCrossLanguageMetricEvidenceTask struct{}

func (t *frenchCrossLanguageMetricEvidenceTask) Name() string {
	return "French Cross-Language Metric Evidence Task"
}

func (t *frenchCrossLanguageMetricEvidenceTask) MetricName() string {
	return "Accuracy"
}

func (t *frenchCrossLanguageMetricEvidenceTask) Run(models []string) (map[string]TaskResult, error) {
	type Evidence struct {
		text     string
		lang     string
		relevant bool
	}
	metric := "How comprehensive are the organization’s employee wellness programs?"
	evidenceChunks := []Evidence{
		{
			text:     "Chez Horizon Inc., nous priorisons le bien-être des employés avec un programme complet comprenant des abonnements à des salles de sport, des services de conseil en santé mentale par un thérapeute sur place et des ateliers de planification financière trimestriels. L’an dernier, nous avons ajouté des sessions de pleine conscience et un plan de repas sains subventionné.",
			lang:     "French",
			relevant: true, // Class A: comprehensive
		},
		{
			text:     "Notre entreprise valorise la santé et offre un centre de fitness au siège avec une foire annuelle de la santé. Les employés bénéficient de tarifs réduits pour les salles de sport et participent à un défi de pas au printemps. Nous nous concentrons sur la forme physique, mais explorons d’autres options basées sur les retours de notre enquête annuelle.",
			lang:     "French",
			relevant: true, // Class B: moderate
		},
		{
			text:     "Horizon Inc. s’engage à fournir des produits de haute qualité, avec des équipes travaillant dur pour respecter les délais. Nous avons récemment modernisé nos bureaux avec des meubles ergonomiques et une salle de pause contemporaine pour améliorer le confort pendant les longues heures de travail.",
			lang:     "French",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "Notre entreprise se concentre sur l’innovation et la productivité. Nous avons introduit un nouvel outil de gestion de projet pour rationaliser les flux de travail et assurer une livraison ponctuelle des projets clients. Des réunions d’équipe hebdomadaires alignent les objectifs et résolvent les défis.",
			lang:     "French",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "Chez Horizon Inc., notre équipe d’ingénieurs travaille sur des projets de pointe. Nous avons investi dans des stations de travail de dernière génération et offrons une formation continue pour maintenir les compétences techniques à jour. La cafétéria a été rénovée pour inclure des options de restauration rapide.",
			lang:     "French",
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
				return nil, fmt.Errorf("error getting personally for evidence %d (%s): %v", i+1, evidence.lang, err)
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
	fmt.Println("Registering French Cross-Language Metric Evidence Task")
	RegisterTask(&frenchCrossLanguageMetricEvidenceTask{})
}
