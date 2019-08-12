package kdb

// #cgo LDFLAGS: -lelektra
// #include <elektra/kdb.h>
// #include <stdlib.h>
//
// static Key * keyNewEmptyWrapper() {
//   return keyNew(0);
// }
//
// static Key * keyNewWrapper(char* k) {
//   return keyNew(k, KEY_END);
// }
//
// static Key * keyNewValueWrapper(char* k, char* v) {
//   return keyNew(k, KEY_VALUE, v, KEY_END);
// }
//
// static KeySet * ksNewWrapper(size_t size) {
// 	 return ksNew(size, KEY_END);
// }
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/pkg/errors"
)

// KDB is an interface to the elektra library.
type KDB interface {
	Open(key Key) error
	Close(key Key) error

	CreateKey(name string, valueAndMeta ...interface{}) (Key, error)
	CreateKeySet(keys ...Key) (KeySet, error)

	Get(keySet KeySet, parentKey Key) error
	Set(keySet KeySet, parentKey Key) error

	Version() (string, error)
}

type kdbC struct {
	handle *C.struct__KDB
}

// New returns a new KDB instance.
func New() KDB {
	return &kdbC{}
}

// Open creates a handle to the elektra library,
// this is mandatory to Get / Set Keys.
func (e *kdbC) Open(key Key) error {
	k, err := toCKey(key)

	if err != nil {
		return err
	}

	handle, ret := C.kdbOpen(k.key)

	if handle == nil {
		return fmt.Errorf("unable to open kdb: %v", ret)
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

// CreateKey creates a new key with an optional value.
func (e *kdbC) CreateKey(name string, value ...interface{}) (Key, error) {
	var key *ckey

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if name == "" {
		key = newKey(C.keyNewEmptyWrapper())
	} else if len(value) > 0 {
		switch v := value[0].(type) {
		case string:
			cValue := C.CString(v)
			key = newKey(C.keyNewValueWrapper(n, cValue))
			defer C.free(unsafe.Pointer(cValue))
		default:
			return nil, errors.New("unsupported key value type")
		}
	} else {
		key = newKey(C.keyNewWrapper(n))
	}

	if key.key == nil {
		return nil, errors.New("could not create key")
	}

	runtime.SetFinalizer(key, freeKey)

	return key, nil
}

// CreateKeySet creates a new KeySet.
func (e *kdbC) CreateKeySet(keys ...Key) (KeySet, error) {
	size := len(keys)
	ks := &ckeySet{C.ksNewWrapper(C.ulong(size))}

	if ks.keySet == nil {
		return nil, errors.New("could not create keyset")
	}

	runtime.SetFinalizer(ks, freeKeySet)

	for _, k := range keys {
		if err := ks.AppendKey(k); err != nil {
			return nil, err
		}
	}

	return ks, nil
}

// Get retrieves parentKey and all Keys beneath it.
func (e *kdbC) Get(keySet KeySet, parentKey Key) error {
	cKey, err := toCKey(parentKey)

	if err != nil {
		return err
	}

	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return err
	}

	C.kdbGet(e.handle, cKeySet.keySet, cKey.key)

	return nil
}

// Set sets all Keys of a KeySet.
func (e *kdbC) Set(keySet KeySet, parentKey Key) error {
	cKey, err := toCKey(parentKey)

	if err != nil {
		return err
	}

	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return err
	}

	C.kdbSet(e.handle, cKeySet.keySet, cKey.key)

	return nil
}

func (e *kdbC) Version() (string, error) {

	k, err := e.CreateKey("system/elektra/version")

	if err != nil {
		return "", err
	}

	ks, err := e.CreateKeySet()

	if err != nil {
		return "", err
	}

	err = e.Get(ks, k)

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
