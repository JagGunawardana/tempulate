package main

import (
	"github.com/oliveagle/jsonpath"
)

func fetchValue(in interface{}, query string) (interface{}, error) {
	return jsonpath.JsonPathLookup(in, query)
}
