# Go Bindings for Elektra

This repository contains the low-level ("kdb" subpackage).

# Prerequisites

* Go (version >1.11) and
* libelektra installed.

## Build

Run `go install` or `go build`.

## Execute Tests

Prerequisite: you have to have KDB and Go installed on your machine.

Execute all tests:
`go test ./...`

Execute tests of a package, e.g. kdb:
`go test ./kdb`

## Use Elektra in your Application

Just _go get_ it like you are used to with Go.

`go get github.com/ElektraInitiative/go-elektra`

In the future we will add a vanity import.

To use it import it in your .go file (error handling was omitted for brevity):

```go
package main

import (
    "fmt"

    "github.com/ElektraInitiative/go-elektra/kdb"
)

func main() {
    // PREREQUISITE: run `kdb set /test/hello_world foo` in your terminal
	ks, _ := kdb.CreateKeySet()

    handle := kdb.New()
    _, _ = handle.Open()

    parentKey, _ := kdb.CreateKey("user/test")
    _, _ = handle.Get(ks, parentKey)

    foundKey := ks.LookupByName("/test/hello_world")

    value := foundKey.Value()

    fmt.Print(value)
}
```

## Documentation

The documentation can be viewed on [godoc.org](https://godoc.org/github.com/ElektraInitiative/go-elektra/kdb)

## Troubleshooting

### Elektra-Go does not compile

Make sure that libelektra is installed.

Go-Elektra leverages [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) to compile the Elektra library.

If the bindings fail to compile you probably need to set the `PKG_CONFIG_PATH` to the installation folder of Elektra, e.g.: `PKG_CONFIG_PATH=/usr/local/lib/pkgconfig`.

