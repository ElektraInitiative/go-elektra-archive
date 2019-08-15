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
