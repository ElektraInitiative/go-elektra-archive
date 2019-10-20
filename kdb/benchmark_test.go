package kdb

import (
	"fmt"
	"testing"

	. "go.libelektra.org/test"
)

func setupTestData(b *testing.B, count int) *CKeySet {
	b.Helper()

	ks := NewKeySet()

	for n := 0; n < count; n++ {
		k, err := NewKey(fmt.Sprintf("proc/tests/go/elektra/benchmark/iterator/callback/%03d", n))
		Checkf(b, err, "kdb.NewKey() failed: %v", err)

		ks.AppendKey(k)
	}

	b.ResetTimer()
	return ks.(*CKeySet)
}

func BenchmarkKeySetExternalCallbackIterator(b *testing.B) {
	ks := setupTestData(b, 1000)

	for n := 0; n < b.N; n++ {
		ks.ForEach(func(k Key) {
		})
	}
}

func BenchmarkKeySetInternalCallbackIterator(b *testing.B) {
	ks := setupTestData(b, 1000)

	for n := 0; n < b.N; n++ {
		ks.forEachInternal(func(k Key) {
		})
	}
}

func BenchmarkKeySetSliceRangeIterator(b *testing.B) {
	ks := setupTestData(b, 1000)

	for n := 0; n < b.N; n++ {
		ksSlice := ks.ToSlice()

		for range ksSlice {
		}
	}
}
