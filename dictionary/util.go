package dictionary

import (
	"sort"
	"strings"
)

type myTokens []string

func (s myTokens) Len() int {
	return len(s)
}

func (s myTokens) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func (s myTokens) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/// Tokenise a string and sort by length desc
func tokenise(input string) []string {
	tokens := strings.Split(input, " ")
	var output []string

	for i := 0; i < len(tokens); i++ {
		for j := i + 1; j <= len(tokens); j++ {
			s := strings.Join(tokens[i:j], " ")
			if len(s) >= 3 {
				output = append(output, s)
			}
		}
	}

	sort.Sort(sort.Reverse(myTokens(output)))

	return output
}
