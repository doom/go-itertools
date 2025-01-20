package itertools_test

import (
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/doom/go-itertools"
)

func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

func IntRange(a, b int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for ; a < b; a++ {
			if !yield(a) {
				return
			}
		}
	}
}

func TestItertools_FromSlice(t *testing.T) {
	is := itertools.FromSlice([]int{0, 1, 2, 3, 4})
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	ss := itertools.FromSlice([]string{})
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Map(t *testing.T) {
	ss := itertools.Map(IntRange(0, 5), strconv.Itoa)
	assert.Equal(t, []string{"0", "1", "2", "3", "4"}, slices.Collect(ss))

	ss = itertools.Map(Empty[int](), strconv.Itoa)
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Filter(t *testing.T) {
	ss := itertools.Filter(IntRange(0, 5), func(i int) bool {
		return i%2 == 0
	})
	assert.Equal(t, []int{0, 2, 4}, slices.Collect(ss))

	ss = itertools.Filter(IntRange(0, 5), func(i int) bool { return false })
	assert.Equal(t, []int(nil), slices.Collect(ss))

	ss = itertools.Filter(Empty[int](), func(_ int) bool { return true })
	assert.Equal(t, []int(nil), slices.Collect(ss))
}
