package kdb_test

import (
	"testing"

	elektra "github.com/ElektraInitiative/go-elektra/kdb"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestCreateKeySet(t *testing.T) {
	kdb := elektra.New()

	err := kdb.Open()
	defer kdb.Close()

	Check(t, err, "could not open KDB")

	k, err := elektra.CreateKey("user/go-elektra/test/createKeySet", "Hello World")

	Check(t, err, "could not create Key")

	ks, err := elektra.CreateKeySet(k)

	Check(t, err, "could not create KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")
}

func TestAddAndRemoveFromKeySet(t *testing.T) {
	kdb := elektra.New()

	err := kdb.Open()
	defer kdb.Close()

	Check(t, err, "could not open KDB")

	ks, err := elektra.CreateKeySet()

	Check(t, err, "could not create KeySet")

	k, err := elektra.CreateKey("user/go-elektra/test/addAndRemoveFromKeySet/1", "Hello World")

	Check(t, err, "could not create Key")

	err = ks.AppendKey(k)

	Check(t, err, "could not append to KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	k2, err := elektra.CreateKey("user/go-elektra/test/addAndRemoveFromKeySet/2", "Hello World")
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
	namespace := "user/go-elektra/test/removeKey"

	parentKey, err := elektra.CreateKey(namespace)
	Check(t, err, "could not create parent Key")

	err = kdb.Open()
	Check(t, err, "could not open KDB")
	defer kdb.Close()

	k, err := elektra.CreateKey(namespace+"/hello_world", "Hello World")
	Check(t, err, "could not create Key")

	k2, err := elektra.CreateKey(namespace+"/hello_world_2", "Hello World 2")
	Check(t, err, "could not create Key")

	ks, err := elektra.CreateKeySet()
	Check(t, err, "could not create KeySet")

	changed, err := kdb.Get(ks, parentKey)
	Assert(t, changed, "kdb.Get() has not retrieved any keys")
	Check(t, err, "could not Get KeySet")

	err = ks.AppendKey(k)
	Check(t, err, "could not append Key to KeySet")

	err = ks.AppendKey(k2)
	Check(t, err, "could not append Key to KeySet")

	changed, err = kdb.Set(ks, parentKey)
	Assert(t, changed, "kdb.Set() has not updated any keys")
	Check(t, err, "could not Set KeySet")

	_, err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey := ks.LookupByName("/go-elektra/test/removeKey/hello_world")
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey, _ = ks.Lookup(k2)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k2.Name())

	err = ks.Remove(k2)
	Check(t, err, "could not delete Key")

	changed, err = kdb.Set(ks, parentKey)
	Assert(t, changed, "kdb.Set() has not updated any keys")
	Check(t, err, "could not set KeySet")

	_, err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey, _ = ks.Lookup(k)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey, _ = ks.Lookup(k2)
	Assertf(t, foundKey == nil, "KeySet contains key %s", k2.Name())
}

func TestClearKeySet(t *testing.T) {
	k, err := elektra.CreateKey("/go-elektra/test/clearKeySet", "Hello World")
	ks, err := elektra.CreateKeySet(k)

	Check(t, err, "could not create KeySet")

	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	err = ks.Clear()

	Check(t, err, "KeySet.Clear() failed")

	Assertf(t, ks.Len() == 0, "after KeySet.Clear() KeySet.Len() should be 0 but is %d", ks.Len())
}

func TestLookupByName(t *testing.T) {
	keyName := "user/go-elektra/test/lookupByName"

	k, err := elektra.CreateKey(keyName, "Hello World")
	ks, err := elektra.CreateKeySet(k)

	Check(t, err, "could not create KeySet")

	foundKey := ks.LookupByName(keyName)

	Assert(t, foundKey != nil, "KeySet.LookupByName() did not find the correct Key")
	Assertf(t, foundKey.Name() == keyName,
		"the name of Key found by LookupByName() should be %q but is %q", k.Name(), foundKey.Name())
}
