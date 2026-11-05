package engine

import (
	"errors"
	"math"
	"strings"

	"github.com/username/myproject/pkg/fuzzy"
	"github.com/username/myproject/pkg/text"
)

type BM25Params struct {
	K1           float64
	B            float64
	AvgDocLength float64
	TotalDocs    int
}

var ErrNilDictionary = errors.New("dictionary cannot be nil")

func Preprocess(rawText string, dictionary []string, threshold int) (ProcessedDoc, error) {
	// 1. Валидация словаря
	if dictionary == nil {
		return ProcessedDoc{}, ErrNilDictionary
	}

	if threshold < 0 {
		return ProcessedDoc{}, errors.New("threshold must be >= 0")
	}

	// 2. Lowercase с обработкой ошибки
	lowered := text.LowerCase(rawText)

	// 3. Tokenize с обработкой ошибки
	tokens := text.Tokenize(lowered)

	// 4. Fuzzy Correction
	correctedTokens := make([]string, len(tokens))
	for i, token := range tokens {
		bestMatch := token
		minDist := threshold + 1

		for _, dictWord := range dictionary {
			dist := fuzzy.LevenshteinDistance(token, dictWord)
			if dist < minDist {
				minDist = dist
				bestMatch = dictWord
			}
		}
		correctedTokens[i] = bestMatch
	}

	return ProcessedDoc{Tokens: correctedTokens}, nil
}

func TermFrequency(target string, doc ProcessedDoc) int {
	count := 0
	for _, t := range doc.Tokens {
		if t == target {
			count++
		}
	}
	return count
}

func AvrDocLength(docs []string) float64 {
	if len(docs) == 0 {
		return 0.0
	}
	var words float64
	for _, doc := range docs {
		words += float64(len(strings.Fields(doc)))
	}

	return float64(words) / float64(len(docs))
}
func normalizeDocLen(docLen, coef float64, docs []string) float64 {
	return 1 - coef + coef*(docLen/AvrDocLength(docs))
}

// CalculateSecondMonomial implements: (tf * (k1 + 1)) / (tf + k1 * (1 - b + b * (docLen / avgLen)))

func CalculateSecondMonomial(tf int, docLen int, params BM25Params) float64 {
	tfFloat := float64(tf)
	L := float64(docLen)

	// Вычисляем коэффициент нормализации длины документа
	// Formula: (1 - b) + b * (L / avgL)
	normalization := (1.0 - params.B) + params.B*(L/params.AvgDocLength)

	numerator := tfFloat * (params.K1 + 1.0)
	denominator := tfFloat + params.K1*normalization

	return numerator / denominator
}

// CalculateIDF реализует формулу IDF для BM25:
// IDF(q_i) = log( (N - n(q_i) + 0.5) / (n(q_i) + 0.5) + 1 )
// Или более распространенный вариант сглаживания:
// IDF(q_i) = log( (N - n(q_i) + 0.5) / (n(q_i) + 0.5) ) + 1 (но это зависит от реализации)
// Ниже приведена самая стабильная версия, используемая в современных поисковиках:
func CalculateIDF(n_qi int, totalDocs int) float64 {
	N := float64(totalDocs)
	n := float64(n_qi)

	// Добавляем 0.5 (сглаживание), чтобы избежать деления на 0
	// и чтобы значение не уходило в минус, если слово есть во всех документах.
	// Формула: log( (N - n + 0.5) / (n + 0.5) + 1 )
	// Добавление +1 в конце гарантирует, что IDF всегда будет > 0.

	idf := math.Log((N-n+0.5)/(n+0.5) + 1.0)

	return idf
}

func BM25Score(tf int, docLen int, params BM25Params) float64 {
	secondMonomial := CalculateSecondMonomial(tf, docLen, params)
	idf := CalculateIDF(tf, params.TotalDocs)
	return idf * secondMonomial
}
