package kdb

// #cgo pkg-config: elektra
// #include <kdb.h>
import "C"

import (
	"github.com/pkg/errors"
)

// KDB is an interface to the Elektra library.
type KDB interface {
	Open() error
	Close() error

	Get(keySet KeySet, parentKey Key) (changed bool, err error)
	Set(keySet KeySet, parentKey Key) (changed bool, err error)

	Version() (string, error)
}

type kdbC struct {
	handle *C.struct__KDB
}

// New returns a new KDB instance.
func New() KDB {
	return &kdbC{}
}

// Open creates a handle to the kdb library,
// this is mandatory to Get / Set Keys.
func (e *kdbC) Open() error {
	key, err := createKey("")

	if err != nil {
		return err
	}

	handle := C.kdbOpen(key.ptr)

	if handle == nil {
		return errFromKey(key)
	}

	e.handle = handle

	return nil
}

// Close closes the kdb handle.
func (e *kdbC) Close() error {
	key, err := createKey("")

	if err != nil {
		return err
	}

	ret := C.kdbClose(e.handle, key.ptr)

	if ret < 0 {
		return errors.New("could not close kdb handle")
	}

	return nil
}

// Get retrieves parentKey and all Keys beneath it.
func (e *kdbC) Get(keySet KeySet, parentKey Key) (bool, error) {
	cKey, err := toCKey(parentKey)

	if err != nil {
		return false, err
	}

	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return false, err
	}

	changed := C.kdbGet(e.handle, cKeySet.ptr, cKey.ptr)

	if changed == -1 {
		return false, errFromKey(cKey)
	}

	return changed == 1, nil
}

// Set sets all Keys of a KeySet.
func (e *kdbC) Set(keySet KeySet, parentKey Key) (bool, error) {
	cKey, err := toCKey(parentKey)

	if err != nil {
		return false, err
	}

	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return false, err
	}

	changed := C.kdbSet(e.handle, cKeySet.ptr, cKey.ptr)

	if changed == -1 {
		return false, errFromKey(cKey)
	}

	return changed == 1, nil
}

// Version returns the current version of Elektra
// in the format Major.Minor.Micro
func (e *kdbC) Version() (string, error) {
	k, err := CreateKey("system/elektra/version")

	if err != nil {
		return "", err
	}

	ks, err := CreateKeySet()

	if err != nil {
		return "", err
	}

	_, err = e.Get(ks, k)

	versionKey := ks.LookupByName("system/elektra/version/constants/KDB_VERSION")
	version := versionKey.Value()

	return version, nil
}
