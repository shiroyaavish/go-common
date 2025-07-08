package common

import (
	"regexp"
	"strings"
)

var stringRg = regexp.MustCompile(`\s+`)
var alpha = regexp.MustCompile("[^a-zA-Z0-9 ]+")

// ToTsQuery takes a query string and returns a transformed version suitable for text search queries.
// It performs the following steps:
// 1. Converts the query string to lowercase.
// 2. Removes all non-alphanumeric characters except spaces.
// 3. Splits the string into substrings using spaces as delimiters.
// 4. Concatenates the substrings with a '|' separator.
//
// Parameters:
// - query: The input query string to be transformed.
//
// Returns:
// - The transformed query string.
func ToTsQuery(query string) string {
	if len(query) == 0 {
		return query
	}

	query = strings.ToLower(query)

	sp := strings.Split(stringRg.ReplaceAllString(alpha.ReplaceAllString(query, ""), "|"), "|")

	appliedCount := 0
	var builder strings.Builder

	for _, s := range sp {
		if len(s) == 0 {
			continue
		}

		if appliedCount != 0 {
			builder.WriteString("|")
		}

		builder.WriteString(s)
		appliedCount += 1
	}

	return builder.String()
}
