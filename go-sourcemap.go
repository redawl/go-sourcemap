package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/redawl/go-sourcemap/sourcemap/spec"
)

func main() {
    url := flag.String("u", "http://localhost:80", "url to download from")
    generateDirectory := flag.String("d", "", "directory to parse sources to")
    flag.Parse()

    if *generateDirectory == "" {
        slog.Error("-d (output directory) is required")
        os.Exit(-1)
    }

    sourceMapFiles := flag.Args()

    for _, mapFile := range sourceMapFiles {
        uri := fmt.Sprintf("%s/%s", *url, mapFile)
        var body string
        if strings.Contains(*url, "http") {
            resp, err := http.DefaultClient.Get(uri)

            if err != nil {
                panic(err)
            }

            bodyBytes, err := io.ReadAll(resp.Body)

            if err != nil {
                panic(err)
            }

            fmt.Printf("Status code: %s\n", resp.Status)
            body = string(bodyBytes)
        } else {
            bodyBytes, err := os.ReadFile(mapFile)

            if err != nil {
                slog.Error("Error opening raw sourceMap file", "mapFile", mapFile)
            }

            body = string(bodyBytes)
        }

        decoded, err := spec.ParseSourceMap(string(body), *url)
        
        if err != nil {
            slog.Error("Error parsing source map", "error", err)
            os.Exit(-1)
        }

        if false {
            err = os.MkdirAll(*generateDirectory, 0700)

            if err != nil {
                slog.Error("Error creating directory", "directory", *generateDirectory)
                os.Exit(-1)
            }

            for _, source := range decoded.Sources {
                index := strings.LastIndexByte(source.Url, os.PathSeparator)

                if index != -1 {
                    err = os.MkdirAll(*generateDirectory + source.Url[:index], 0700)

                    if err != nil {
                        slog.Error("Error creating directory", "directory", *generateDirectory + source.Url[:index])
                        os.Exit(-1)
                    }

                    err := os.WriteFile(*generateDirectory + source.Url, []byte(source.Content), 0600)

                    if err != nil {
                        slog.Error("Error writing file contents", "path",*generateDirectory + source.Url)
                        os.Exit(-1)
                    }
                }
            }

            slog.Info("Successfully wrote source files", "count", len(decoded.Sources))
        }

    }
}

