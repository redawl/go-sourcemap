//go:build jsonv2

// This file will be used when Go 1.25+ is available and the user builds with:
// GOEXPERIMENT=jsonv2 go build -tags=jsonv2 ./...
// or just: go build -tags=jsonv2 ./... (if json/v2 is stable)
//
// The json/v2 package provides better performance and additional features
// compared to the standard encoding/json package.

package spec

import "encoding/json/v2"

// ParseJSON parses str into a SourceMap object using json.Unmarshal for better performance
// Returns error if str is not valid json, or if the json object is not a SourceMap
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-ParseJSON
func ParseJSON(str string) (*SourceMap, error) {
	sourceMap := &SourceMap{}

	// Use json.Unmarshal directly which is often faster for strings
	err := json.Unmarshal([]byte(str), sourceMap)
	if err != nil {
		return nil, err
	}

	return sourceMap, nil
}
