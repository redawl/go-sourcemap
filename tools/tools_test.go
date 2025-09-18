package tools

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var testFiles = []string{
	"test1.js.map",
	"test2.js.map",
}

func TestParseSourceMapFromUrl(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("../testdata")))

	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			_, err := ParseSourceMapFromUrl(server.URL + "/" + testFile)

			if err != nil {
				t.Errorf("Error parsing %s: %v", testFile, err)
			}
		})
	}
}

func TestParseSourceMapFromFile(t *testing.T) {
	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			_, err := ParseSourceMapFromFile("../testdata/" + testFile)

			if err != nil {
				t.Errorf("Error parsing %s: %v", testFile, err)
			}
		})
	}
}
