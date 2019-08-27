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
	AppendKey(key Key) error
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

	Lookup(key Key) (Key, error)
	LookupByName(name string) Key
}

type ckeySet struct {
	keySet *C.struct__KeySet
}

// CreateKeySet creates a new KeySet.
func CreateKeySet(keys ...Key) (KeySet, error) {
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

	ret := C.ksAppend(ks.keySet, cKeySet.keySet)

	if ret < 0 {
		return errors.New("could not append keySet to keyset")
	}

	return nil
}

// NeedSync returns true if KDB.Set() has to be called.
func (ks *ckeySet) NeedSync() bool {
	ret := C.ksNeedSync(ks.keySet)

	return ret == 1
}

func (ks *ckeySet) Cut(key Key) KeySet {
	k, err := toCKey(key)

	if err != nil {
		return nil
	}

	newKs := C.ksCut(ks.keySet, k.key)

	return &ckeySet{newKs}
}

// Head returns the first Element of the KeySet - or nil if empty.
func (ks *ckeySet) Head() Key {
	key := newKey(C.ksHead(ks.keySet))

	if key.isNil() {
		return nil
	}

	return key
}

func (ks *ckeySet) Rewind() {
	C.ksRewind(ks.keySet)
}

// Copy copies the entire KeySet to a new one.
func (ks *ckeySet) Copy(keySet KeySet) error {
	cKeySet, err := toCKeySet(keySet)

	if err != nil {
		return err
	}

	C.ksCopy(cKeySet.keySet, ks.keySet)

	return nil
}

// Tail returns the last Element of the KeySet - or nil if empty.
func (ks *ckeySet) Tail() Key {
	key := newKey(C.ksTail(ks.keySet))

	if key.isNil() {
		return nil
	}

	return key
}

// Pop removes and returns the last Element that was added to the KeySet.
func (ks *ckeySet) Pop() Key {
	key := newKey(C.ksPop(ks.keySet))

	if key.isNil() {
		return nil
	}

	return key
}

func (ks *ckeySet) Remove(key Key) error {
	ckey, err := toCKey(key)

	if err != nil {
		return err
	}

	removed := C.ksLookup(ks.keySet, ckey.key, C.KDB_O_POP)

	if removed == nil {
		return errors.New("not found")
	}

	return nil
}

// AppendKey adds a Key to the KeySet.
func (ks *ckeySet) AppendKey(key Key) error {
	ckey, err := toCKey(key)

	if err != nil {
		return err
	}

	ret := C.ksAppendKey(ks.keySet, ckey.key)

	if ret < 0 {
		return errors.New("could not append key to keyset")
	}

	return nil
}

// Clear removes all Keys from the KeySet.
func (ks *ckeySet) Clear() error {
	ret := C.ksClear(ks.keySet)

	if ret != 0 {
		return errors.New("unable to clear keyset")
	}

	return nil
}

// Next moves the Cursor to the next Key.
func (ks *ckeySet) Next() Key {
	key := newKey(C.ksNext(ks.keySet))

	if key.isNil() {
		return nil
	}

	return key
}

// Lookup searches the KeySet for a certain Key.
func (ks *ckeySet) Lookup(key Key) (Key, error) {
	ckey, err := toCKey(key)

	if err != nil {
		return nil, err
	}

	foundKey := newKey(C.ksLookup(ks.keySet, ckey.key, 0))

	if foundKey.isNil() {
		return nil, nil
	}

	return foundKey, nil
}

// LookupByName searches the KeySet for a Key by name.
func (ks *ckeySet) LookupByName(name string) Key {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	key := newKey(C.ksLookupByName(ks.keySet, n, 0))

	if key.isNil() {
		return nil
	}

	return key
}

func (ks *ckeySet) KeyNames() []string {
	keys := []string{}

	// save cursor
	cursor := C.ksGetCursor(ks.keySet)

	ks.Rewind()

	for key := ks.Next(); key != nil; key = ks.Next() {
		keys = append(keys, key.Name())
	}

	// and reset it after iterating over the keys
	C.ksSetCursor(ks.keySet, cursor)

	return keys
}

// Len returns the length of the KeySet.
func (ks *ckeySet) Len() int {
	return int(C.ksGetSize(ks.keySet))
}

func freeKeySet(k *ckeySet) {
	C.ksDel(k.keySet)
}
