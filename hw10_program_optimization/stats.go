package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domains, err := countDomains(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return domains, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + domain
	var userLine string
	var etIdx, quoteIdx int

	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		userLine = string(line)
		etIdx = strings.Index(userLine, "@")
		quoteIdx = strings.Index(userLine[etIdx:], `"`)
		userDomain := strings.ToLower(userLine[etIdx+1 : etIdx+quoteIdx])

		if strings.HasSuffix(userDomain, domain) {
			result[userDomain]++
		}
	}

	return result, nil
}
