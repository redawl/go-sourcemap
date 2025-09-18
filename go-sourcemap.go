// go-sourcemap parses a sourcemap, and either prints the sourcemap to stdout, or saves the maps sources to a directory.
//
//	Usage:
//	    go-sourcemap [flags]
//	The flags are:
//	    -u
//	        Url to download the source map from. Cannot be specified at the same time as -f.
//	    -f
//	        File to read the source map from. Cannot be specified at the same time as -u.
//	    -d
//	        Directory to save decoded source files. If not specifed the decoded source map will be printed to stdout.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/redawl/go-sourcemap/spec"
	"github.com/redawl/go-sourcemap/tools"
)

type sourceMapArgs struct {
	url    string
	file   string
	outDir string
}

func main() {
	args := sourceMapArgs{}

	flag.StringVar(&args.url, "u", "", "url to download from")
	flag.StringVar(&args.file, "f", "", "path to location of sourcemap file")
	flag.StringVar(&args.outDir, "d", "", "Directory to save decoded source files. If not specified, decoded source map will be printed to stdout")

	flag.Parse()

	if args.url == "" && args.file == "" {
		fmt.Println("Either -u or -f is required")
		os.Exit(-1)
	}

	if args.url != "" && args.file != "" {
		fmt.Println("Cannot specify both -u and -f")
		os.Exit(-1)
	}

	var decoded *spec.DecodedSourceMapRecord
	var err error

	if args.url != "" {
		decoded, err = tools.ParseSourceMapFromUrl(args.url)
	} else if args.file != "" {
		decoded, err = tools.ParseSourceMapFromFile(args.file)
	}

	if err != nil {
		fmt.Printf("Error parsing source map from %s: %v\n", args.url, err)
		os.Exit(-1)
	}

	if args.outDir != "" {
		err := tools.SaveSourcesToDirectory(decoded, args.outDir)

		if err != nil {
			fmt.Printf("Error saving sources to %s: %v\n", args.outDir, err)
			os.Exit(-1)
		}
	} else {
		decodedStr, err := tools.MarshalDecodedSourceMapRecord(decoded)

		if err != nil {
			fmt.Printf("Error stringifying decodedStr: %v\n", err)
			os.Exit(-1)
		}

		fmt.Println(decodedStr)
	}
}
