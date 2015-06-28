package utils

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

func CellarLocation() string {
	wd, _ := os.Getwd()
	location := "Cellar/github.com/jcarley/docker-vault"

	new_path := path.Join(wd, location)
	return new_path
}

func Encode(writer io.Writer, data interface{}) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(data)
}

func Decode(reader io.Reader, data interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(data)
}
