package spec

import (
	"os"
	"testing"
)

var benchmarkFiles = map[string]string{
	"Tiny":   "test1.js.map",         // 348 bytes
	"Small":  "test2.js.map",         // 459 bytes
	"Medium": "jquery.min.map",       // 135 KB
	"Large":  "angular-core.mjs.map", // 2.7 MB
	"XLarge": "babylon.js.map",       // 18 MB
}

func loadBenchmarkFile(filename string) string {
	contents, err := os.ReadFile("../testdata/" + filename)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

func BenchmarkParseSourceMap(b *testing.B) {
	for name, filename := range benchmarkFiles {
		contents := loadBenchmarkFile(filename)

		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(contents)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := ParseSourceMap(contents, "")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseJSON(b *testing.B) {
	for name, filename := range benchmarkFiles {
		contents := loadBenchmarkFile(filename)

		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(contents)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := ParseJSON(contents)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeSourceMap(b *testing.B) {
	for name, filename := range benchmarkFiles {
		contents := loadBenchmarkFile(filename)
		sourceMap, err := ParseJSON(contents)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := DecodeSourceMap(sourceMap, "")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeMappings(b *testing.B) {
	for name, filename := range benchmarkFiles {
		contents := loadBenchmarkFile(filename)
		sourceMap, err := ParseJSON(contents)
		if err != nil {
			b.Fatal(err)
		}

		sources, err := DecodeSourceMapSources("", sourceMap.SourceRoot, sourceMap.Sources, sourceMap.SourcesContent, sourceMap.IgnoreList)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := DecodeMappings(sourceMap.Mappings, sourceMap.Names, sources)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
