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
        return nil, fmt.Errorf("Error reading response body: %w", err)
    }

    return spec.ParseSourceMap(string(contents), "")
}
// ParseSourceMapFromFile parses a source map file.
// Returns an error if the file is unreadable, or the file is not a valid source map file.
func ParseSourceMapFromFile(filename string) (*spec.DecodedSourceMapRecord, error) {
    contents, err := os.ReadFile(filename) 

    if err != nil {
        return nil, fmt.Errorf("Error reading contents of %s: %w", filename, err)
    }

    return spec.ParseSourceMap(string(contents), filename)
}

// SaveSourcesToDirectory saves mapRecord.Sources to dir. 
// If dir doesn't exist, it is recursively created with 0700 permissions.
// Files are saved with 0600 permissions.
func SaveSourcesToDirectory(mapRecord *spec.DecodedSourceMapRecord, dir string) error {
    err := os.MkdirAll(dir, 0700)

    if err != nil {
        return fmt.Errorf("Error creating %s: %w", dir, err)
    }

    for _, source := range mapRecord.Sources {
        index := strings.LastIndexByte(source.Url, os.PathSeparator)

        if index != -1 {
            basePath := dir + "/" + source.Url[:index]
            fullPath := basePath + source.Url[index:]
            err = os.MkdirAll(basePath, 0700)

            if err != nil {
                return fmt.Errorf("Error creating %s: %w", dir + source.Url[:index], err)
            }

            err := os.WriteFile(fullPath, []byte(source.Content), 0600)

            if err != nil {
                return fmt.Errorf("Error writing file contents to %s: %w", fullPath, err)
            }
        }
    }

    return nil
}
// MarshalDecodedSourceMapRecord returns the JSON encoding of mapRecord
func MarshalDecodedSourceMapRecord(mapRecord *spec.DecodedSourceMapRecord) (string, error) {
    str, err := json.Marshal(mapRecord)

    if err != nil {
        return "", fmt.Errorf("Error stringifying mapRecord: %w", err)
    }

    return string(str), nil
}

