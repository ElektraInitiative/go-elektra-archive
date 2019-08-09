package elektra_test

import (
	"testing"

	"github.com/ElektraInitiative/go-elektra"
)

func TestOpen(t *testing.T) {
	kdb := elektra.New()

	k, _ := kdb.CreateKey("Test", "Value")

	err := kdb.Open(k)

	Checkf(t, err, "kdb.Open() failed: %v", err)
}

func TestSet(t *testing.T) {
	kdb := elektra.New()

	parent, _ := kdb.CreateKey("/bla")
	err := kdb.Open(parent)

	Checkf(t, err, "kdb.Open() failed: %v", err)

	key, _ := kdb.CreateKey("/bla/bla2")
	ks, _ := kdb.CreateKeySet(key)

	Assert(t, ks.Len() == 1, "KeySet.Len() should be 1")

	err = kdb.Set(ks, parent)

	Checkf(t, err, "kdb Set failed %v", err)
}

func TestGet(t *testing.T) {
	t.Skip()

	kdb := elektra.New()

	key, _ := kdb.CreateKey("/bla")
	err := kdb.Open(key)

	Checkf(t, err, "kdb.Open() failed: %v", err)

	ks, _ := kdb.CreateKeySet()

	err = kdb.Get(ks, key)

	Checkf(t, err, "kdb.Open() failed: %v", err)

	t.Log(ks.Len())

	for next := ks.Next(); next != nil; next = ks.Next() {
		t.Log(next)
	}

	t.Log(ks.LookupByName("/bla"))
}

func TestVersion(t *testing.T) {
	kdb := elektra.New()

	key, _ := kdb.CreateKey("/bla")
	err := kdb.Open(key)

	version, err := kdb.Version()

	Checkf(t, err, "kdb.Version() failed: %v", err)
	Assert(t, version != "", "kdb.Version() is empty")
}
