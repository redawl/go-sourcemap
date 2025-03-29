// Package tools contains miscellaneous utilities for using the sourcemap implementation in package spec.
// These tools are not part of the spec, and are included for convenience.
package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/redawl/go-sourcemap/spec"
)

// Parse functions

// ParseSourceMapFromUrl parses a source map file located at url.
// Returns an error if url is unreachable, returns a status != 200, or url is not a valid source map file.
func ParseSourceMapFromUrl(url string) (*spec.DecodedSourceMapRecord, error) {
    response, err := http.Get(url)

    if err != nil {
        return nil, err
    }
    if response.StatusCode != 200 {
        return nil, fmt.Errorf("Error retrieving %s: %s", url, response.Status)
    }

    contents, err := io.ReadAll(response.Body)

    if err != nil {
        return nil, fmt.Errorf("Error reading response body: %v", err)
    }

    return spec.ParseSourceMap(string(contents), url)
}
// ParseSourceMapFromFile
// TODO: Docs
func ParseSourceMapFromFile(filename string) (*spec.DecodedSourceMapRecord, error) {
    contents, err := os.ReadFile(filename) 

    if err != nil {
        return nil, fmt.Errorf("Error reading contents of %s: %v", filename, err)
    }

    return spec.ParseSourceMap(string(contents), filename)
}

// SaveSourcesToDirectory saves mapRecord.Sources to dir. 
// If dir doesn't exist, it is recursively created with 0700 permissions.
// Files are saved with 0600 permissions.
func SaveSourcesToDirectory(mapRecord *spec.DecodedSourceMapRecord, dir string) error {
    err := os.MkdirAll(dir, 0700)

    if err != nil {
        return fmt.Errorf("Error creating %s: %v", dir, err)
    }

    for _, source := range mapRecord.Sources {
        index := strings.LastIndexByte(source.Url, os.PathSeparator)

        if index != -1 {
            err = os.MkdirAll(dir + source.Url[:index], 0700)

            if err != nil {
                return fmt.Errorf("Error creating %s: %v", dir + source.Url[:index], err)
            }

            err := os.WriteFile(dir + source.Url, []byte(source.Content), 0600)

            if err != nil {
                return fmt.Errorf("Error writing file contents to %s: %v", dir + source.Url, err)
            }
        }
    }

    return nil
}
// StringifyDecodedSourceMapRecord
// TODO: Docs
func StringifyDecodedSourceMapRecord(mapRecord *spec.DecodedSourceMapRecord) (string, error) {
    str, err := json.Marshal(mapRecord)

    if err != nil {
        return "", fmt.Errorf("Error stringifying mapRecord: %v", err)
    }

    return string(str), nil
}

