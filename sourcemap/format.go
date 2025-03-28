package sourcemap

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type SourceMap struct {
    Version int             `json:"version"`             // must always be 3
    File string             `json:"file"`                // optional, name of the compiled output i.e. *.map.js
    SourceRoot string       `json:"sourceRoot"`          // optional
    Sources []string        `json:"sources"`             // original mapped sources names
    SourcesContent []string `json:"sourcesContent"`      // original mapped sources contents
    Names []string          `json:"names"`               // optional, symbol names which can be used by mappings field
    Mappings string         `json:"mappings"`            // Encoded mapping data
    IgnoreList []int        `json:"ingoreList"`          // optional, list of indices that should be considered third-party code
    XGoogleIgnoreList []int `json:"x_google_ignoreList"` // Deprecated, only checked if ignoreList is not present
}

type SourceRecord struct {
    Url string
    Content string
    Ingored bool
}

type Source struct {
    sourceMap SourceMap
    isParsed bool
    sourceRecord SourceRecord
}

func NewSource(sourceMapText string) (*Source, error) {
    decoder := json.NewDecoder(strings.NewReader(sourceMapText))

    source := Source{}
    err := decoder.Decode(source)

    if err != nil {
        return nil, fmt.Errorf("Failed to parse source map, %v", err)
        
    }

    return &source, nil
}
