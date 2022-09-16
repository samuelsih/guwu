package model

import (
	_ "embed"
)

//go:embed migration/user.sql
var schema string