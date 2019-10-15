package kdb

import (
	"fmt"
	"testing"

	. "github.com/ElektraInitiative/go-elektra/test"
)

func setupInMemoryKeySet(b *testing.B, count int) *ckeySet {
	b.Helper()

	keys := []Key{}

	ks := NewKeySet(keys...)

	for n := 0; n < count; n++ {
		k, err := NewKey(fmt.Sprintf("proc/tests/go/elektra/benchmark/iterator/callback/%03d", n))
		Checkf(b, err, "kdb.NewKey() failed: %v", err)

		ks.AppendKey(k)
	}

	return ks.(*ckeySet)
}

func BenchmarkKeySetExternalCallbackIterator(b *testing.B) {
	ks := setupInMemoryKeySet(b, 1000)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ks.Each(func(k Key) bool {
			return true
		})
	}
}

func BenchmarkKeySetInternalCallbackIterator(b *testing.B) {
	ks := setupInMemoryKeySet(b, 1000)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ks.loopInternal(func(k Key) bool {
			return true
		})
	}
}

func BenchmarkKeySetSliceRangeIterator(b *testing.B) {
	ks := setupInMemoryKeySet(b, 1000)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ksSlice := ks.Slice()

		for range ksSlice {
		}
	}
}
