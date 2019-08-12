package kdb_test

import (
	"testing"

	"github.com/ElektraInitiative/go-elektra"
)

func TestCreateKeySet(t *testing.T) {
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world", "Hello World")

	Check(t, err, "could not create Key")

	err = kdb.Open(k)

	Check(t, err, "could not open KDB")

	ks, err := kdb.CreateKeySet(k)

	Check(t, err, "could not create KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")
}

func TestAddAndRemoveFromKeySet(t *testing.T) {
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world", "Hello World")

	Check(t, err, "could not create Key")

	err = kdb.Open(k)

	Check(t, err, "could not open KDB")

	ks, err := kdb.CreateKeySet()

	Check(t, err, "could not create KeySet")

	err = ks.AppendKey(k)

	Check(t, err, "could not append to KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	k2, err := kdb.CreateKey("user/hello_world_2", "Hello World")
	Check(t, err, "could not create Key")

	err = ks.AppendKey(k2)

	Check(t, err, "could not append to KeySet")
	Assert(t, ks.Len() == 2, "KeySet should have len 2")

	k3 := ks.Pop()

	Assert(t, k3 != nil, "could not pop key from KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")
}

func TestRemoveKey(t *testing.T) {
	kdb := elektra.New()
	namespace := "user/test"

	parentKey, err := kdb.CreateKey(namespace)
	Check(t, err, "could not create parent Key")

	err = kdb.Open(parentKey)
	Check(t, err, "could not open KDB")

	k, err := kdb.CreateKey(namespace+"/hello_world", "Hello World")
	Check(t, err, "could not create Key")

	k2, err := kdb.CreateKey(namespace+"/hello_world_2", "Hello World 2")
	Check(t, err, "could not create Key")

	ks, err := kdb.CreateKeySet()
	Check(t, err, "could not create KeySet")

	err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	err = ks.AppendKey(k)
	Check(t, err, "could not append Key to KeySet")

	err = ks.AppendKey(k2)
	Check(t, err, "could not append Key to KeySet")

	err = kdb.Set(ks, parentKey)
	Check(t, err, "could not Set KeySet")

	err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey := ks.LookupByName("/test/hello_world")
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey, _ = ks.Lookup(k2)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k2.Name())

	err = ks.Remove(k2)
	Check(t, err, "could not delete Key")

	err = kdb.Set(ks, parentKey)
	Check(t, err, "could not set KeySet")

	err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey, _ = ks.Lookup(k)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey, _ = ks.Lookup(k2)
	Assertf(t, foundKey == nil, "KeySet contains key %s", k2.Name())
}

func TestClearKeySet(t *testing.T) {
	kdb := elektra.New()

	k, err := kdb.CreateKey("user/hello_world", "Hello World")
	ks, err := kdb.CreateKeySet(k)

	Check(t, err, "could not create KeySet")

	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	err = ks.Clear()

	Check(t, err, "KeySet.Clear() failed")

	Assertf(t, ks.Len() == 0, "after KeySet.Clear() KeySet.Len() should be 0 but is %d", ks.Len())
}

func TestLookupByName(t *testing.T) {
	kdb := elektra.New()

	keyName := "user/hello_world"

	k, err := kdb.CreateKey(keyName, "Hello World")
	ks, err := kdb.CreateKeySet(k)

	Check(t, err, "could not create KeySet")

	foundKey := ks.LookupByName(keyName)

	Assert(t, foundKey != nil, "KeySet.LookupByName() did not find the correct Key")
	Assertf(t, foundKey.Name() == keyName,
		"the name of Key found by LookupByName() should be %q but is %q", k.Name(), foundKey.Name())
}
