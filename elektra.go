package elektra

/*
#cgo pkg-config: elektra-highlevel
#include <elektra.h>
#include <stdlib.h>

void err_callback(ElektraError *err);

static void _register_callback(Elektra * elektra) {
	elektraFatalErrorHandler(elektra, err_callback);
}
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/pkg/errors"
)

type Elektra interface {
	Open(string) error
	Close()

	SetValue(name, value string) error
	Value(name string) string
	SetLong(name string, value int64) error
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


// Open creates a handle to the Elektra library,
// this is mandatory to Get / Set Keys.
func (e *elektraC) Open(application string) error {
	n := C.CString(application)
	defer C.free(unsafe.Pointer(n))

	e.handle = C.elektraOpen(n, nil, &e.err)

	if err := e.lastError(); err != nil {
		return err
	}

	C._register_callback(e.handle)

	return nil
}

//export err_callback
func err_callback(err *C.struct__ElektraError) {
	errDescription := C.GoString(C.elektraErrorDescription(err))

	fmt.Printf("elektra err: %s", errDescription)
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

func (e *elektraC) SetValue(name, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	C.elektraSetString(e.handle, n, v, &e.err)

	return e.lastError()
}

func (e *elektraC) Value(name string) string {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	val := C.elektraGetString(e.handle, n)

	return C.GoString(val)
}

func (e *elektraC) SetLong(name string, value int64) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	C.elektraSetLong(e.handle, n, C.int(value), &e.err)

	return e.lastError()
}

func (e *elektraC) Long(name string) int64 {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	val := C.elektraGetLong(e.handle, n)

	return int64(val)
}
