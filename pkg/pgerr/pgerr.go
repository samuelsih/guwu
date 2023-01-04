package pgerr

import (
	"strings"

	"github.com/lib/pq"
)

func UniqueColumn(err error) (string, error) {
	if e, ok := err.(*pq.Error); ok {
		if(e.Code == "23505") {
			return getUniqueColumn(e.Constraint), e
		}
	}

	return "", nil
}

func getUniqueColumn(str string) string {
	// users_column_key
	s := strings.Split(str, "_")

	return strings.Join(s[1:len(s) - 1], " ")
}