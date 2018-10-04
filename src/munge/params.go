package munge

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/oliveagle/jsonpath"
)

func unmarshalEncoded(in []byte) (interface{}, error) {
	var inJson interface{}
	err := json.Unmarshal([]byte(in), &inJson)
	if err == nil {
		return inJson, nil
	}
	err = yaml.Unmarshal([]byte(in), &inJson)
	if err != nil {
		return "", err
	}
	return inJson, nil
}

// ExtractPath takes valid YAML or JSON and carries out a JSON path query
func ExtractPath(in string, query string) (interface{}, error) {
	inJson, err := unmarshalEncoded([]byte(in))
	if err != nil { // Could not decode
		return "", err
	}
	res, err := jsonpath.JsonPathLookup(inJson, query)
	if err != nil {
		return res, err
	}
	return res, nil
}
