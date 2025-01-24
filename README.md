[![Go Reference](https://pkg.go.dev/badge/github.com/tigerwill90/foxdump.svg)](https://pkg.go.dev/github.com/tigerwill90/foxdump)
[![tests](https://github.com/tigerwill90/foxdump/actions/workflows/tests.yaml/badge.svg)](https://github.com/tigerwill90/foxdump/actions?query=workflow%3Atests)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigerwill90/foxdump)](https://goreportcard.com/report/github.com/tigerwill90/foxdump)
[![codecov](https://codecov.io/gh/tigerwill90/foxdump/branch/master/graph/badge.svg?token=D6qSTlzEcE)](https://codecov.io/gh/tigerwill90/foxdump)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tigerwill90/foxdump)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tigerwill90/foxdump)

# Foxdump
Foxdump is a middleware for [Fox](https://github.com/tigerwill90/fox) that provides an efficient way to dump 
HTTP request and response bodies. This feature can be extremely useful for debugging, logging, testing, and 
monitoring HTTP traffic.

## Disclaimer
Foxdump's API is linked to Fox router, and it will only reach v1 when the router is stabilized.
During the pre-v1 phase, breaking changes may occur and will be documented in the release notes.

## Getting started
### Installation
````shell
go get -u github.com/tigerwill90/foxdump
````

### Features
- Efficient body dumping with minimal performance impact (zero allocation).
- Can be configured to dump either request, response, or both bodies.
- Easily integrate with Fox ecosystem.

### Usage

Here's a simple example of how to use the Foxdump middleware:
````go
package main

import (
	"errors"
	"github.com/tigerwill90/fox"
	"github.com/tigerwill90/foxdump"
	"io"
	"log"
	"net/http"
)

func DumpRequest(c fox.Context, buf []byte) {
	log.Println("request:", string(buf))
}

func DumpResponse(c fox.Context, buf []byte) {
	log.Println("response:", string(buf))
}

func main() {
	f, err := fox.New(
		fox.WithMiddleware(foxdump.Middleware(DumpRequest, DumpResponse)),
	)
	if err != nil {
		panic(err)
	}

	f.MustHandle(http.MethodPost, "/hello/fox", func(c fox.Context) {
		_, _ = io.Copy(c.Writer(), c.Request().Body)
	})

	if err = http.ListenAndServe(":8080", f); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
````

Note that the `buf` slice is transient, and its data is only guaranteed to be valid during the execution of the 
`foxdump.BodyHandler` function. If the data needs to be persisted or used outside the scope of this function, it should be copied 
to a new byte slice (e.g. using `copy`). Furthermore, `buf` should be treated as read-only to prevent any unintended 
side effects.

## Benchmark
````
goos: linux
goarch: amd64
pkg: github.com/tigerwill90/foxdump
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkFoxDumpMiddleware-16            7333161             157.9 ns/op               0 B/op          0 allocs/op
BenchmarkEchoDumpMiddleware-16           1215985              1618 ns/op            2665 B/op         10 allocs/op
````
