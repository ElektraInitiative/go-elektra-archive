package kdb_test

import (
	"testing"

	elektra "github.com/ElektraInitiative/go-elektra/kdb"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestOpen(t *testing.T) {
	kdb := elektra.New()

	err := kdb.Open()
	defer kdb.Close()

	Checkf(t, err, "kdb.Open() failed: %v", err)
}

func TestSet(t *testing.T) {
	kdb := elektra.New()

	err := kdb.Open()
	defer kdb.Close()

	Checkf(t, err, "kdb.Open() failed: %v", err)

	ks := elektra.CreateKeySet()
	key, _ := elektra.CreateKey("user/tests/go-elektra/set")
	_, _ = kdb.Get(ks, key)

	ks.AppendKey(key)

	_, err = kdb.Set(ks, key)
	Checkf(t, err, "kdb Set failed %v", err)
}

func TestConflict(t *testing.T) {
	kdb1 := elektra.New()
	kdb2 := elektra.New()

	ks1 := elektra.CreateKeySet()
	ks2 := elektra.CreateKeySet()

	rootKey1, _ := elektra.CreateKey("user/tests/go-elektra/conflict")
	rootKey2, _ := elektra.CreateKey("user/tests/go-elektra/conflict")
	firstKey, _ := elektra.CreateKey("user/tests/go-elektra/conflict/first")
	secondKey, _ := elektra.CreateKey("user/tests/go-elektra/conflict/second")
	conflictKey, _ := elektra.CreateKey("user/tests/go-elektra/conflict/second")

	_ = kdb1.Open()
	defer kdb1.Close()

	_, _ = kdb1.Get(ks1, rootKey1)
	ks1.AppendKey(firstKey)
	_, _ = kdb1.Set(ks1, rootKey1)

	_ = kdb2.Open()
	defer kdb2.Close()
	_, _ = kdb2.Get(ks2, rootKey2)

	ks1.AppendKey(secondKey)
	_, _ = kdb1.Set(ks1, rootKey1)

	ks2.AppendKey(conflictKey)
	_, err := kdb2.Set(ks2, rootKey2)

	Assertf(t, err == elektra.ErrConflictingState, "expected conflict err: %v", err)
}

func TestGet(t *testing.T) {
	t.Skip()

	kdb := elektra.New()

	key, _ := elektra.CreateKey("user/tests/go-elektra/get")
	err := kdb.Open()
	defer kdb.Close()

	Checkf(t, err, "kdb.Open() failed: %v", err)

	ks := elektra.CreateKeySet()

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

	err := kdb.Open()
	defer kdb.Close()

	version, err := kdb.Version()

	Checkf(t, err, "kdb.Version() failed: %v", err)
	Assert(t, version != "", "kdb.Version() is empty")
}
