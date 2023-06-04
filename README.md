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
	"github.com/tigerwill90/foxdump"
	"github.com/tigerwill90/fox"
	"log"
	"net/http"
)

func main() {
	f := fox.New(
		fox.WithMiddleware(foxdump.Middleware(func(c fox.Context, req, res []byte) {
			log.Printf("Request Body: %s\n", string(req))
			log.Printf("Response Body: %s\n", string(res))
		})),
	)

	f.MustHandle(http.MethodGet, "/hello/{name}", func(c fox.Context) {
		_ = c.String(http.StatusOK, "hello %s", c.Param("name"))
	})

	log.Fatalln(http.ListenAndServe(":8080", f))
}
````

Note that the 'req' and 'res' byte slices are transient, and their data is only guaranteed to be valid
during the execution of the BodyDumpHandler function. If the data needs to be persisted or
used outside the scope of this function, it should be copied to a new byte slice (e.g., using 'copy').
Furthermore, these slices should be treated as read-only to prevent any unintended side effects.