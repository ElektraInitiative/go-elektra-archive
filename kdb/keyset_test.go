package kdb_test

import (
	"fmt"
	"testing"
	"time"

	elektra "go.libelektra.org/kdb"
	. "go.libelektra.org/test"
)

func setupKeySetData(count int) {
	kdb := elektra.New()

	kdb.Open()
	defer kdb.Close()

	ks := elektra.NewKeySet()
	defer ks.Close()

	root, _ := elektra.NewKey("user/tests/go/elektra/memory")

	kdb.Get(ks, root)

	for n := 0; n < count; n++ {
		k, _ := elektra.NewKey(fmt.Sprintf("user/tests/go/elektra/memory/%06d", n))
		ks.AppendKey(k)
	}

	kdb.Set(ks, root)
}

func TestKeySetMemoryWithDelayedClose(t *testing.T) {
	setupKeySetData(100000)

	keysets := [20]struct {
		ks     elektra.KeySet
		handle elektra.KDB
	}{}

	for i := range keysets {
		root, _ := elektra.NewKey("/")

		kdb := elektra.New()

		kdb.Open()

		ks := elektra.NewKeySet()

		_, err := kdb.Get(ks, root)

		if err != nil {
			t.Error(err)
		}

		fmt.Printf("Loaded %d keys\n", ks.Len())

		time.Sleep(1 * time.Second)

		if err != nil {
			t.Fatalf("could not close kdb handle")
		}

		keysets[i].handle = kdb
		keysets[i].ks = ks
	}

	for i := range keysets {
		keysets[i].ks.Close()
		keysets[i].handle.Close()

		keysets[i].ks = nil
		keysets[i].handle = nil
	}

	time.Sleep(20 * time.Second)
}

func TestKeySetMemory(t *testing.T) {
	setupKeySetData(100000)

	root, _ := elektra.NewKey("/")

	for i := 0; i < 20; i++ {
		kdb := elektra.New()

		kdb.Open()

		ks := elektra.NewKeySet()

		_, err := kdb.Get(ks, root)

		if err != nil {
			t.Error(err)
		}

		fmt.Printf("Loaded %d keys\n", ks.Len())

		time.Sleep(1 * time.Second)

		err = kdb.Close()

		if err != nil {
			t.Fatalf("could not close kdb handle")
		}
		ks.Close()
	}

	time.Sleep(20 * time.Second)
}

func TestCreateKeySet(t *testing.T) {
	k, err := elektra.NewKey("user/tests/go/elektra/createkeyset", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)
	Assert(t, ks.Len() == 1, "KeySet should have len 1")
}

func TestSlice(t *testing.T) {
	k, err := elektra.NewKey("user/tests/go/elektra/addandremovefromkeyset/1", "Hello World")
	Check(t, err, "could not create Key")
	k2, err := elektra.NewKey("user/tests/go/elektra/addandremovefromkeyset/2", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k, k2)

	slice := ks.ToSlice()
	Assert(t, len(slice) == 2, "Slice should have length 2")
	Assert(t, slice[0].Compare(k) == 0 && slice[1].Compare(k2) == 0, "Slice does not contain the correct keys")
}

func TestKeyNames(t *testing.T) {
	keyName1 := "user/tests/go/elektra/addandremovefromkeyset/1"
	k, err := elektra.NewKey(keyName1, "Hello World")
	Check(t, err, "could not create Key")

	keyName2 := "user/tests/go/elektra/addandremovefromkeyset/2"
	k2, err := elektra.NewKey(keyName2, "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k, k2)

	keyNames := ks.KeyNames()
	Assert(t, len(keyNames) == 2, "KeyNames should have length 2")
	Assert(t, keyNames[0] == keyName1 && keyNames[1] == keyName2, "")

}

func TestAddAndRemoveFromKeySet(t *testing.T) {
	ks := elektra.NewKeySet()

	k, err := elektra.NewKey("user/tests/go/elektra/addandremovefromkeyset/1", "Hello World")
	Check(t, err, "could not create Key")

	size := ks.AppendKey(k)
	Assert(t, size == 1, "KeySet should have len 1")

	k2, err := elektra.NewKey("user/tests/go/elektra/addandremovefromkeyset//2", "Hello World")
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

func TestRemove(t *testing.T) {

	k1, err := elektra.NewKey("user/tests/go/elektra/remove/1", "Hello World")
	Check(t, err, "could not create Key")
	k2, err := elektra.NewKey("user/tests/go/elektra/remove/2", "Hello World")
	Check(t, err, "could not create Key")
	k3, err := elektra.NewKey("user/tests/go/elektra/remove/3", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k1, k2, k3)

	Assert(t, ks.Len() == 3, "KeySet should have length 3")

	removed := ks.Remove(k1)
	Assert(t, removed != nil, "Remove failed")
	Assert(t, ks.Len() == 2, "KeySet should have length 2")

	removed = ks.RemoveByName("user/tests/go/elektra/remove/2")
	Assert(t, removed != nil, "RemoveByName failed")
	Assert(t, ks.Len() == 1, "KeySet should have length 2")
}

func TestClearKeySet(t *testing.T) {
	k, err := elektra.NewKey("user/tests/go/elektra/clearkeyset/1", "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)
	Check(t, err, "could not create KeySet")
	Assert(t, ks.Len() == 1, "KeySet should have len 1")

	k2, err := elektra.NewKey("user/tests/go/elektra/clearkeyset/2", "Hello World")
	Check(t, err, "could not create Key")
	ks.AppendKey(k2)

	ks.Clear()
	Check(t, err, "KeySet.Clear() failed")
	Assertf(t, ks.Len() == 0, "after KeySet.Clear() KeySet.Len() should be 0 but is %d", ks.Len())
}

func TestLookupByName(t *testing.T) {
	keyName := "user/tests/go/elektra/lookupbyname"

	k, err := elektra.NewKey(keyName, "Hello World")
	Check(t, err, "could not create Key")

	ks := elektra.NewKeySet(k)

	foundKey := ks.LookupByName(keyName)
	Assert(t, foundKey != nil, "KeySet.LookupByName() did not find the correct Key")
	Assertf(t, foundKey.Name() == keyName,
		"the name of Key found by LookupByName() should be %q but is %q", k.Name(), foundKey.Name())
}
