package munge

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// joinString joins a string with a delimiter
func joinString(in []interface{}, delim string) string {
	var holder []string
	for _, v := range in {
		s, ok := v.(string)
		if !ok {
			log.Fatal("Not a string for joinString")
		}
		holder = append(holder, s)
	}
	return strings.Join(holder, delim)
}

// joinStringComma joins a string with a comma delimiter
func joinStringComma(in []interface{}) string {
	return joinString(in, ",")
}

// getEnvDefault gets a environment variable value, using a default if not set
func getEnvDefault(name string, defaultVal string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultVal
	}
	return v
}

// createValue generates the value function for extracting from YAML/JSON
func createValue(paramFiles []string) func(string) interface{} {
	var fileContents []string
	for _, fileName := range paramFiles {
		contents, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal("Failed to read from file")
		}
		fileContents = append(fileContents, string(contents))
	}
	return func(jsonPath string) interface{} {
		for _, content := range fileContents {
			val, err := ExtractPath(content, jsonPath)
			if err == nil {
				return val
			}
		}
		return ""
	}
}
