package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	probes "embedding-probes/probes"
)

type Config struct {
	Models []string `json:"models"`
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Debug: Print registered tasks
	fmt.Printf("Registered tasks: %d\n", len(probes.TaskRegistry))
	for i, task := range probes.TaskRegistry {
		fmt.Printf("Task %d: %s\n", i+1, task.Name())
	}

	configPath, err := filepath.Abs("config.json")
	if err != nil {
		fmt.Printf("Error resolving config path: %v\n", err)
		return
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("Error parsing config file: %v", err)
		return
	}

	results := make(map[int]map[string]probes.TaskResult)
	for i, task := range probes.TaskRegistry {
		taskNum := i + 1
		fmt.Printf("\nTask %d: %s\n", taskNum, task.Name())
		taskResults, err := task.Run(config.Models)
		if err != nil {
			fmt.Printf("Error running task %d: %s: %v\n", taskNum, task.Name(), err)
			return
		}
		results[taskNum] = taskResults
	}

	// Calculate maximum widths for each column
	maxTaskWidth := len("Task")
	maxTaskNameWidth := len("Task Name")
	maxMetricWidth := len("Metric")
	maxGraniteWidth := len("granite-embedding:latest")
	maxNomicWidth := len("nomic-embed-text")
	maxWinnerWidth := len("Winner")

	for taskNum := 1; taskNum <= len(results); taskNum++ {
		task := probes.TaskRegistry[taskNum-1]
		taskNumStr := fmt.Sprintf("%d", taskNum)
		graniteMetric := fmt.Sprintf("%.4f", results[taskNum]["granite-embedding:latest"].Metric)
		nomicMetric := fmt.Sprintf("%.4f", results[taskNum]["nomic-embed-text"].Metric)
		winner := results[taskNum]["granite-embedding:latest"].Winner
		if winner == "" {
			winner = results[taskNum]["nomic-embed-text"].Winner
		}
		if winner == "" {
			winner = "Tie"
		}

		maxTaskWidth = max(maxTaskWidth, len(taskNumStr))
		maxTaskNameWidth = max(maxTaskNameWidth, len(task.Name()))
		maxMetricWidth = max(maxMetricWidth, len(task.MetricName()))
		maxGraniteWidth = max(maxGraniteWidth, len(graniteMetric))
		maxNomicWidth = max(maxNomicWidth, len(nomicMetric))
		maxWinnerWidth = max(maxWinnerWidth, len(winner))
	}

	// Print table with dynamic widths
	fmt.Println("\nFinal Results Table:")
	headerFormat := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds |",
		maxTaskWidth, maxTaskNameWidth, maxMetricWidth, maxGraniteWidth, maxNomicWidth, maxWinnerWidth)
	// Generate separator with proper width including borders and spaces
	separator := fmt.Sprintf("|%%-%ds|%%-%ds|%%-%ds|%%-%ds|%%-%ds|%%-%ds|",
		maxTaskWidth+2, maxTaskNameWidth+2, maxMetricWidth+2,
		maxGraniteWidth+2, maxNomicWidth+2, maxWinnerWidth+2)
	separator = fmt.Sprintf(separator,
		strings.Repeat("-", maxTaskWidth),
		strings.Repeat("-", maxTaskNameWidth),
		strings.Repeat("-", maxMetricWidth),
		strings.Repeat("-", maxGraniteWidth),
		strings.Repeat("-", maxNomicWidth),
		strings.Repeat("-", maxWinnerWidth))

	fmt.Fprintf(os.Stdout, headerFormat+"\n", "Task", "Task Name", "Metric", "granite-embedding:latest", "nomic-embed-text", "Winner")
	fmt.Println(separator)

	for taskNum := 1; taskNum <= len(results); taskNum++ {
		task := probes.TaskRegistry[taskNum-1]
		graniteMetric := fmt.Sprintf("%.4f", results[taskNum]["granite-embedding:latest"].Metric)
		nomicMetric := fmt.Sprintf("%.4f", results[taskNum]["nomic-embed-text"].Metric)
		winner := results[taskNum]["granite-embedding:latest"].Winner
		if winner == "" {
			winner = results[taskNum]["nomic-embed-text"].Winner
		}
		if winner == "" {
			winner = "Tie"
		}
		fmt.Fprintf(os.Stdout, headerFormat+"\n",
			fmt.Sprintf("%d", taskNum), task.Name(), task.MetricName(), graniteMetric, nomicMetric, winner)
	}

	graniteWins := 0
	nomicWins := 0
	for taskNum := 1; taskNum <= len(results); taskNum++ {
		if results[taskNum]["granite-embedding:latest"].Winner == "granite-embedding:latest" {
			graniteWins++
		} else if results[taskNum]["nomic-embed-text"].Winner == "nomic-embed-text" {
			nomicWins++
		}
	}
	fmt.Printf("\nOverall Reliability:\n")
	fmt.Printf("granite-embedding:latest wins: %d\n", graniteWins)
	fmt.Printf("nomic-embed-text wins: %d\n", nomicWins)
	if graniteWins > nomicWins {
		fmt.Printf("granite-embedding:latest is more reliable (%d vs. %d wins).\n", graniteWins, nomicWins)
	} else if nomicWins > graniteWins {
		fmt.Printf("nomic-embed-text is more reliable (%d vs. %d wins).\n", nomicWins, graniteWins)
	} else {
		fmt.Println("Both models are equally reliable.")
	}
}
