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

func TestItertools_TakeWhile(t *testing.T) {
	is := itertools.TakeWhile(IntRange(0, 5), func(i int) bool { return i < 3 })
	assert.Equal(t, []int{0, 1, 2}, slices.Collect(is))

	is = itertools.TakeWhile(IntRange(0, 5), func(i int) bool { return true })
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	is = itertools.TakeWhile(IntRange(0, 5), func(i int) bool { return false })
	assert.Equal(t, []int(nil), slices.Collect(is))

	ss := itertools.TakeWhile(Empty[string](), func(i string) bool { return true })
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Take(t *testing.T) {
	is := itertools.Take(IntRange(0, 5), 3)
	assert.Equal(t, []int{0, 1, 2}, slices.Collect(is))

	is = itertools.Take(IntRange(0, 5), 10)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	is = itertools.Take(IntRange(0, 5), 0)
	assert.Equal(t, []int(nil), slices.Collect(is))

	ss := itertools.Take(Empty[string](), 5)
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_DropWhile(t *testing.T) {
	is := itertools.DropWhile(IntRange(0, 5), func(i int) bool { return i < 3 })
	assert.Equal(t, []int{3, 4}, slices.Collect(is))

	is = itertools.DropWhile(IntRange(0, 5), func(i int) bool { return true })
	assert.Equal(t, []int(nil), slices.Collect(is))

	is = itertools.DropWhile(IntRange(0, 5), func(i int) bool { return false })
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	ss := itertools.DropWhile(Empty[string](), func(i string) bool { return false })
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Drop(t *testing.T) {
	is := itertools.Drop(IntRange(0, 5), 3)
	assert.Equal(t, []int{3, 4}, slices.Collect(is))

	is = itertools.Drop(IntRange(0, 5), 10)
	assert.Equal(t, []int(nil), slices.Collect(is))

	is = itertools.Drop(IntRange(0, 5), 0)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	ss := itertools.Drop(Empty[string](), 0)
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Chain(t *testing.T) {
	is := itertools.Chain(Empty[int](), Empty[int]())
	assert.Equal(t, []int(nil), slices.Collect(is))

	is = itertools.Chain(Empty[int](), IntRange(0, 5))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	is = itertools.Chain(IntRange(0, 5), Empty[int]())
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(is))

	is = itertools.Chain(IntRange(0, 5), IntRange(5, 10))
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, slices.Collect(is))
}
