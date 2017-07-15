package db

import (
	"fmt"
)

type DbError string
type NotFound string

func (e DbError) Error() string {
	return fmt.Sprintf("Generic DB error: %v", e)
}

func (e NotFound) Error() string {
	return string(e)
}
