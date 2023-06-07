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

	f := fox.New(
		fox.WithMiddleware(foxdump.Middleware(DumpRequest, DumpResponse)),
	)

	f.MustHandle(http.MethodPost, "/hello/fox", func(c fox.Context) {
		_, _ = io.Copy(c.Writer(), c.Request().Body)
	})

	log.Fatalln(http.ListenAndServe(":8080", f))
}

````

Note that the `buf` slice is transient, and its data is only guaranteed to be valid during the execution of the 
`foxdump.BodyHandler` function. If the data needs to be persisted or used outside the scope of this function, it should be copied 
to a new byte slice (e.g., using `copy`). Furthermore, `buf` should be treated as read-only to prevent any unintended 
side effects.