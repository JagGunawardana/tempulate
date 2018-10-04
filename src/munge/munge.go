// package munge will take a template and a list of parameter files
// (either YAML or JSON) and turn them into a templated output file

package munge

import (
	"bytes"
	"text/template"
)

// MungeFile takes the input template, parameter files and returns the templated file
// templateContent is the template file contents (not file name)
func MungeFile(templateContent string, paramFiles []string) (outputFile string, err error) {
	t, err := template.New("template").Funcs(
		template.FuncMap{
			"envdef": getEnvDefault,
			"value":  createValue(paramFiles),
		},
	).Parse(templateContent)
	if err != nil {
		return "", err
	}
	tParams := map[string]string{
		"one": "two",
	}
	var outBuff bytes.Buffer
	err = t.Execute(&outBuff, tParams)
	if err != nil {
		return "", err
	}
	return outBuff.String(), nil
}
