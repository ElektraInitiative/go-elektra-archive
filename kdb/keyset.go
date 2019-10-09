package kdb

// TODO REVIEW: cleanup?
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

	Pop() Key
	Head() Key
	Tail() Key
	// TODO REVIEW API: We should remove the internal iterator and provide an external instead
	Next() Key
	Len() int

	// TODO REVIEW API: We should remove the internal iterator and provide an external instead
	Rewind()

	Cut(key Key) KeySet

	KeyNames() []string

	NeedSync() bool

	Clear()

	Lookup(key Key) Key
	LookupByName(name string) Key
}

type ckeySet struct {
	ptr *C.struct__KeySet
}

// TODO REVIEW API: Why not NewKeySet?
// CreateKeySet creates a new KeySet.
func CreateKeySet(keys ...Key) KeySet {
	size := len(keys)
	ks := &ckeySet{C.ksNewWrapper(C.ulong(size))}

	runtime.SetFinalizer(ks, freeKeySet)

	for _, k := range keys {
		if k != nil {
			ks.AppendKey(k)
		}
	}

	return ks
}

func toCKeySet(keySet KeySet) (*ckeySet, error) {
	if keySet == nil {
		return nil, errors.New("keyset is nil")
	}

	ckeySet, ok := keySet.(*ckeySet)

	if !ok {
		// TODO REVIEW: What is a ckeySet? (Error message not helpful)
		return nil, errors.New("only pointer to ckeySet struct allowed")
	}

	return ckeySet, nil
}

// TODO REVIEW: Confusing description of Append

// Append adds a KeySet and returns the new length of the KeySet
// after appending or -1 if keySet is not a pointer of type ckeySet.
func (ks *ckeySet) Append(keySet KeySet) int {
	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return -1
	}

	ret := int(C.ksAppend(ks.ptr, cKeySet.ptr))

	return ret
}

// TODO REVIEW: Confusing description of AppendKey

// AppendKey adds a Key to the KeySet  and returns the new
// length of the KeySet after appending or -1 if the key is
// not a pointer of type ckey.
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

	return &ckeySet{newKs}
}

// Head returns the first Element of the KeySet - or nil if the KeySet is empty.
func (ks *ckeySet) Head() Key {
	return newKey(C.ksHead(ks.ptr))
}

// Rewind resets the internal KeySet cursor.
func (ks *ckeySet) Rewind() {
	C.ksRewind(ks.ptr)
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
	return newKey(C.ksTail(ks.ptr))
}

// Pop removes and returns the last Element that was added to the KeySet.
func (ks *ckeySet) Pop() Key {
	return newKey(C.ksPop(ks.ptr))
}

// Remove removes a key from the KeySet and returns it if found.
func (ks *ckeySet) Remove(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	removed := C.ksLookup(ks.ptr, ckey.ptr, C.KDB_O_POP)

	return newKey(removed)
}

// Clear removes all Keys from the KeySet.
func (ks *ckeySet) Clear() {
	C.ksClear(ks.ptr)
}

// Next moves the Cursor to the next Key and returns it.
func (ks *ckeySet) Next() Key {
	key := newKey(C.ksNext(ks.ptr))

	if key == nil {
		return nil
	}

	return key
}

// Lookup searches the KeySet for a certain Key.
func (ks *ckeySet) Lookup(key Key) Key {
	ckey, err := toCKey(key)

	if err != nil {
		return nil
	}

	if foundKey := newKey(C.ksLookup(ks.ptr, ckey.ptr, 0)); foundKey != nil {
		return foundKey
	}

	return nil
}

// LookupByName searches the KeySet for a Key by name.
func (ks *ckeySet) LookupByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	if key := newKey(C.ksLookupByName(ks.ptr, n, 0)); key != nil {
		return key
	}

	return nil
}

// KeyNames returns a slice of the name of every Key in the KeySet.
func (ks *ckeySet) KeyNames() []string {
	keys := []string{}

	ks.Rewind()

	for key := ks.Next(); key != nil; key = ks.Next() {
		keys = append(keys, key.Name())
	}

	return keys
}

// Len returns the length of the KeySet.
func (ks *ckeySet) Len() int {
	return int(C.ksGetSize(ks.ptr))
}

func freeKeySet(k *ckeySet) {
	C.ksDel(k.ptr)
}
