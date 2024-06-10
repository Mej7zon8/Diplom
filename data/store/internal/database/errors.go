package database

import "strings"

func IsErrorRecordNotFound(e error) bool {
	return e != nil && strings.Contains(e.Error(), "record not found")
}
