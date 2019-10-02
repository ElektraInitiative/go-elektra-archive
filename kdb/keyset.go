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
	Copy(keySet KeySet) error
	Append(keySet KeySet) error
	AppendKey(key Key)
	Remove(key Key) error

	Pop() Key
	Head() Key
	Tail() Key
	Next() Key
	Len() int
	Rewind()

	Cut(key Key) KeySet

	KeyNames() []string

	NeedSync() bool

	Clear() error

	Lookup(key Key) Key
	LookupByName(name string) Key
}

type ckeySet struct {
	ptr *C.struct__KeySet
}

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
		return nil, errors.New("only pointer to ckeySet struct allowed")
	}

	return ckeySet, nil
}

// Append adds a Key to a KeySet.
func (ks *ckeySet) Append(key KeySet) error {
	cKeySet, err := toCKeySet(key)

	if err != nil {
		return err
	}

	ret := C.ksAppend(ks.ptr, cKeySet.ptr)

	if ret < 0 {
		return errors.New("could not append keySet to keyset")
	}

	return nil
}

// NeedSync returns true if KDB.Set() has to be called.
func (ks *ckeySet) NeedSync() bool {
	ret := C.ksNeedSync(ks.ptr)

	return ret == 1
}

func (ks *ckeySet) Cut(key Key) KeySet {
	k, err := toCKey(key)

	if err != nil {
		return nil
	}

	newKs := C.ksCut(ks.ptr, k.ptr)

	return &ckeySet{newKs}
}

// Head returns the first Element of the KeySet - or nil if empty.
func (ks *ckeySet) Head() Key {
	return newKey(C.ksHead(ks.ptr))
}

func (ks *ckeySet) Rewind() {
	C.ksRewind(ks.ptr)
}

// Copy copies the entire KeySet to a new one.
func (ks *ckeySet) Copy(keySet KeySet) error {
	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return err
	}

	C.ksCopy(cKeySet.ptr, ks.ptr)

	return nil
}

// Tail returns the last Element of the KeySet - or nil if empty.
func (ks *ckeySet) Tail() Key {
	return newKey(C.ksTail(ks.ptr))
}

// Pop removes and returns the last Element that was added to the KeySet.
func (ks *ckeySet) Pop() Key {
	return newKey(C.ksPop(ks.ptr))
}

func (ks *ckeySet) Remove(key Key) error {
	ckey, err := toCKey(key)

	if err != nil {
		return err
	}

	removed := C.ksLookup(ks.ptr, ckey.ptr, C.KDB_O_POP)

	if removed == nil {
		return errors.New("not found")
	}

	return nil
}

// AppendKey adds a Key to the KeySet.
func (ks *ckeySet) AppendKey(key Key) {
	ckey, err := toCKey(key)

	if err != nil {
		return
	}

	C.ksAppendKey(ks.ptr, ckey.ptr)
}

// Clear removes all Keys from the KeySet.
func (ks *ckeySet) Clear() error {
	ret := C.ksClear(ks.ptr)

	if ret != 0 {
		return errors.New("unable to clear keyset")
	}

	return nil
}

// Next moves the Cursor to the next Key.
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
