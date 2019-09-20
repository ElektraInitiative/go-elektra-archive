package kdb_test

import (
	"testing"

	elektra "github.com/ElektraInitiative/go-elektra/kdb"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestOpen(t *testing.T) {
	kdb := elektra.New()

	k, _ := elektra.CreateKey("Test", "Value")

	err := kdb.Open(k)

	Checkf(t, err, "kdb.Open() failed: %v", err)
}

func TestSet(t *testing.T) {
	kdb := elektra.New()

	parent, _ := elektra.CreateKey("/error")
	err := kdb.Open(parent)

	Checkf(t, err, "kdb.Open() failed: %v", err)

	ks, _ := elektra.CreateKeySet()
	key, _ := elektra.CreateKey("user/go-binding-low/test")
	_, _ = kdb.Get(ks, key)

	err = ks.AppendKey(key)
	Checkf(t, err, "KeySet.AppendKey() failed: %v", err)

	_, err = kdb.Set(ks, parent)
	Checkf(t, err, "kdb Set failed %v", err)
}

func TestConflict(t *testing.T) {
	kdb1 := elektra.New()
	kdb2 := elektra.New()

	ks1, _ := elektra.CreateKeySet()
	ks2, _ := elektra.CreateKeySet()

	rootKey1, _ := elektra.CreateKey("user/go-elektra/test/conflict")
	rootKey2, _ := elektra.CreateKey("user/go-elektra/test/conflict")
	firstKey, _ := elektra.CreateKey("user/go-elektra/test/conflict/first")
	secondKey, _ := elektra.CreateKey("user/go-elektra/test/conflict/second")
	conflictKey, _ := elektra.CreateKey("user/go-elektra/test/conflict/second")

	_ = kdb1.Open(rootKey1)
	_, _ = kdb1.Get(ks1, rootKey1)
	_ = ks1.AppendKey(firstKey)
	_, _ = kdb1.Set(ks1, rootKey1)

	_ = kdb2.Open(rootKey2)
	_, _ = kdb2.Get(ks2, rootKey2)

	_ = ks1.AppendKey(secondKey)
	_, _ = kdb1.Set(ks1, rootKey1)

	_ = ks2.AppendKey(conflictKey)
	_, err := kdb2.Set(ks2, rootKey2)

	Assertf(t, err == elektra.ErrConflictingState, "expected conflict err: %v", err)
}

func TestGet(t *testing.T) {
	t.Skip()

	kdb := elektra.New()

	key, _ := elektra.CreateKey("/bla")
	err := kdb.Open(key)

	Checkf(t, err, "kdb.Open() failed: %v", err)

	ks, _ := elektra.CreateKeySet()

	changed, err := kdb.Get(ks, key)

	Assert(t, changed, "kdb.Get() has not retrieved any keys")
	Checkf(t, err, "kdb.Get() failed: %v", err)

	t.Log(ks.Len())

	for next := ks.Next(); next != nil; next = ks.Next() {
		t.Log(next)
	}

	t.Log(ks.LookupByName("/bla"))
}

func TestVersion(t *testing.T) {
	kdb := elektra.New()

	key, _ := elektra.CreateKey("/bla")
	err := kdb.Open(key)

	version, err := kdb.Version()

	Checkf(t, err, "kdb.Version() failed: %v", err)
	Assert(t, version != "", "kdb.Version() is empty")
}
