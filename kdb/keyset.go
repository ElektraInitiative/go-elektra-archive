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

	ForEach(iterator Iterator)
	ToSlice() []Key
	KeyNames() []string

	NeedSync() bool

	Clear()

	Lookup(key Key) Key
	LookupByName(name string) Key
}

type CKeySet struct {
	ptr  *C.struct__KeySet
	keys map[*C.struct__Key]*CKey
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

func wrapKeySet(ks *C.struct__KeySet) *CKeySet {
	if ks == nil {
		return nil
	}

	keySet := &CKeySet{
		ptr:  ks,
		keys: make(map[*C.struct__Key]*CKey),
	}

	keySet.forEach(func(key Key) {
		keySet.rememberKey(key.(*CKey))
	})

	runtime.SetFinalizer(keySet, freeKeySet)

	return keySet
}

// freeKeySet frees the keySet's memory when it
// goes out of scope.
func freeKeySet(k *CKeySet) {
	if k.ptr != nil {
		C.ksDel(k.ptr)
	}
}

func toCKeySet(keySet KeySet) (*CKeySet, error) {
	if keySet == nil {
		return nil, errors.New("keyset is nil")
	}

	ckeySet, ok := keySet.(*CKeySet)

	if !ok {
		return nil, errors.New("only instances of KeySet that were created by elektra/kdb may be passed to this function")
	}

	return ckeySet, nil
}

// Append appends all Keys from `other` to this KeySet and returns the
// new length of this KeySet or -1 if `other` is not a KeySet which was
// created by elektra/kdb.
func (ks *CKeySet) Append(other KeySet) int {
	ckeySet, err := toCKeySet(other)

	if err != nil {
		return -1
	}

	ret := int(C.ksAppend(ks.ptr, ckeySet.ptr))

	ckeySet.forEach(func(key Key) {
		ks.rememberKey(key.(*CKey))
	})

	return ret
}

// AppendKey appends a Key to this KeySet and returns the new
// length of this KeySet or -1 if the key is
// not a Key created by elektra/kdb.
func (ks *CKeySet) AppendKey(key Key) int {
	ckey, err := toCKey(key)

	if err != nil {
		return -1
	}

	ks.rememberKey(ckey)

	size := int(C.ksAppendKey(ks.ptr, ckey.ptr))

	return size
}

// NeedSync returns true if KDB.Set() has to be called.
func (ks *CKeySet) NeedSync() bool {
	ret := C.ksNeedSync(ks.ptr)

	return ret == 1
}

// Cut cuts out a new KeySet at the cutpoint key and returns it.
func (ks *CKeySet) Cut(key Key) KeySet {
	k, err := toCKey(key)

	if err != nil {
		return nil
	}

	newKs := wrapKeySet(C.ksCut(ks.ptr, k.ptr))

	newKs.forEach(func(key Key) {
		ks.forgetKey(k.ptr)
	})

	return newKs
}

// ToSlice returns a slice containing all Keys.
func (ks *CKeySet) ToSlice() []Key {
	var keys []Key

	ks.forEach(func(k Key) {
		keys = append(keys, k)
	})

	return keys
}

// Iterator is a function that loops over Keys.
type Iterator func(k Key)

// toKey returns a cached Key that wraps the *C.struct__Key -
// or creates a new wrapped *CKey.
func (ks *CKeySet) toKey(k *C.struct__Key) *CKey {
	if k == nil {
		return nil
	}

	if key := ks.keys[k]; key != nil {
		return key
	} else {
		return wrapKey(k)
	}
}

// rememberKey remembers the relationship between instances of *CKey
// and *C.struct__Key. This is important because we don't want multiple
// instances of *CKey pointing to the same *C.struct__Key since this
// causes troubles with Garbage Collection, which runs in parallel and
// freeing of keys is not threadsafe.
func (ks *CKeySet) rememberKey(key *CKey) {
	ks.keys[key.ptr] = key
}

// forgetKey forgets about the reference *CKey <-> *C.struct__Key. Calls this
// when a key gets removed from the underlying *ckeyset.
func (ks *CKeySet) forgetKey(k *C.struct__Key) *CKey {
	if k == nil {
		return nil
	}

	key := ks.keys[k]

	delete(ks.keys, k)

	return key
}

// forEach provides an easy way of looping of the keyset by passing
// an iterator function.
func (ks *CKeySet) forEach(iterator Iterator) {
	cursor := C.cursor_t(0)

	next := func() Key {
		key := ks.toKey(C.ksAtCursor(ks.ptr, cursor))
		cursor++

		if key == nil {
			return nil
		}

		return key
	}

	for key := next(); key != nil; key = next() {
		iterator(key)
	}
}

// ForEach accepts an `Iterator` that loops over every Key in the KeySet.
func (ks *CKeySet) ForEach(iterator Iterator) {
	ks.forEach(iterator)
}

// KeyNames returns a slice of the name of every Key in the KeySet.
func (ks *CKeySet) KeyNames() []string {
	var keys []string

	ks.forEach(func(k Key) {
		keys = append(keys, k.Name())
	})

	return keys
}

// Head returns the first Element of the KeySet - or nil if the KeySet is empty.
func (ks *CKeySet) Head() Key {
	return ks.toKey(C.ksHead(ks.ptr))
}

// Copy copies the entire KeySet to the passed KeySet.
func (ks *CKeySet) Copy(keySet KeySet) {
	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return
	}

	C.ksCopy(cKeySet.ptr, ks.ptr)

	return
}

