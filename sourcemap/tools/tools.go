package sourcemap

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/redawl/go-sourcemap/sourcemap/spec"
)

// Parse functions


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
        return nil, fmt.Errorf("Error reading response body: %s", err.Error())
    }

    return spec.ParseSourceMap(string(contents), url)
}

// Save functions

func SaveSourcesToDirectory(mapRecord *spec.DecodedSourceMapRecord, dir string) error {
    err := os.MkdirAll(dir, 0700)

    if err != nil {
        return fmt.Errorf("Error creating %s: %s", dir, err.Error())
    }

    for _, source := range mapRecord.Sources {
        index := strings.LastIndexByte(source.Url, os.PathSeparator)

        if index != -1 {
            err = os.MkdirAll(dir + source.Url[:index], 0700)

            if err != nil {
                return fmt.Errorf("Error creating %s: %s", dir + source.Url[:index], err.Error())
            }

            err := os.WriteFile(dir + source.Url, []byte(source.Content), 0600)

            if err != nil {
                return fmt.Errorf("Error writing file contents to %s: %s", dir + source.Url, err.Error())
            }
        }
    }

    return nil
}
