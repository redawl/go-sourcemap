package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {
    url := flag.String("u", "http://localhost:80", "url to download from")
    flag.Parse()
    sourceMapFiles := flag.Args()

    for _, mapFile := range sourceMapFiles {
        resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/%s", *url, mapFile))

        if err != nil {
            panic(err)
        }

        body, err := io.ReadAll(resp.Body)

        if err != nil {
            panic(err)
        }

        fmt.Printf("Status code: %s\n", resp.Status)
        fmt.Println(string(body))
    }
}
