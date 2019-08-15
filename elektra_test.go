package elektra_test

import (
	"testing"

	"github.com/ElektraInitiative/go-elektra"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestElektraOpen(t *testing.T) {
	handle := elektra.New()
	err := handle.Open("user/go-binding-high/test")

	Check(t, err, "could not open elektra instance")
}

func TestLong(t *testing.T) {
	handle := elektra.New()
	_ = handle.Open("user/go-binding-high/test")

	err := handle.SetLong("mylong", 5)
	Check(t, err, "elektra.SetLong() failed")

	val := handle.Long("mylong")

	Assert(t, val == 5, "Long() should be 5")
}

func TestString(t *testing.T) {
	handle := elektra.New()
	_ = handle.Open("user/go-binding-high/test")

	err := handle.SetValue("mystring", "foo")
	Check(t, err, "elektra.SetValue() failed")

	val := handle.Value("mystring")

	Assert(t, val == "foo", "Value() should be foo")
}
