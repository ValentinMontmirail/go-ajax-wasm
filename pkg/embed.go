package pkg

import (
	_ "embed"
)

//go:embed authors/schema.sql
var AuthorSchema string
