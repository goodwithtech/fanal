# fanal
Static Analysis Library for Containers

[![GoDoc](https://godoc.org/github.com/goodwithtech/deckoder?status.svg)](https://godoc.org/github.com/goodwithtech/deckoder)
[![Build Status](https://travis-ci.org/goodwithtech/deckoder.svg?branch=master)](https://travis-ci.org/goodwithtech/deckoder)
<!-- [![Coverage Status](https://coveralls.io/repos/github/goodwithtech/deckoder/badge.svg?branch=master)](https://coveralls.io/github/goodwithtech/deckoder?branch=master) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/goodwithtech/deckoder)](https://goreportcard.com/report/github.com/goodwithtech/deckoder)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/goodwithtech/deckoder/blob/master/LICENSE)

## Feature
- Detect OS
- Extract OS packages
- Extract libraries used by an application
  - Bundler, Composer, npm, Yarn, Pipenv, Poetry, Cargo

## Example
See [`cmd/fanal/`](cmd/fanal)

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/xerrors"

	"github.com/goodwithtech/deckoder/cache"

	"github.com/goodwithtech/deckoder/analyzer"
	_ "github.com/goodwithtech/deckoder/analyzer/library/bundler"
	_ "github.com/goodwithtech/deckoder/analyzer/library/composer"
	_ "github.com/goodwithtech/deckoder/analyzer/library/npm"
	_ "github.com/goodwithtech/deckoder/analyzer/library/pipenv"
	_ "github.com/goodwithtech/deckoder/analyzer/library/poetry"
	_ "github.com/goodwithtech/deckoder/analyzer/library/yarn"
	_ "github.com/goodwithtech/deckoder/analyzer/library/cargo"
	_ "github.com/goodwithtech/deckoder/analyzer/os/alpine"
	_ "github.com/goodwithtech/deckoder/analyzer/os/amazonlinux"
	_ "github.com/goodwithtech/deckoder/analyzer/os/debianbase"
	_ "github.com/goodwithtech/deckoder/analyzer/os/suse"
	_ "github.com/goodwithtech/deckoder/analyzer/os/redhatbase"
	_ "github.com/goodwithtech/deckoder/analyzer/pkg/apk"
	_ "github.com/goodwithtech/deckoder/analyzer/pkg/dpkg"
	_ "github.com/goodwithtech/deckoder/analyzer/pkg/rpm"
	"github.com/goodwithtech/deckoder/extractor"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	ctx := context.Background()
	tarPath := flag.String("f", "-", "layer.tar path")
	clearCache := flag.Bool("clear", false, "clear cache")
	flag.Parse()

	if *clearCache {
		if err = cache.Clear(); err != nil {
			return xerrors.Errorf("error in cache clear: %w", err)
		}
	}

	args := flag.Args()

	var files extractor.FileMap
	if len(args) > 0 {
		files, err = analyzer.Analyze(ctx, args[0])
		if err != nil {
			return err
		}
	} else {
		rc, err := openStream(*tarPath)
		if err != nil {
			return err
		}

		files, err = analyzer.AnalyzeFromFile(ctx, rc)
		if err != nil {
			return err
		}
	}

	os, err := analyzer.GetOS(files)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", os)

	pkgs, err := analyzer.GetPackages(files)
	if err != nil {
		return err
	}
	fmt.Printf("Packages: %d\n", len(pkgs))

	libs, err := analyzer.GetLibraries(files)
	if err != nil {
		return err
	}
	for filepath, libList := range libs {
		fmt.Printf("%s: %d\n", filepath, len(libList))
	}
	return nil
}

func openStream(path string) (*os.File, error) {
	if path == "-" {
		if terminal.IsTerminal(0) {
			flag.Usage()
			os.Exit(64)
		} else {
			return os.Stdin, nil
		}
	}
	return os.Open(path)
}

```


## Notes
When using `latest` tag, that image will be cached. After `latest` tag is updated, you need to clear cache.



