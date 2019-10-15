package kdb

// #include <kdb.h>
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
import "C"

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

// Key is the wrapper around the Elektra Key.
type Key interface {
	Name() string
	Namespace() string
	BaseName() string

	String() string
	Bytes() []byte

	Meta(name string) string
	MetaMap() map[string]string
	RemoveMeta(name string) error
	MetaSlice() []Key
	NextMeta() Key

	IsBelow(key Key) bool
	IsBelowOrSame(key Key) bool
	IsDirectlyBelow(key Key) bool
	Compare(key Key) int

	Duplicate() Key

	SetMeta(name, value string) error
	SetName(name string) error
	SetString(value string) error
	SetBytes(value []byte) error
}

type ckey struct {
	ptr *C.struct__Key
}

func errFromKey(k *ckey) error {
	description := k.Meta("error/description")
	number := k.Meta("error/number")

	if err, ok := errCodeMap[number]; ok {
		return err
	}

	return fmt.Errorf("%s (%s)", description, number)
}

// NewKey creates a new `Key` with an optional value.
func NewKey(name string, value ...interface{}) (Key, error) {
	return newKey(name, value...)
}

// newKey should be used internally because the C pointer
// can be used directly without having to cast from `Key` first.
func newKey(name string, value ...interface{}) (*ckey, error) {
	var key *ckey

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if name == "" {
		key = wrapKey(C.keyNewEmptyWrapper())
	} else if len(value) > 0 {
		switch v := value[0].(type) {
		case string:
			cValue := C.CString(v)
			key = wrapKey(C.keyNewValueWrapper(n, cValue))
			defer C.free(unsafe.Pointer(cValue))
		default:
			return nil, errors.New("unsupported key value type")
		}
	} else {
		key = wrapKey(C.keyNewWrapper(n))
	}

	if key == nil {
		return nil, errors.New("could not create key")
	}

	return key, nil
}

func wrapKey(k *C.struct__Key) *ckey {
	if k == nil {
		return nil
	}

	C.keyIncRef(k)

	key := &ckey{ptr: k}

	runtime.SetFinalizer(key, freeKey)

	return key
}

// freeKey frees the resources of the Key.
func freeKey(k *ckey) {
	if k.ptr == nil {
		return
	}

	refs := C.keyDecRef(k.ptr)

	if refs == 0 {
		C.keyDel(k.ptr)
	}
}

func toCKey(key Key) (*ckey, error) {
	if key == nil {
		return nil, errors.New("key is nil")
	}

	ckey, ok := key.(*ckey)

	if !ok {
		return nil, errors.New("only pointer to ckey struct allowed")
	}

	return ckey, nil
}

// BaseName returns the basename of the Key.
// Some examples:
// - BaseName of system/some/keyname is keyname
// - BaseName of user/tmp/some key is "some key"
func (k *ckey) BaseName() string {
	name := C.keyBaseName(k.ptr)

	return C.GoString(name)
}

// Name returns the name of the Key.
func (k *ckey) Name() string {
	name := C.keyName(k.ptr)

	return C.GoString(name)
}

// SetBytes sets the value of a key to a byte slice.
func (k *ckey) SetBytes(value []byte) error {
	v := C.CBytes(value)
	defer C.free(unsafe.Pointer(v))

	size := C.ulong(len(value))

	ret := C.keySetBinary(k.ptr, unsafe.Pointer(v), size)

	fmt.Print(ret)

	return nil
}

// SetString sets the string of a key.
func (k *ckey) SetString(value string) error {
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	_ = C.keySetString(k.ptr, v)

	return nil
}

// SetBoolean sets the string of a key to a boolean
// where true is represented as "1" and false as "0".
func (k *ckey) SetBoolean(value bool) error {
	strValue := "0"

	if value {
		strValue = "1"
	}

	return k.SetString(strValue)
}

// SetName sets the name of the Key.
func (k *ckey) SetName(name string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if ret := C.keySetName(k.ptr, n); ret < 0 {
		return errors.New("could not set key name")
	}

	return nil
}

// Bytes returns the value of the Key as a byte slice.
func (k *ckey) Bytes() []byte {
	size := (C.ulong)(C.keyGetValueSize(k.ptr))

	buffer := unsafe.Pointer((*C.char)(C.malloc(size)))
	defer C.free(buffer)

	ret := C.keyGetBinary(k.ptr, buffer, C.ulong(size))

	if ret <= 0 {
		return []byte{}
	}

	bytes := C.GoBytes(buffer, C.int(size))

	return bytes
}

// String returns the string value of the Key.
func (k *ckey) String() string {
	str := C.keyString(k.ptr)

	return C.GoString(str)
}

