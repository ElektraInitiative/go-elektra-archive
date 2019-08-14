package elektra_test

import (
	"testing"

	"github.com/ElektraInitiative/go-elektra"
	. "github.com/ElektraInitiative/go-elektra/test"
)

func TestElektraOpen(t *testing.T) {
	handle := elektra.New()
	err := handle.Open("/sw/org/myapp/#0/current")

	Check(t, err, "could not open elektra instance")
}

func TestLong(t *testing.T) {
	handle := elektra.New()
	_ = handle.Open("/sw/org/myapp/#0/current")

	val := handle.Long("mylong")

	Assert(t, val == 5, "Long() should be 5")
}

func TestString(t *testing.T) {
	handle := elektra.New()
	_ = handle.Open("/sw/org/myapp/#0/current")

	val := handle.Value("mylong")

	Assert(t, val == "5", "Value() should be 5")
}
