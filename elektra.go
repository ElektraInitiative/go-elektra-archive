package elektra

// #cgo pkg-config: elektra-highlevel
// #include <elektra.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"github.com/pkg/errors"
)

type Elektra interface {
	Open(string) error
	Close()
}

type elektraC struct {
	handle *C.struct__Elektra
	error  *C.struct__ElektraError
}

// New returns a new Elektra instance.
func New() Elektra {
	return &elektraC{}
}

// Open creates a handle to the elektra library,
// this is mandatory to Get / Set Keys.
func (e *elektraC) Open(namespace string) error {
	n := C.CString(namespace)
	defer C.free(unsafe.Pointer(n))

	e.handle = C.elektraOpen(n, nil, &e.error)

	if e.error != nil {
		return errors.New("TODO")
	}

	return nil
}

// Close closes the elektra handle.
func (e *elektraC) Close() {
	C.elektraClose(e.handle)
}
