package utils

import (
	"os"
	"path"
)

func CellarLocation() string {
	wd, _ := os.Getwd()
	location := "Cellar/github.com/jcarley/docker-vault"

	new_path := path.Join(wd, location)
	return new_path
}
