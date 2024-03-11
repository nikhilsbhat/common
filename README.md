# Common Library

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/common)](https://goreportcard.com/report/github.com/nikhilsbhat/common)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/common/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/common?status.svg)](https://godoc.org/github.com/nikhilsbhat/common)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/common.svg)](https://github.com/nikhilsbhat/common/tags)

## Introduction

Generic functions that would be handy while building your Golang utility.

This library stands on the shoulders of various libraries built by some awesome folks.

## Installation

Get the latest version of GoCD sdk using `go get` command.

```shell
go get github.com/nikhilsbhat/common@latest
```

Get specific version of the same.

```shell
go get github.com/nikhilsbhat/common@v0.0.2
```

## Usage

```go
package main

import (
	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type Object struct {
	Name string
	Date string
}

func main() {
	newObject := []Object{
		{Name: "nikhil", Date: "01-01-2024"},
		{Name: "john", Date: "01-02-2024"},
	}

	logger := logrus.New()
	render := renderer.GetRenderer(os.Stdout, logger, true, true, false, false, false)

	if err := render.Render(newObject); err != nil {
		log.Fatal(err)
	}
}
```

Above code should generate yaml as below:

```yaml
---
- Date: 01-01-2024
  Name: nikhil
- Date: 01-02-2024
  Name: john
```

More example of the libraries can be found [here](https://github.com/nikhilsbhat/common/blob/main/example).
