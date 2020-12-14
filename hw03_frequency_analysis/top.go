package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
	"unicode"
)

type wordCount struct {
	Word  string
	Count int
}

type wordFrequency []*wordCount

func (c *wordFrequency) count(word string) {
	var incremented = false

	for _, v := range *c {
		if v.Word == word {
			v.Count++
			incremented = true
		}
	}

	if !incremented {
		*c = append(*c, &wordCount{Word: word, Count: 1})
	}
}

func Top10(s string) []string {
	s = strings.ToLower(s)
	f := func(c rune) bool {
		return c != 45 && !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	var counter wordFrequency
	for _, word := range strings.FieldsFunc(s, f) {
		// когда вместо тире используют дефис
		if word == "-" {
			continue
		}

		counter.count(word)
	}

	sort.Slice(counter, func(i, j int) bool {
		return counter[i].Count > counter[j].Count
	})

	var top = 10
	if len(counter) < 10 {
		top = len(counter)
	}

	result := make([]string, top)
	for i, c := range counter[:top] {
		result[i] = c.Word
	}

	return result
}
