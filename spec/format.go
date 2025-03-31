package spec

// SourceMap represents the raw json map, without any parsing
type SourceMap struct {
    // Version must always be 3
    Version int             `json:"version"`             
    // File is the *optional* name of the compiled output i.e. *.map.js
    File string             `json:"file"`                
    // SourceRoot is optional 
    SourceRoot string       `json:"sourceRoot"`          
    // Sources is the original mapped sources names
    Sources []string        `json:"sources"`             
    // SourcesContent is the original mapped sources contents
    SourcesContent []string `json:"sourcesContent"`      
    // Name is the optional symbol names which can be used by mappings field
    Names []string          `json:"names"`               
    // Mappings is the encoded mapping data
    Mappings string         `json:"mappings"`            
    // IgnoreList is an optional list of indices that should be considered third-party code
    IgnoreList []int        `json:"ignoreList"`          
    // Deprecated: XGoogleIgnoreList is only checked if IgnoreList is not present
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

// DecodedMappingRecord represents a mapping of a source symbol to the 
// generated symbol.
type DecodedMappingRecord struct {
    GeneratedLine int                   `json:"generatedLine"`
    GeneratedColumn int                 `json:"generatedColumn"`
    OriginalSource *DecodedSourceRecord `json:"originalSource"`
    OriginalLine int                    `json:"originalLine"`
    OriginalColumn int                  `json:"originalColumn"`
    Name string                         `json:"name"`
}

// DecodedSourceMapRecord represents a fully decoded source map record.
type DecodedSourceMapRecord struct {
    // File is a *optional* name of the compiled output i.e. *.map.js
    File string                      `json:"file"`
    // Sources is the original source records
    Sources []*DecodedSourceRecord   `json:"sources"`
    // Mappings is the symbol mappings from source records to compuled output map record
    Mappings []*DecodedMappingRecord `json:"mappings"`
}

