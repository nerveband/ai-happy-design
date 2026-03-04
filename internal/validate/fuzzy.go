package validate

import "strings"

// FuzzyMatchEnum finds the closest enum value for a given input.
// Returns the match and true if found, or "" and false.
func FuzzyMatchEnum(input string, enum []string) (string, bool) {
	lower := strings.ToLower(strings.TrimSpace(input))

	// Exact case-insensitive match first
	for _, v := range enum {
		if strings.ToLower(v) == lower {
			return v, true
		}
	}

	// Strip hyphens/underscores/spaces and retry
	normalized := strings.NewReplacer("-", "", "_", "", " ", "").Replace(lower)
	for _, v := range enum {
		norm := strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(v))
		if norm == normalized {
			return v, true
		}
	}

	// Levenshtein distance <= 3
	best := ""
	bestDist := 4
	for _, v := range enum {
		d := levenshtein(lower, strings.ToLower(v))
		if d < bestDist {
			bestDist = d
			best = v
		}
	}
	if bestDist <= 3 {
		return best, true
	}
	return "", false
}

// FuzzyMatchCommand finds the closest command name for a given input.
func FuzzyMatchCommand(input string, commands []string) (string, bool) {
	lower := strings.ToLower(input)

	// Exact match
	for _, c := range commands {
		if strings.ToLower(c) == lower {
			return c, true
		}
	}

	// Prefix/contains
	for _, c := range commands {
		cl := strings.ToLower(c)
		if strings.HasPrefix(cl, lower) || strings.Contains(cl, lower) {
			return c, true
		}
	}

	// Levenshtein
	best := ""
	bestDist := 4
	for _, c := range commands {
		d := levenshtein(lower, strings.ToLower(c))
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	if bestDist <= 3 {
		return best, true
	}
	return "", false
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			ins := curr[j-1] + 1
			del := prev[j] + 1
			sub := prev[j-1] + cost
			m := ins
			if del < m {
				m = del
			}
			if sub < m {
				m = sub
			}
			curr[j] = m
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}
