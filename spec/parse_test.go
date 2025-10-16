package spec

import (
	"os"
	"testing"
)

var testFiles = []string{
	"test1.js.map",
	"test2.js.map",
}

func getTestFileContents(filename string) (string, error) {
	contents, err := os.ReadFile("../testdata/" + filename)

	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func TestParseSourceMap(t *testing.T) {
	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			contents, err := getTestFileContents(testFile)

			if err != nil {
				t.Errorf("Error getting contents of %s: %v", testFile, err)
				return
			}

			_, err = ParseSourceMap(contents, "")

			if err != nil {
				t.Errorf("Error parsing %s: %v", testFile, err)
			}
		})
	}
}
