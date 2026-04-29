package annotations

import "strings"

func trimQuotes(s string) string {
	s = strings.TrimSpace(s) // 👈 Сначала убираем внешние пробелы
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			s = s[1 : len(s)-1]
		}
	}
	return strings.TrimSpace(s) // 👈 Убираем пробелы, которые могли быть внутри кавычек
}
