package fuzzy

func LevenshteinDistance(s1, s2 string) int {
	// 1. Преобразуем в []rune для поддержки Unicode (кириллицы, эмодзи и т.д.)
	r1 := []rune(s1)
	r2 := []rune(s2)

	n, m := len(r1), len(r2)

	// Оптимизация: чтобы n всегда было меньше или равно m (экономим память)
	if n > m {
		r1, r2 = r2, r1
		n, m = m, n
	}

	// Инициализируем строки
	current_row := make([]int, n+1)
	previous_row := make([]int, n+1)

	// Базовый случай: превращение пустой строки в r1[:i]
	for i := 0; i <= n; i++ {
		previous_row[i] = i
	}

	for i := 1; i <= m; i++ {
		current_row[0] = i
		for j := 1; j <= n; j++ {
			// Расчет стоимостей
			add := previous_row[j] + 1
			delete := current_row[j-1] + 1
			change := previous_row[j-1]

			if r1[j-1] != r2[i-1] {
				change++
			}

			// Выбираем минимум
			current_row[j] = min(add, delete, change)
		}
		// Копируем текущую строку в предыдущую
		copy(previous_row, current_row)
	}

	return previous_row[n]
}

func IsMatch(word, target string, threshold int) bool {
	return LevenshteinDistance(word, target) <= threshold
}
