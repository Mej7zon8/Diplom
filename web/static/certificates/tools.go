package certificates

import "os"

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
