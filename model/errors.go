package model

import (
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func wrapErr(errSQL error, domain string) error {
	if err, ok := errSQL.(*pgconn.PgError); ok {
		switch err.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("%s already exists", domain)

		case pgerrcode.NotNullViolation:
			return fmt.Errorf("%s is required", beautifyColumn(err.ColumnName))
		}
	}
	
	return errSQL
}

func beautifyColumn(column string) string {
	var sb strings.Builder
	
	str := strings.Split(column, "_")

	if len(str) == 1 {
		return str[0]
	}

	for i := 0; i < len(str); i++ {
		if i == len(str)-1 {
			sb.WriteString(str[i])
		} else {
			sb.WriteString(str[i] + " ")
		}
	}

	return sb.String()
}