// Tail returns the last Element of the KeySet - or nil if empty.
func (ks *CKeySet) Tail() Key {
	return ks.toKey(C.ksTail(ks.ptr))
}

// Pop removes and returns the last Element that was added to the KeySet.
func (ks *CKeySet) Pop() Key {
	key := C.ksPop(ks.ptr)

	return ks.forgetKey(key)
}

// Remove removes a key from the KeySet and returns it if found.
func (ks *CKeySet) Remove(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	removed := C.ksLookup(ks.ptr, ckey.ptr, C.KDB_O_POP)

	return ks.forgetKey(removed)
}

// RemoveByName removes a key by its name from the KeySet and returns it if found.
func (ks *CKeySet) RemoveByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	key := C.ksLookupByName(ks.ptr, n, C.KDB_O_POP)

	return ks.forgetKey(key)
}

// Clear removes all Keys from the KeySet.
func (ks *CKeySet) Clear() {
	root, _ := newKey("/")

	ks.forEach(func(k Key) {
		ks.forgetKey(k.(*CKey).ptr)
	})

	// don't use `ksClear` because it is internal
	// and renders the KeySet unusable
	newKs := C.ksCut(ks.ptr, root.ptr)

	// we don't need this keyset
	C.ksDel(newKs)
}

// Lookup searches the KeySet for a certain Key.
func (ks *CKeySet) Lookup(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	if foundKey := ks.toKey(C.ksLookup(ks.ptr, ckey.ptr, 0)); foundKey != nil {
		return foundKey
	}

	return nil
}

// LookupByName searches the KeySet for a Key by name.
func (ks *CKeySet) LookupByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if key := ks.toKey(C.ksLookupByName(ks.ptr, n, 0)); key != nil {
		return key
	}

	return nil
}

// Len returns the length of the KeySet.
func (ks *CKeySet) Len() int {
	return int(C.ksGetSize(ks.ptr))
}

/*****
	The following functions are for benchmarks only
	and should not be exported
*****/

func (ks *CKeySet) forEachInternal(iterator Iterator) {
	cursor := C.ksGetCursor(ks.ptr)
	defer C.ksSetCursor(ks.ptr, cursor)

	next := func() Key {
		key := ks.toKey(C.ksNext(ks.ptr))

		if key == nil {
			return nil
		}

		return key
	}

	C.ksRewind(ks.ptr)

	for key := next(); key != nil; key = next() {
		iterator(key)
	}
}
