# go bindings for Elektra

This repository contains the low level ("kdb" subpackage) and high level (root package, TODO) go bindings for the elektra library.

## Execute tests
Prerequisite: you have to have KDB and go installed on your machine.

Execute all tests:
`go test ./...`

Execute tests of a package, e.g. kdb:
`go test ./kdb`

## Use Elektra in your application
Just *go get* it like you are used to with go.

`go get github.com/ElektraInitiative/go-elektra`

In the future we will add a vanity import.

To use it import it in your .go file (error handling was omitted for brevity):

```go
import "github.com/ElektraInitiative/go-elektra/kdb"

func main() {
    parentKey, _ := kdb.CreateKey("user/test")
	ks, _ := kdb.CreateKeySet()

    handle := kdb.New()
	_ = handle.Open(parentKey)
    _ = handle.Get(ks, parentKey)
    
    foundKey := ks.LookupByName("/test/hello_world")
    
    value := foundKey.Value()

    // do something with the value
}
```
