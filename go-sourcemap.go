package main

import (
	"flag"
	"log/slog"
	"os"
)

func main() {
    _ = flag.String("u", "http://localhost:80", "url to download from")
    generateDirectory := flag.String("d", "", "directory to parse sources to")
    flag.Parse()

    if *generateDirectory == "" {
        slog.Error("-d (output directory) is required")
        os.Exit(-1)
    }

    _ = flag.Args()

}

