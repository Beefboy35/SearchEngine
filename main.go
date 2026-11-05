package main

import (
	"fmt"
	"sort"

	"github.com/username/myproject/pkg/engine"
)

func main() {
	// 1. Коллекция документов
	rawDocs := []string{
		"Machine learning is a branch of artificial intelligence",
		"Deep learning uses neural networks for machine learning",
		"Artificial intelligence is transforming modern technology",
		"Machine learning algorithms learn patterns from data",
		"Bananas and apples are delicious fruits",
	}

	// 2. Словарь для Preprocess
	dict := []string{
		"machine", "learning", "artificial",
		"intelligence", "deep", "neural",
		"networks", "algorithms", "data",
	}

	// 3. Предобработка документов
	docs := make([]engine.ProcessedDoc, 0, len(rawDocs))

	for _, raw := range rawDocs {
		doc, err := engine.Preprocess(raw, dict, 0)
		if err != nil {
			panic(err)
		}

		docs = append(docs, doc)
	}

	// 4. Средняя длина документов
	avgLen := engine.AvrDocLength(docs)

	// 5. Параметры BM25
	params := engine.BM25Params{
		K1:           1.2,
		B:            0.75,
		AvgDocLength: avgLen,
		TotalDocs:    len(docs),
	}

	// 6. Поисковый запрос
	queryDoc, err := engine.Preprocess(
		"machine learning",
		dict,
		0,
	)
	if err != nil {
		panic(err)
	}

	// 7. Вычисляем оценки
	type Result struct {
		Index int
		Score float64
	}

	results := make([]Result, 0, len(docs))

	for i, doc := range docs {
		score := engine.BM25Score(
			queryDoc.Tokens,
			doc,
			docs,
			params,
		)

		results = append(results, Result{
			Index: i,
			Score: score,
		})
	}

	// 8. Сортируем по убыванию релевантности
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 9. Выводим результаты
	fmt.Println("Query: machine learning")
	fmt.Printf("Documents: %d\n", len(docs))
	fmt.Printf("Average document length: %.2f\n\n", avgLen)

	for _, result := range results {
		fmt.Printf(
			"Score: %.6f | Document %d: %s\n",
			result.Score,
			result.Index+1,
			rawDocs[result.Index],
		)
	}
}
