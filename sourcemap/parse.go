package sourcemap

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

func ParseSourceMap(str string, baseURL string) (*DecodedSourceMapRecord, error) {
    sourceMap, err := ParseJSON(str)

    if err != nil {
        return nil, fmt.Errorf("Error parsing str: %s", err.Error())
    }

    // TODO: call DecodeIndexSourceMap(sourceMap, baseURL) here

    return DecodeSourceMap(sourceMap, baseURL)
}

func ParseJSON(str string) (*SourceMap, error) {
    decoder := json.NewDecoder(strings.NewReader(str))

    sourceMap := &SourceMap{}

    err := decoder.Decode(sourceMap)

    if err != nil {
        return nil, err
    }
    
    return sourceMap, nil
}

func DecodeSourceMap(sourceMap *SourceMap, baseURL string) (*DecodedSourceMapRecord, error) {
    if sourceMap.Version != 3 {
        slog.Warn("Version was not 3, parsing may fail", "version", sourceMap.Version)
    }

    sources, err := DecodeSourceMapSources(baseURL, sourceMap.SourceRoot, sourceMap.Sources, sourceMap.SourcesContent, sourceMap.IgnoreList)

    if err != nil {
        return nil, fmt.Errorf("Failed decoding source map sources: %s", err.Error())
    }

    mappings, err := DecodeMappings(sourceMap.Mappings, sourceMap.Names, sources)

    if err != nil {
        return nil, fmt.Errorf("Error decoding mappings: %s", err.Error())
    }

    return &DecodedSourceMapRecord{
        File: sourceMap.File,
        Sources: sources,
        Mappings: mappings,
    }, nil
}

func DecodeSourceMapSources(baseURL string, sourceRoot string, sources []string, sourcesContent []string, ignoreList []uint) ([]*SourceRecord, error) {
    decodedSources := make([]*SourceRecord, len(sources))

    sourcesContentCount := len(sourcesContent)

    var sourceUrlPrefix string

    if sourceRoot != "" {
        if strings.Contains(sourceRoot, "\x2F") {
            idx := strings.Index(sourceRoot, "\x2F") 

            sourceUrlPrefix = sourceRoot[0:idx+1]
        } else {
            sourceUrlPrefix = sourceRoot + "/" 
        }
    }

    for index, source := range sources {
        decodedSource := &SourceRecord{
            Ignored: false,
        }

        if source != "" {
            decodedSource.Url = baseURL + sourceUrlPrefix + source
        }

        if slices.Contains(ignoreList, uint(index)) {
            decodedSource.Ignored = true 
        }

        if sourcesContentCount > index {
           decodedSource.Content = sourcesContent[index]  
        }

        decodedSources[index] = decodedSource
    }

    return decodedSources, nil
}

