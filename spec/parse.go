// Package spec contains types and implementations of the abstract methods defined in [Draft ECMA-426].
//
// [Draft ECMA-426]: https://tc39.es/ecma426/
package spec

import (
	"bytes"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"sync"
)

// Pre-compiled base64 lookup table for faster decoding
var base64Lookup [256]int

func init() {
	// Initialize lookup table with -1 (invalid)
	for i := range base64Lookup {
		base64Lookup[i] = -1
	}

	// Fill in valid base64 characters
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	for i, ch := range chars {
		base64Lookup[ch] = i
	}
}

// Buffer pool for reducing allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// ParseSourceMap parses str into a DecodedSourceMapRecord
// Returns an error if parsing was not successfull
// TODO: Support sourcemaps with the optional "sections" extension
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-ParseSourceMap
func ParseSourceMap(str string, baseURL string) (*DecodedSourceMapRecord, error) {
	sourceMap, err := ParseJSON(str)
	if err != nil {
		return nil, fmt.Errorf("Error parsing str: %w", err)
	}

	// TODO: call DecodeIndexSourceMap(sourceMap, baseURL) here if "sections" exists

	return DecodeSourceMap(sourceMap, baseURL)
}

// DecodeSourceMap decodes sourceMap into a DecodedSourceMapRecord.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-DecodeSourceMap
func DecodeSourceMap(sourceMap *SourceMap, baseURL string) (*DecodedSourceMapRecord, error) {
	if sourceMap.Version != 3 {
		slog.Warn("Version was not 3, parsing may fail", "version", sourceMap.Version)
	}

	ignoreList := sourceMap.IgnoreList
	if ignoreList == nil {
		// Check deprecanted x_google_ignore_list if ignoreList is null
		ignoreList = sourceMap.XGoogleIgnoreList
	}

	sources, err := DecodeSourceMapSources(baseURL, sourceMap.SourceRoot, sourceMap.Sources, sourceMap.SourcesContent, ignoreList)
	if err != nil {
		return nil, fmt.Errorf("Failed decoding source map sources: %w", err)
	}

	mappings, err := DecodeMappings(sourceMap.Mappings, sourceMap.Names, sources)
	if err != nil {
		return nil, fmt.Errorf("Error decoding mappings: %w", err)
	}

	return &DecodedSourceMapRecord{
		File:     sourceMap.File,
		Sources:  sources,
		Mappings: mappings,
	}, nil
}

// DecodeSourceMapSources decodes source map source information and returns a DecodedSourceRecord.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-DecodeSourceMapSources
func DecodeSourceMapSources(baseURL string, sourceRoot string, sources []string, sourcesContent []string, ignoreList []int) ([]*DecodedSourceRecord, error) {
	decodedSources := make([]*DecodedSourceRecord, len(sources))
	sourcesContentCount := len(sourcesContent)

	// Build sourceUrlPrefix more efficiently
	var sourceUrlPrefix strings.Builder
	if sourceRoot != "" {
		if idx := strings.Index(sourceRoot, "/"); idx != -1 {
			sourceUrlPrefix.WriteString(sourceRoot[:idx+1])
		} else {
			sourceUrlPrefix.WriteString(sourceRoot)
			sourceUrlPrefix.WriteByte('/')
		}
	}
	prefix := sourceUrlPrefix.String()

	// Convert ignoreList to map for O(1) lookup
	ignoreMap := make(map[int]bool, len(ignoreList))
	for _, idx := range ignoreList {
		ignoreMap[idx] = true
	}

	for index, source := range sources {
		decodedSource := &DecodedSourceRecord{
			Ignored: ignoreMap[index],
		}

		if source != "" {
			decodedSource.Url = baseURL + prefix + source
		}

		if sourcesContentCount > index {
			decodedSource.Content = sourcesContent[index]
		}

		decodedSources[index] = decodedSource
	}

	return decodedSources, nil
}

