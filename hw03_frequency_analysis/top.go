package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	splitSlice := strings.Fields(input)

	freqMap := make(map[string]int)

	for _, word := range splitSlice {
		freqMap[word]++
	}

	uniqSlice := make([]string, len(freqMap))
	i := 0
	for word := range freqMap {
		uniqSlice[i] = word
		i++
	}
	sort.Slice(uniqSlice, func(i, j int) bool {
		f1 := freqMap[uniqSlice[i]]
		f2 := freqMap[uniqSlice[j]]
		if f1 == f2 {
			return uniqSlice[i] < uniqSlice[j]
		}
		return f1 > f2
	})

	if len(uniqSlice) < 10 {
		return uniqSlice
	}
	return uniqSlice[:10]
}
