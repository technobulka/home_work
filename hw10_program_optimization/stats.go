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
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
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
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	result := make(DomainStat)
	domain = "." + domain
	var user User

	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		user = User{}
		err = json.Unmarshal(line, &user)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
