package kdb

// #include <kdb.h>
// #include <stdlib.h>
//
// static KeySet * ksNewWrapper(size_t size) {
// 	 return ksNew(size, KEY_END);
// }
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/pkg/errors"
)

// KeySet represents a collection of Keys.
type KeySet interface {
	Copy(keySet KeySet)
	Append(keySet KeySet) int
	AppendKey(key Key) int
	Remove(key Key) Key
	RemoveByName(name string) Key

	Pop() Key
	Head() Key
	Tail() Key
	Len() int

	Cut(key Key) KeySet

	Each(iterator Iterator)
	Slice() []Key
	KeyNames() []string

	NeedSync() bool

	Clear()

	Lookup(key Key) Key
	LookupByName(name string) Key
}

type ckeySet struct {
	ptr *C.struct__KeySet
}

// NewKeySet creates a new KeySet.
func NewKeySet(keys ...Key) KeySet {
	size := len(keys)
	ks := wrapKeySet(C.ksNewWrapper(C.ulong(size)))

	for _, k := range keys {
		ks.AppendKey(k)
	}

	return ks
}

func wrapKeySet(ks *C.struct__KeySet) *ckeySet {
	if ks == nil {
		return nil
	}

	keySet := &ckeySet{ks}

	runtime.SetFinalizer(keySet, freeKeySet)

	return keySet
}

func freeKeySet(k *ckeySet) {
	if k.ptr != nil {
		C.ksDel(k.ptr)
	}
}

func toCKeySet(keySet KeySet) (*ckeySet, error) {
	if keySet == nil {
		return nil, errors.New("keyset is nil")
	}

	ckeySet, ok := keySet.(*ckeySet)

	if !ok {
		return nil, errors.New("only instances of KeySet that were created by elektra/kdb may be passed to this function")
	}

	return ckeySet, nil
}

// Append appends all Keys from `other` to this KeySet and returns the
// new length of this KeySet or -1 if `other` is not a KeySet which was
// created by elektra/kdb.
func (ks *ckeySet) Append(other KeySet) int {
	cKeySet, err := toCKeySet(other)

	if err != nil {
		return -1
	}

	ret := int(C.ksAppend(ks.ptr, cKeySet.ptr))

	return ret
}

// AppendKey appends a Key to this KeySet and returns the new
// length of this KeySet or -1 if the key is
// not a Key created by elektra/kdb.
func (ks *ckeySet) AppendKey(key Key) int {
	ckey, err := toCKey(key)

	if err != nil {
		return -1
	}

	size := int(C.ksAppendKey(ks.ptr, ckey.ptr))

	return size
}

// NeedSync returns true if KDB.Set() has to be called.
func (ks *ckeySet) NeedSync() bool {
	ret := C.ksNeedSync(ks.ptr)

	return ret == 1
}

// Cut cuts out a new KeySet at the cutpoint key and returns it.
func (ks *ckeySet) Cut(key Key) KeySet {
	k, err := toCKey(key)

	if err != nil {
		return nil
	}

	newKs := C.ksCut(ks.ptr, k.ptr)

	return wrapKeySet(newKs)
}

// Slice returns a slice containing Keys.
func (ks *ckeySet) Slice() []Key {
	var metaKeys []Key

	ks.loop(func(k Key) bool {
		metaKeys = append(metaKeys, k)
		return true
	})

	return metaKeys
}

// Iterator is a function that loops over Keys.
// Return false if you want break out of the loop.
type Iterator func(k Key) bool

func (ks *ckeySet) loop(iterator Iterator) {
	for cursor := C.cursor_t(0); C.ksAtCursor(ks.ptr, cursor) != nil; cursor++ {
		key := wrapKey(C.ksAtCursor(ks.ptr, cursor))
		cont := iterator(key)

		if !cont {
			break
		}
	}

}

// Loop accepts an `Iterator` that loops over every Key in the KeySet.
func (ks *ckeySet) Each(iterator Iterator) {
	ks.loop(iterator)
}

// KeyNames returns a slice of the name of every Key in the KeySet.
func (ks *ckeySet) KeyNames() []string {
	var keys []string

	ks.loop(func(k Key) bool {
		keys = append(keys, k.Name())
		return true
	})

	return keys
}

// Head returns the first Element of the KeySet - or nil if the KeySet is empty.
func (ks *ckeySet) Head() Key {
	return wrapKey(C.ksHead(ks.ptr))
}

// Copy copies the entire KeySet to the passed KeySet.
func (ks *ckeySet) Copy(keySet KeySet) {
	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return
	}

	C.ksCopy(cKeySet.ptr, ks.ptr)

	return
}

// Tail returns the last Element of the KeySet - or nil if empty.
func (ks *ckeySet) Tail() Key {
	return wrapKey(C.ksTail(ks.ptr))
}

// Pop removes and returns the last Element that was added to the KeySet.
func (ks *ckeySet) Pop() Key {
	return wrapKey(C.ksPop(ks.ptr))
}

// Remove removes a key from the KeySet and returns it if found.
func (ks *ckeySet) Remove(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	removed := C.ksLookup(ks.ptr, ckey.ptr, C.KDB_O_POP)

	return wrapKey(removed)
}

// RemoveByName removes a key by its name from the KeySet and returns it if found.
func (ks *ckeySet) RemoveByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if key := wrapKey(C.ksLookupByName(ks.ptr, n, C.KDB_O_POP)); key != nil {
		return key
	}

	return nil
}

// Clear removes all Keys from the KeySet.
func (ks *ckeySet) Clear() {
	C.ksClear(ks.ptr)
}

// Lookup searches the KeySet for a certain Key.
func (ks *ckeySet) Lookup(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	if foundKey := wrapKey(C.ksLookup(ks.ptr, ckey.ptr, 0)); foundKey != nil {
		return foundKey
	}

	return nil
}

// LookupByName searches the KeySet for a Key by name.
func (ks *ckeySet) LookupByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if key := wrapKey(C.ksLookupByName(ks.ptr, n, 0)); key != nil {
		return key
	}

	return nil
}

// Len returns the length of the KeySet.
func (ks *ckeySet) Len() int {
	return int(C.ksGetSize(ks.ptr))
}

/*****
	The following functions are for benchmarks only
	and should not be exported
*****/

func (ks *ckeySet) loopInternal(iterator Iterator) {
	cursor := C.ksGetCursor(ks.ptr)
	defer C.ksSetCursor(ks.ptr, cursor)

	next := func() Key {
		key := wrapKey(C.ksNext(ks.ptr))

		if key == nil {
			return nil
		}

		return key
	}

	C.ksRewind(ks.ptr)

	for key := next(); key != nil; key = next() {
		cont := iterator(key)

		if !cont {
			break
		}
	}
}
