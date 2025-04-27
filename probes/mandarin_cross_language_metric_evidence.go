package probes

import (
	"fmt"
)

type mandarinCrossLanguageMetricEvidenceTask struct{}

func (t *mandarinCrossLanguageMetricEvidenceTask) Name() string {
	return "Mandarin Cross-Language Metric Evidence Task"
}

func (t *mandarinCrossLanguageMetricEvidenceTask) MetricName() string {
	return "Accuracy"
}

func (t *mandarinCrossLanguageMetricEvidenceTask) Run(models []string) (map[string]TaskResult, error) {
	type Evidence struct {
		text     string
		lang     string
		relevant bool
	}
	metric := "How comprehensive are the organization’s employee wellness programs?"
	evidenceChunks := []Evidence{
		{
			text:     "在Horizon公司，我们优先考虑员工的福祉，提供全面的健康计划，包括健身房会员、现场心理健康咨询和每季度的财务规划研讨会。去年，我们增加了正念课程和补贴健康饮食计划，以支持员工的全面健康。",
			lang:     "Mandarin",
			relevant: true, // Class A: comprehensive
		},
		{
			text:     "我们公司重视健康，在总部设有健身中心并举办年度健康博览会。员工可享受健身房折扣价并参加春季步数挑战赛。我们专注于身体健康，但根据年度员工调查的反馈正在探索更多选择。",
			lang:     "Mandarin",
			relevant: true, // Class B: moderate
		},
		{
			text:     "Horizon公司致力于交付高质量产品，团队努力工作以按时完成任务。我们最近升级了办公室，配备人体工学家具和现代休息室，以提高长时间工作的舒适度。",
			lang:     "Mandarin",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "我们公司专注于创新和生产力。我们引入了新的项目管理工具，以优化工作流程并确保客户项目按时交付。每周团队会议帮助统一目标并解决问题。",
			lang:     "Mandarin",
			relevant: false, // Class C: minimal/no wellness programs
		},
		{
			text:     "在Horizon公司，我们的工程师团队致力于尖端项目。我们投资了最先进的工作站并提供持续培训以保持技术技能的更新。公司食堂已翻新，增加了快餐选择。",
			lang:     "Mandarin",
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
	fmt.Println("Registering Mandarin Cross-Language Metric Evidence Task")
	RegisterTask(&mandarinCrossLanguageMetricEvidenceTask{})
}
