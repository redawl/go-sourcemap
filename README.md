# go-sourcemap
[![Go Reference](https://pkg.go.dev/badge/github.com/redawl/go-sourcemap.svg)](https://pkg.go.dev/github.com/redawl/go-sourcemap)
![tests badge](https://github.com/redawl/go-sourcemap/actions/workflows/tests.yml/badge.svg)

go-sourcemap is a pure-go implementation of the [Source Map](https://tc39.es/ecma426/) specification.

```bash
user@workstation ~ $ go-sourcemap -h
Usage of go-sourcemap:
  -d string
        Directory to save decoded source files. If not specified, decoded source map will be printed to stdout
  -f string
        path to location of sourcemap file
  -u string
        url to download from
```