// SetMeta sets the meta value of a Key.
func (k *ckey) SetMeta(name, value string) error {
	cName, cValue := C.CString(name), C.CString(value)

	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cValue))

	ret := C.keySetMeta(k.ptr, cName, cValue)

	if ret < 0 {
		return errors.New("could not set meta")
	}

	return nil
}

// DeleteMeta deletes a meta Key.
func (k *ckey) RemoveMeta(name string) error {
	cName := C.CString(name)

	defer C.free(unsafe.Pointer(cName))

	ret := C.keySetMeta(k.ptr, cName, nil)

	if ret < 0 {
		return errors.New("could not delete meta")
	}

	return nil
}

// Meta retrieves the Meta value of a Key.
func (k *ckey) Meta(name string) string {
	cName := C.CString(name)

	defer C.free(unsafe.Pointer(cName))

	metaKey := wrapKey(C.keyGetMeta(k.ptr, cName))

	if metaKey == nil {
		return ""
	}

	return metaKey.String()
}

// NextMeta returns the next meta Key.
func (k *ckey) NextMeta() Key {
	key := wrapKey(C.keyNextMeta(k.ptr))

	if key == nil {
		return nil
	}

	return key
}

// MetaSlice builds a slice of all meta Keys.
func (k *ckey) MetaSlice() []Key {
	dup := k.Duplicate().(*ckey)
	C.keyRewindMeta(dup.ptr)

	var metaKeys []Key

	for key := dup.NextMeta(); key != nil; key = dup.NextMeta() {
		metaKeys = append(metaKeys, key)
	}

	return metaKeys
}

// MetaMap builds a Key/Value map of all meta Keys.
func (k *ckey) MetaMap() map[string]string {
	dup := k.Duplicate().(*ckey)
	C.keyRewindMeta(dup.ptr)

	m := make(map[string]string)

	for key := dup.NextMeta(); key != nil; key = dup.NextMeta() {
		m[key.Name()] = key.String()
	}

	return m
}

// Duplicate duplicates a Key.
func (k *ckey) Duplicate() Key {
	return wrapKey(C.keyDup(k.ptr))
}

// IsBelow checks if this key is below the `other` key.
func (k *ckey) IsBelow(other Key) bool {
	otherKey, err := toCKey(other)

	if err != nil {
		return false
	}

	ret := C.keyIsBelow(otherKey.ptr, k.ptr)

	return ret != 0
}

// IsBelowOrSame checks if this key is below or the same as the `other` key.
func (k *ckey) IsBelowOrSame(other Key) bool {
	otherKey, err := toCKey(other)

	if err != nil {
		return false
	}

	ret := C.keyIsBelowOrSame(otherKey.ptr, k.ptr)

	return ret != 0
}

// IsDirectlyBelow checks if this key is directly below the `other` Key.
func (k *ckey) IsDirectlyBelow(other Key) bool {
	otherKey, err := toCKey(other)

	if err != nil {
		return false
	}

	ret := C.keyIsDirectlyBelow(otherKey.ptr, k.ptr)

	return ret != 0
}

// Compare the name of two keys. It returns 0 if the keys are equal,
// < 0 if this key is less than `other` Key and
// > 0 if this key is greater than `other` Key.
// This function defines the sorting order of a KeySet.
func (k *ckey) Compare(other Key) int {
	otherKey, _ := toCKey(other)

	return int(C.keyCmp(k.ptr, otherKey.ptr))
}

// Namespace returns the namespace of a Key.
func (k *ckey) Namespace() string {
	name := k.Name()
	index := strings.Index(name, "/")

	if index < 0 {
		return ""
	}

	return name[:index]
}

func nameWithoutNamespace(key Key) string {
	name := key.Name()
	index := strings.Index(name, "/")

	if index < 0 {
		return "/"
	}

	return name[index:]
}

// CommonKeyName returns the common path of two Keys.
func CommonKeyName(key1, key2 Key) string {
	key1Name := key1.Name()
	key2Name := key2.Name()

	if key1.IsBelowOrSame(key2) {
		return key2Name
	}
	if key2.IsBelowOrSame(key1) {
		return key1Name
	}

	if key1.Namespace() != key2.Namespace() {
		key1Name = nameWithoutNamespace(key1)
		key2Name = nameWithoutNamespace(key2)
	}

	index := 0
	k1Parts, k2Parts := strings.Split(key1Name, "/"), strings.Split(key2Name, "/")

	for ; index < len(k1Parts) && index < len(k2Parts) && k1Parts[index] == k2Parts[index]; index++ {
	}

	return strings.Join(k1Parts[:index], "/")
}
