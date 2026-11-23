package postgres

import "strings"

var likePatternEscaper = strings.NewReplacer(
	`\`, `\\`,
	`%`, `\%`,
	`_`, `\_`,
)

// escapeLikePattern escapes PostgreSQL LIKE/ILIKE metacharacters so user search
// terms are matched literally inside a surrounding "%...%" pattern.
func escapeLikePattern(input string) string {
	return "%" + likePatternEscaper.Replace(input) + "%"
}
