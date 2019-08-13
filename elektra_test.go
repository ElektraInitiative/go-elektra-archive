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
