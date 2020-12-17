package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
	"unicode"
)

type wordCount struct {
	word  string
	count int
}

type wordFrequency struct {
	data []*wordCount
}

func (c *wordFrequency) count(word string) {
	var incremented = false

	for _, v := range c.data {
		if v.word == word {
			v.count++
			incremented = true
		}
	}

	if !incremented {
		c.data = append(c.data, &wordCount{word: word, count: 1})
	}
}

func (c wordFrequency) less(i, j int) bool {
	return c.data[i].count > c.data[j].count
}

func wordEdge(c rune) bool {
	return c != 45 && !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

func Top10(s string) []string {
	s = strings.ToLower(s)

	var counter wordFrequency
	for _, word := range strings.FieldsFunc(s, wordEdge) {
		// когда вместо тире используют дефис
		if word == "-" {
			continue
		}

		counter.count(word)
	}

	sort.Slice(counter.data, counter.less)

	var top = 10
	if len(counter.data) < 10 {
		top = len(counter.data)
	}

	result := make([]string, top)
	for i, c := range counter.data[:top] {
		result[i] = c.word
	}

	return result
}
