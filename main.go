package main

import (
	"fmt"

	"github.com/username/myproject/pkg/engine"
)

func main() {
	dict := []string{"apple", "banana", "eat"}
	text := "I ate a banana and an apple"

	doc, _ := engine.Preprocess(text, dict, 1)

	tf := engine.TermFrequency("apple", doc)
	idf := engine.CalculateIDF(1, 1) // Simplified for example
	fmt.Printf("Tokens: %v, TF: %d, IDF: %f\n", doc.Tokens, tf, idf)
	bm25_score := engine.BM25Score(tf, len(doc.Tokens), engine.BM25Params{K1: 0.05, B: 0.05, AvgDocLength: engine.AvrDocLength(doc.Tokens), TotalDocs: 1})
	fmt.Printf("BM25 Score: %f\n", bm25_score)

}
