package probes

import (
	"fmt"
)

type russianCrossLanguageMetricEvidenceTask struct{}

func (t *russianCrossLanguageMetricEvidenceTask) Name() string {
	return "Russian Cross-Language Metric Evidence Task"
}

func (t *russianCrossLanguageMetricEvidenceTask) MetricName() string {
	return "Accuracy"
}

func (t *russianCrossLanguageMetricEvidenceTask) Run(models []string) (map[string]TaskResult, error) {
	type Evidence struct {
		text     string
		lang     string
		relevant bool
	}
	metric := "How comprehensive are the organization’s employee wellness programs?"
	evidenceChunks := []Evidence{
		{
			text:     "В Horizon Inc. мы уделяем приоритетное внимание благополучию сотрудников с помощью комплексной программы, включающей абонементы в спортзал, консультации по психическому здоровью с терапевтом на месте и ежеквартальные семинары по финансовому планированию. В прошлом году мы добавили занятия по осознанности и субсидируемый план здорового питания.",
			lang:     "Russian",
			relevant: true, // Class A: comprehensive
		},
		{
			text:     "Наша компания ценит здоровье и предоставляет фитнес-центр в штаб-квартире с ежегодной ярмаркой здоровья. Сотрудники получают скидки на абонементы в спортзал и участвуют в весеннем шаговом марафоне. Мы сосредоточены на физической форме, но изучаем дополнительные возможности на основе отзывов из ежегодного опроса.",
			lang:     "Russian",
			relevant: true, // Class B: moderate
		},
		{
			text:     "Horizon Inc. стремится поставлять продукцию высокого качества, а команды усердно работают, чтобы соблюдать сроки. Недавно мы обновили офис эргономичной мебелью и современной комнатой отдыха, чтобы повысить комфорт во время долгих рабочих часов.",
			lang:     "Russian",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "Наша компания сосредоточена на инновациях и производительности. Мы внедрили новый инструмент управления проектами для оптимизации рабочих процессов и своевременной доставки проектов клиентам. Еженедельные встречи команды помогают согласовывать цели и решать проблемы.",
			lang:     "Russian",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "В Horizon Inc. наша команда инженеров работает над передовыми проектами. Мы инвестировали в современные рабочие станции и предлагаем непрерывное обучение для поддержания технических навыков. Столовая компании была обновлена, чтобы включить варианты быстрого питания.",
			lang:     "Russian",
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
	fmt.Println("Registering Russian Cross-Language Metric Evidence Task")
	RegisterTask(&russianCrossLanguageMetricEvidenceTask{})
}
