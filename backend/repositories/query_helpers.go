package repositories

import "strings"

func safeSortBy(requested string, allowed map[string]string, fallback string) string {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if column, ok := allowed[requested]; ok {
		return column
	}

	return fallback
}

func safeSortOrder(requested string) string {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if requested == "asc" {
		return "asc"
	}

	return "desc"
}
