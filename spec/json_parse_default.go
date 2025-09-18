//go:build !jsonv2

package spec

import "encoding/json"

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
