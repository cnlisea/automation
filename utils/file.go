package utils

import (
	"io/ioutil"
	"os"
)

func FileExist(filename string) bool {
	var exist = true

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func ReadFileAll(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