// DecodeMappings decodes mappings from a source map, and returns a slice of DecodedMappingRecords.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-DecodeMappings
func DecodeMappings(mappings string, names []string, sources []*DecodedSourceRecord) ([]*DecodedMappingRecord, error) {
	// Skip validation for performance (can be re-enabled if needed)
	// err := ValidateBase64VLQGroupings(mappings)
	// if err != nil {
	//     return nil, err
	// }

	// Pre-allocate with estimated capacity
	estimatedSize := strings.Count(mappings, ",") + strings.Count(mappings, ";")
	decodedMappings := make([]*DecodedMappingRecord, 0, estimatedSize)

	// Use byte slice for faster iteration
	mappingBytes := []byte(mappings)

	generatedLine := 0
	originalLine := 0
	originalColumn := 0
	sourceIndex := 0
	nameIndex := 0

	// Manual parsing instead of strings.Split to reduce allocations
	lineStart := 0
	for i := 0; i <= len(mappingBytes); i++ {
		if i == len(mappingBytes) || mappingBytes[i] == ';' {
			if i > lineStart {
				line := mappingBytes[lineStart:i]
				if len(line) > 0 {
					generatedColumn := 0
					segmentStart := 0

					for j := 0; j <= len(line); j++ {
						if j == len(line) || line[j] == ',' {
							if j > segmentStart {
								segment := string(line[segmentStart:j])
								position := 0

								relativeGeneratedColumn, err := DecodeBase64VLQ(segment, &position)
								if err != nil {
									return nil, fmt.Errorf("Error decoding base64 VLQ: %w", err)
								}

								generatedColumn += relativeGeneratedColumn

								if generatedColumn >= 0 {
									decodedMapping := &DecodedMappingRecord{
										GeneratedLine:   generatedLine,
										GeneratedColumn: generatedColumn,
									}

									decodedMappings = append(decodedMappings, decodedMapping)

									// Decode additional fields if present
									relativeSourceIndex, err := DecodeBase64VLQ(segment, &position)
									if err != nil {
										return nil, fmt.Errorf("Error decoding base64 VLQ: %w", err)
									}

									relativeOriginalLine, err := DecodeBase64VLQ(segment, &position)
									if err != nil {
										return nil, fmt.Errorf("Error decoding base64 VLQ: %w", err)
									}

									relativeOriginalColumn, err := DecodeBase64VLQ(segment, &position)
									if err != nil {
										return nil, fmt.Errorf("Error decoding base64 VLQ: %w", err)
									}

									if relativeOriginalColumn != math.MaxInt && relativeSourceIndex != math.MaxInt {
										sourceIndex += relativeSourceIndex
										originalLine += relativeOriginalLine
										originalColumn += relativeOriginalColumn

										if sourceIndex >= 0 && sourceIndex < len(sources) &&
											originalLine >= 0 && originalColumn >= 0 {
											decodedMapping.OriginalSource = sources[sourceIndex]
											decodedMapping.OriginalLine = originalLine
											decodedMapping.OriginalColumn = originalColumn
										} else {
											slog.Error(
												"Error: an index was less than zero, or sourceIndex >= len(sources)",
												"sourceIndex",
												sourceIndex,
												"originalLine",
												originalLine,
												"originalColumn",
												originalColumn,
												"len(sources)",
												len(sources),
											)
										}

										relativeNameIndex, err := DecodeBase64VLQ(segment, &position)
										if err != nil {
											return nil, fmt.Errorf("Error decoding base64 VLQ: %w", err)
										}

										if relativeNameIndex != math.MaxInt {
											nameIndex += relativeNameIndex
											if nameIndex >= 0 && nameIndex < len(names) {
												decodedMapping.Name = names[nameIndex]
											} else {
												slog.Error(
													"Error: nameIndex < 0 or nameIndex >= len(names)",
													"nameIndex",
													nameIndex,
													"len(names)",
													len(names),
												)
											}
										}
									}

									if int(position) != len(segment) {
										slog.Error("Error: position != len(segment)", "position", position, "len(segments)", len(segment))
									}
								} else {
									slog.Error("Error: generatedColumn was less than 0", "generatedColumn", generatedColumn)
								}
							}
							segmentStart = j + 1
						}
					}
				}
			}
			generatedLine++
			lineStart = i + 1
		}
	}

	return decodedMappings, nil
}

// ValidateBase64VLQGroupings validates that all chars in groupings are valid base64VLQ chars.
// Returns an error if any char in groupings is not a valid base64VLQ char.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-ValidateBase64VLQGroupings
func ValidateBase64VLQGroupings(groupings string) error {
	for i := 0; i < len(groupings); i++ {
		ch := groupings[i]
		if base64Lookup[ch] == -1 && ch != ',' && ch != ';' {
			return fmt.Errorf("Error: groupings contains invalid chars at position %d: %c", i, ch)
		}
	}
	return nil
}

// DecodeBase64VLQ attempts to base64VLQ decode the char in segment at position,
// and returns the int value of the decoded char.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#sec-DecodeBase64VLQ
func DecodeBase64VLQ(segment string, position *int) (int, error) {
	segmentLen := len(segment)
	if *position >= segmentLen {
		return math.MaxInt, nil
	}

	first, err := ConsumeBase64ValueAt(segment, position)
	if err != nil {
		return -1, err
	}

	if first >= 64 {
		panic(fmt.Sprintf("First >= 64: %d", first))
	}

	sign := 1
	if first%2 == 1 {
		sign = -1
	}

	value := (first % 32) / 2
	nextShift := 16
	currentByte := first

	for currentByte/32 == 1 {
		if *position >= segmentLen {
			return -1, fmt.Errorf("Error: position == segmentLen: %d", segmentLen)
		}

		currentByte, err = ConsumeBase64ValueAt(segment, position)
		if err != nil {
			return -1, err
		}

		chunk := currentByte % 32
		value += chunk * nextShift

		if value >= 2147483648 {
			return -1, fmt.Errorf("Error: value >= 2 ^ 31: %d", value)
		}

		nextShift *= 32
	}

	if value == 0 && sign == -1 {
		return -2147483648, nil
	}

	return sign * value, nil
}

// ConsumeBase64ValueAt attempts to base64 decode the char at position, and if successful returns the int value of the decoded char.
// Returns an error if position is out of bounds of str, or if the char at position is not a valid base64 char.
// If decoding of the char is successful, position is incremented.
// Uses lookup table instead of string search for better performance.
//
// [Source map format specification]
//
// [Source map format specification]: https://tc39.es/ecma426/#ConsumeBase64ValueAt
func ConsumeBase64ValueAt(str string, position *int) (int, error) {
	if *position >= len(str) {
		return -1, fmt.Errorf("Position was out of bounds of str!")
	}

	ch := str[*position]
	chIndex := base64Lookup[ch]

	if chIndex == -1 {
		return -1, fmt.Errorf("Invalid base64 char: %c", ch)
	}

	*position++
	return chIndex, nil
}
