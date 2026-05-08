package counter

import (
	"strings"
	"sync"
	"time"
	"unicode"
)

type WordCountResult struct {
	TotalWords      int
	UniqueWords     int
	WordFrequencies map[string]int
	ProcessingTime  time.Duration
}

type Task struct {
	Content string
}

type Result struct {
	WordCounts map[string]int
	TotalWords int
}

// BNF Grammar for word parsing:
// <word> ::= <letter> { <letter> | <digit> }
// <letter> ::= "a" | "b" | ... | "z" | "A" | ... | "Z"
// <digit> ::= "0" | "1" | ... | "9"

type BNFParseState struct {
	input string
	pos   int
}

// NewBNFParser creates a new BNF parser for the given input string
func NewBNFParser(input string) *BNFParseState {
	return &BNFParseState{input: strings.ToLower(input), pos: 0}
}

func (p *BNFParseState) isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func (p *BNFParseState) isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func (p *BNFParseState) current() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	return rune(p.input[p.pos])
}

func (p *BNFParseState) advance() {
	if p.pos < len(p.input) {
		p.pos++
	}
}

func (p *BNFParseState) skipNonWordChars() {
	for p.pos < len(p.input) && !p.isLetter(p.current()) && !p.isDigit(p.current()) {
		p.advance()
	}
}

// implements <word> ::= <letter> { <letter> | <digit> }
func (p *BNFParseState) parseWord() string {
	p.skipNonWordChars()

	if p.pos >= len(p.input) || (!p.isLetter(p.current()) && !p.isDigit(p.current())) {
		return ""
	}

	start := p.pos

	// Parse first character (must be letter according to BNF)
	if p.isLetter(p.current()) {
		p.advance()
	} else {
		return ""
	}

	// Parse remaining characters (letters or digits)
	for p.pos < len(p.input) && (p.isLetter(p.current()) || p.isDigit(p.current())) {
		p.advance()
	}

	return p.input[start:p.pos]
}

// extracts all words from the input string
func (p *BNFParseState) parseAllWords() []string {
	var words []string

	for p.pos < len(p.input) {
		word := p.parseWord()
		if word != "" {
			words = append(words, word)
		} else {
			p.advance()
		}
	}

	return words
}

func worker(tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		counts := make(map[string]int)
		total := 0

		// Parse words using BNF grammar
		parser := NewBNFParser(task.Content)
		words := parser.parseAllWords()

		for _, word := range words {
			if word != "" {
				counts[word]++
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
		TotalWords:      totalWords,
		UniqueWords:     len(finalCounts),
		WordFrequencies: finalCounts,
		ProcessingTime:  time.Since(startTime),
	}
}