func DecodeMappings(mappings string, names []string, sources []*SourceRecord) ([]*DecodedMappingRecord, error) {
    err := ValidateBase64VLQGroupings(mappings)

    if err != nil {
        return nil, err
    }

    decodedMappings := make([]*DecodedMappingRecord, 0)

    groups := strings.Split(mappings, ";")

    generatedLine := 0
    originalLine := 0
    originalColumn := 0
    sourceIndex := 0
    nameIndex := 0

    for generatedLine < len(groups) {
        if groups[generatedLine] != "" {
            segments := strings.Split(groups[generatedLine], ",")

            generatedColumn := 0
            for _, segment := range(segments) {
                position := uint(0)
                // TODO: position needs to be a reference so that its value can be updated with each call to DecodeBase64VLQ
                relativeGeneratedColumn, err := DecodeBase64VLQ(segment, &position)

                if err != nil {
                    return nil, fmt.Errorf("Error decoding base64 VLQ: %s", err.Error())
                } else {
                    generatedColumn = generatedColumn + relativeGeneratedColumn

                    if generatedColumn < 0 {
                        slog.Error("Error: generatedColumn was less than 0", "generatedColumn", 0)
                    } else {
                        decodedMapping := &DecodedMappingRecord{
                            GeneratedLine: uint(generatedLine),
                            GeneratedColumn: uint(generatedColumn),
                        }

                        decodedMappings = append(decodedMappings, decodedMapping)

                        relativeSourceIndex, err := DecodeBase64VLQ(segment, &position)

                        if err != nil {
                            return nil, fmt.Errorf("Error decoding base64 VLQ: %s", err.Error())
                        }

                        relativeOriginalLine, err := DecodeBase64VLQ(segment, &position)

                        if err != nil {
                            return nil, fmt.Errorf("Error decoding base64 VLQ: %s", err.Error())
                        }

                        relativeOriginalColumn, err := DecodeBase64VLQ(segment, &position)

                        if err != nil {
                            return nil, fmt.Errorf("Error decoding base64 VLQ: %s", err.Error())
                        }
                        
                        if relativeOriginalColumn == -1 && relativeSourceIndex != -1 {
                            slog.Error("Error: relativeOriginalColumn was -1 when relativeSourceIndex was not -1", "relativeSourceIndex", relativeSourceIndex)
                        } else if relativeOriginalColumn != -1 {
                            sourceIndex += relativeSourceIndex
                            originalLine += relativeOriginalLine
                            originalColumn += relativeOriginalColumn
                            
                            // TODO: Docs says source, but that doesn't exist. Is this correct?
                            if sourceIndex < 0 || originalLine < 0 || originalColumn < 0 || sourceIndex >= len(sources) {
                                 slog.Error("Error: an index was less than zero, or sourceIndex >= len(sources)", "sourceIndex", sourceIndex, "originalLine", originalLine, "originalColumn", originalColumn, "len(sources)", len(sources))
                            } else {
                                decodedMapping.OriginalSource = sources[sourceIndex]
                                decodedMapping.OriginalLine = uint(originalLine)
                                decodedMapping.OriginalColumn = uint(originalColumn) 
                            }

                            relativeNameIndex, err := DecodeBase64VLQ(segment, &position)

                            if err != nil {
                                return nil, fmt.Errorf("Error decoding base64 VLQ: %s", err.Error())
                            }
                            
                            if relativeNameIndex != -1 {
                                nameIndex += relativeNameIndex
                                if nameIndex < 0 || nameIndex >= len(names) {
                                    slog.Error("Error: nameIndex < 0 or nameIndex >= len(names)", "nameIndex", nameIndex, "len(names)", len(names))
                                } else {
                                    decodedMapping.Name = names[nameIndex]
                                }
                            }
                        }

                        if int(position) != len(segment) {
                            slog.Error("Error: position != len(segment)", "position", position, "len(segments)", len(segment))
                        }
                    }
                }
            }
        }

        generatedLine++
    }

    return decodedMappings, nil
}

func ValidateBase64VLQGroupings(groupings string) error {
    if strings.ContainsAny(groupings, "")  {
        return fmt.Errorf("Error: groupings contains invalid chars: %s", groupings)
    }

    return nil 
}

func DecodeBase64VLQ(segment string, position *uint) (int, error) {
    segmentLen := len(segment)
    if int(*position) == segmentLen {
        return -1, nil
    }

    first, err := ConsumeBase64ValueAt(segment, position)

    if err != nil {
        return -1, err
    }

    if first >= 64 {
        // Panic, since ConsumeBase64ValueAt should return error for this condition
        panic(fmt.Sprintf("First >= 64: %d", first))
    }

    sign := 0

    if first % 2 == 0 {
        sign = 1
    } else {
        sign = -1
    }

    value := (first % 32) / 2

    nextShift := 16
    currentByte := first

    for currentByte / 32 == 1 {
        if *position == uint(segmentLen) {
            return -1, fmt.Errorf("Error: position == segmentLen => %d == %d", *position, segmentLen)
        }

        currentByte, err = ConsumeBase64ValueAt(segment, position)

        if err != nil {
            return -1, err
        }

        chunk := currentByte % 32

        value += chunk * nextShift

        if value >= 2147483648 {
            return -1 , fmt.Errorf("Error: value >= 2 ^ 31: %d", value)
        }

        nextShift *= 32
    }

    if value == 0 && sign == -1 {
        return -2147483648, nil
    }
    
    slog.Info("Returning from DecodeBase64VLQ", "segment", segment, "ret", sign * value)
    return sign * value, nil
}

// ConsumeBase64ValueAt attempts to base64 decode the char at position, and if successful returns the int value of the decoded char. 
// Returns an error if position is out of bounds of str, or if the char at position is not a valid base64 char
func ConsumeBase64ValueAt(str string, position *uint) (int, error) {
    alph := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

    if int(*position) >= len(str) {
        return -1, fmt.Errorf("Position was out of bounds of str!")
    }

    ch := str[*position]
    chIndex := strings.IndexByte(alph, ch)

    if chIndex == -1 {
        return -1, fmt.Errorf("Invalid base64 char: %c", ch)
    }

    *position++

    return chIndex, nil
}
    
