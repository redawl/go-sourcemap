package spec

// SourceMap represents the raw json map, without any parsing
type SourceMap struct {
    // Version must always be 3
    Version int             `json:"version"`             
    // File is the *optional* name of the compiled output i.e. *.map.js
    File string             `json:"file"`                
    // SourceRoot is optional
    SourceRoot string       `json:"sourceRoot"`          
    // original mapped sources names
    Sources []string        `json:"sources"`             
    // original mapped sources contents
    SourcesContent []string `json:"sourcesContent"`      
    // optional, symbol names which can be used by mappings field
    Names []string          `json:"names"`               
    // Encoded mapping data
    Mappings string         `json:"mappings"`            
    // optional, list of indices that should be considered third-party code
    IgnoreList []int        `json:"ingoreList"`          
    // Deprecated, only checked if ignoreList is not present
    XGoogleIgnoreList []int `json:"x_google_ignoreList"` 
}

// DecodedSourceRecord represents an original source file 
// which was used to generate the compiled output.
type DecodedSourceRecord struct {
    // Url is the filepath of the source file
    Url string      `json:"url"`
    // Content is the contents of the source file
    Content string  `json:"content"`
    // Ignored is whether the source file should be ignored by analyzers
    Ignored bool    `json:"ignored"`
}

type DecodedMappingRecord struct {
    GeneratedLine int            `json:"generatedLine"`
    GeneratedColumn int          `json:"generatedColumn"`
    OriginalSource *DecodedSourceRecord `json:"originalSource"`
    OriginalLine int             `json:"originalLine"`
    OriginalColumn int           `json:"originalColumn"`
    Name string                  `json:"name"`
}

type DecodedSourceMapRecord struct {
    File string                      `json:"file"`
    Sources []*DecodedSourceRecord          `json:"sources"`
    Mappings []*DecodedMappingRecord `json:"mappings"`
}

