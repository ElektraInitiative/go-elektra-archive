package elektra

// #cgo pkg-config: elektra-highlevel
// #include <elektra.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/pkg/errors"
)

type Elektra interface {
	Open(string) error
	Close()

	Value(name string) string
	Long(name string) int64
}

type elektraC struct {
	handle *C.struct__Elektra
	err    *C.struct__ElektraError
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

	e.handle = C.elektraOpen(n, nil, &e.err)

	if err := e.lastError(); err != nil {
		return err
	}

	C.elektraFatalErrorHandler(e.handle, errCallback)

	return nil
}

func errCallback(err *C.struct__ElektraError) {
	errDescription := C.GoString(C.elektraErrorDescription(err))

	fmt.Printf(errDescription)
}

func (e *elektraC) lastError() error {
	if e.err == nil {
		return nil
	}

	errDescription := C.GoString(C.elektraErrorDescription(e.err))

	C.elektraErrorReset(&e.err)

	return errors.New(errDescription)
}

// Close closes the elektra handle.
func (e *elektraC) Close() {
	C.elektraClose(e.handle)
}

func (e *elektraC) Value(name string) string {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	val := C.elektraGetString(e.handle, n)

	return C.GoString(val)
}

func (e *elektraC) Long(name string) int64 {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	val := C.elektraGetLong(e.handle, n)

	return int64(val)
}
