package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domains, err := countDomains(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return domains, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	var user User
	var line []byte
	var err error

	json := jsoniter.ConfigFastest
	result := make(DomainStat)

	br := bufio.NewReader(r)
	for {
		line, _, err = br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		err = json.Unmarshal(line, &user)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) &&
			strings.Contains(user.Email, "@") {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
