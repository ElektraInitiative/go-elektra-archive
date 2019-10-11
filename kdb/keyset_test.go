package kdb_test

import (
	"testing"

	elektra "github.com/ElektraInitiative/go-elektra/kdb"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestCreateKeySet(t *testing.T) {
	// TODO REVIEW: In Elektra all keys should be lower-case without any separator (except /), so you would use user/tests/go/elektra/create/key/set (or createkeyset)
	k, err := elektra.NewKey("user/tests/go/elektra/createKeySet", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)
	Assert(t, ks.Len() == 1, "KeySet should have len 1")
}

func TestAddAndRemoveFromKeySet(t *testing.T) {
	ks := elektra.NewKeySet()

	k, err := elektra.NewKey("user/tests/go/elektra/addAndRemoveFromKeySet/1", "Hello World")
	Check(t, err, "could not create Key")

	size := ks.AppendKey(k)
	Assert(t, size == 1, "KeySet should have len 1")

	k2, err := elektra.NewKey("user/tests/go/elektra/addAndRemoveFromKeySet/2", "Hello World")
	Check(t, err, "could not create Key")

	size = ks.AppendKey(k2)
	Assert(t, ks.Len() == 2, "KeySet should have len 2")

	k3 := ks.Pop()
	Assert(t, k3 != nil, "could not pop key from KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	k4 := ks.Pop()
	Assert(t, k4 != nil, "could not pop key from KeySet")
	Assert(t, ks.Len() == 0, "KeySet should have len 0")
}

func TestRemoveKey(t *testing.T) {
	// TODO REVIEW: why is kdb opened here in the KeySet tests? (should be in separated file)
	kdb := elektra.New()
	namespace := "user/tests/go/elektra/removeKey"

	parentKey, err := elektra.NewKey(namespace)
	Check(t, err, "could not create parent Key")

	err = kdb.Open()
	Check(t, err, "could not open KDB")
	defer kdb.Close()

	k, err := elektra.NewKey(namespace+"/helloworld", "Hello World")
	Check(t, err, "could not create Key")

	k2, err := elektra.NewKey(namespace+"/helloworld2", "Hello World 2")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet()
	Check(t, err, "could not create KeySet")

	changed, err := kdb.Get(ks, parentKey)
	Assert(t, changed, "kdb.Get() has not retrieved any keys")
	Check(t, err, "could not Get KeySet")

	ks.AppendKey(k)
	Check(t, err, "could not append Key to KeySet")

	ks.AppendKey(k2)
	Check(t, err, "could not append Key to KeySet")

	changed, err = kdb.Set(ks, parentKey)
	Assert(t, changed, "kdb.Set() has not updated any keys")
	Check(t, err, "could not Set KeySet")

	_, err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey := ks.LookupByName("/tests/go/elektra/removeKey/helloworld")
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey = ks.Lookup(k2)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k2.Name())

	removed := ks.Remove(k2)
	Assert(t, removed != nil, "could not delete Key")

	changed, err = kdb.Set(ks, parentKey)
	Assert(t, changed, "kdb.Set() has not updated any keys")
	Check(t, err, "could not set KeySet")

	_, err = kdb.Get(ks, parentKey)
	Check(t, err, "could not Get KeySet")

	foundKey = ks.Lookup(k)
	Assertf(t, foundKey != nil, "KeySet does not contain key %s", k.Name())

	foundKey = ks.Lookup(k2)
	Assertf(t, foundKey == nil, "KeySet contains key %s", k2.Name())

	// TODO REVIEW: How to pop via ksLookup? Or do you have Pop with argument? (which is maybe nicer API)
}

func TestClearKeySet(t *testing.T) {
	k, err := elektra.NewKey("user/tests/go/elektra/clearKeySet", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)

	Check(t, err, "could not create KeySet")

	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	ks.Clear()
	Check(t, err, "KeySet.Clear() failed")

	Assertf(t, ks.Len() == 0, "after KeySet.Clear() KeySet.Len() should be 0 but is %d", ks.Len())
}

func TestLookupByName(t *testing.T) {
	keyName := "user/tests/go/elektra/lookupByName"

	k, err := elektra.NewKey(keyName, "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)

	foundKey := ks.LookupByName(keyName)

	Assert(t, foundKey != nil, "KeySet.LookupByName() did not find the correct Key")
	Assertf(t, foundKey.Name() == keyName,
		"the name of Key found by LookupByName() should be %q but is %q", k.Name(), foundKey.Name())
}
