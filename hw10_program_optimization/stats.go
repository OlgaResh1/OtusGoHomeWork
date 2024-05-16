package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	input := bufio.NewScanner(r)
	result := make(DomainStat)

	var err error
	dotDomain := "." + domain

	for input.Scan() {
		line := input.Bytes()

		var user User
		if err = user.UnmarshalJSON(line); err != nil {
			return nil, fmt.Errorf("error json : %w", err)
		}
		if strings.HasSuffix(user.Email, dotDomain) {
			n := strings.LastIndex(user.Email, "@")
			if n == -1 {
				return nil, fmt.Errorf("wrong email: %s", user.Email)
			}
			result[strings.ToLower(user.Email[n+1:])]++
		}
	}
	return result, nil
}
