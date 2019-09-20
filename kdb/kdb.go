package kdb

// #cgo pkg-config: elektra
// #include <kdb.h>
import "C"

import (
	"github.com/pkg/errors"
)

// KDB is an interface to the Elektra library.
type KDB interface {
	Open(key Key) error
	Close(key Key) error

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
func (e *kdbC) Open(key Key) error {
	k, err := toCKey(key)

	if err != nil {
		return err
	}

	handle := C.kdbOpen(k.key)

	if handle == nil {
		return errFromKey(k)
	}

	e.handle = handle

	return nil
}

// Close closes the kdb handle.
func (e *kdbC) Close(key Key) error {
	ckey, err := toCKey(key)

	if err != nil {
		return err
	}

	ret := C.kdbClose(e.handle, ckey.key)

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

	changed := C.kdbGet(e.handle, cKeySet.keySet, cKey.key)

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

	changed := C.kdbSet(e.handle, cKeySet.keySet, cKey.key)

	if changed == -1 {
		return false, errFromKey(cKey)
	}

	return changed == 1, nil
}

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

const (
	KeyName          uint = C.KEY_NAME
	KeyValue         uint = C.KEY_VALUE
	KeyFlags         uint = C.KEY_FLAGS
	KeyOwner         uint = C.KEY_OWNER
	KeyComment       uint = C.KEY_COMMENT
	KeyBinary        uint = C.KEY_BINARY
	KeyUid           uint = C.KEY_UID
	KeyGid           uint = C.KEY_GID
	KeyMode          uint = C.KEY_MODE
	KeyAtime         uint = C.KEY_ATIME
	KeyMtime         uint = C.KEY_MTIME
	KeyCtime         uint = C.KEY_CTIME
	KeySize          uint = C.KEY_SIZE
	KeyDir           uint = C.KEY_DIR
	KeyMeta          uint = C.KEY_META
	KeyNull          uint = C.KEY_NULL
	KeyCascadingName uint = C.KEY_CASCADING_NAME
	KeyMetaName      uint = C.KEY_META_NAME
	KeyEnd           uint = C.KEY_END
)
