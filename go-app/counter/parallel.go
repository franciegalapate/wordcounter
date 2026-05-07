package counter

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

type WordCountResult struct {
	TotalWords int
	UniqueWords int
	WordFrequencies map[string]int
	ProcessingTime time.Duration
}

type Task struct {
	Content string
}

type Result struct {
	WordCounts map[string]int
	TotalWords int
}

func worker(tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	cleanRe := regexp.MustCompile(`[^\w]`)

	for task := range tasks {
		counts := make(map[string]int)
		total := 0

		// Split into words and count
		words := strings.Fields(strings.ToLower(task.Content))
		for _, word := range words {
			// Clean word: remove punctuation
			cleanWord := cleanRe.ReplaceAllString(word, "")
			if cleanWord != "" {
				counts[cleanWord]++
				total++
			}
		}

		results <- Result{
			WordCounts: counts,
			TotalWords: total,
		}
	}
}

func ParallelWordCount(chunks []string, numWorkers int) WordCountResult {
	startTime := time.Now()

	tasks := make(chan Task, len(chunks))
	results := make(chan Result, len(chunks))
	
	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(tasks, results, &wg)
	}

	go func() {
		for _, chunk := range chunks {
			tasks <- Task{Content: chunk}
		}
		close(tasks)
	}()

	// Close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and count words
	finalCounts := make(map[string]int)
	totalWords := 0

	for result := range results {
		totalWords += result.TotalWords
		for word, count := range result.WordCounts {
			finalCounts[word] += count
		}
	}

	return WordCountResult{
		TotalWords: totalWords,
		UniqueWords: len(finalCounts),
		WordFrequencies: finalCounts,
		ProcessingTime: time.Since(startTime),
	}
}
