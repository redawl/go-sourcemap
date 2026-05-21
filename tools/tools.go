// Package tools contains miscellaneous utilities for using the sourcemap implementation in package spec.
// These tools are not part of the spec, and are included for convenience.
package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return spec.ParseSourceMap(string(contents), "")
}

// ParseSourceMapFromFile parses a source map file.
// Returns an error if the file is unreadable, or the file is not a valid source map file.
func ParseSourceMapFromFile(filename string) (*spec.DecodedSourceMapRecord, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading contents of %s: %w", filename, err)
	}

	baseURL, _ := filepath.Split(filename)

	return spec.ParseSourceMap(string(contents), baseURL)
}

func trimPrefixAll(s, prefix string) string {
	result := strings.TrimPrefix(s, prefix)

	if strings.HasPrefix(s, prefix) {
		return trimPrefixAll(result, prefix)
	}

	return result
}

// SaveSourcesToDirectory saves mapRecord.Sources to dir.
// If dir doesn't exist, it is recursively created with 0700 permissions.
// Files are saved with 0600 permissions.
func SaveSourcesToDirectory(mapRecord *spec.DecodedSourceMapRecord, dir string) error {
	err := os.MkdirAll(dir, 0o700)
	if err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("opening root %s: %w", dir, err)
	}

	for _, source := range mapRecord.Sources {
		cleanedPath := trimPrefixAll(source.Url, "../")
		basePath := filepath.Dir(cleanedPath)

		err := root.MkdirAll(basePath, 0o700)
		if err != nil {
			return fmt.Errorf("creating %s: %w", basePath, err)
		}

		err = root.WriteFile(cleanedPath, []byte(source.Content), 0o600)
		if err != nil {
			return fmt.Errorf("writing file contents to %s: %w", cleanedPath, err)
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
