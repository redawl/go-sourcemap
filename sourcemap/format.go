package sourcemap

type SourceMap struct {
    Version int             `json:"version"`             // must always be 3
    File string             `json:"file"`                // optional, name of the compiled output i.e. *.map.js
    SourceRoot string       `json:"sourceRoot"`          // optional
    Sources []string        `json:"sources"`             // original mapped sources names
    SourcesContent []string `json:"sourcesContent"`      // original mapped sources contents
    Names []string          `json:"names"`               // optional, symbol names which can be used by mappings field
    Mappings string         `json:"mappings"`            // Encoded mapping data
    IgnoreList []uint        `json:"ingoreList"`          // optional, list of indices that should be considered third-party code
    XGoogleIgnoreList []uint `json:"x_google_ignoreList"` // Deprecated, only checked if ignoreList is not present
}

type SourceRecord struct {
    Url string
    Content string
    Ignored bool
}

type DecodedMappingRecord struct {
    GeneratedLine uint
    GeneratedColumn uint
    OriginalSource *SourceRecord
    OriginalLine uint
    OriginalColumn uint
    Name string
}

type DecodedSourceMapRecord struct {
    File string
    Sources []*SourceRecord
    Mappings []*DecodedMappingRecord
}
