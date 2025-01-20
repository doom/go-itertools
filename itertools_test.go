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

func TestItertools_WithFunc(t *testing.T) {
	is := itertools.WithFunc(func() int { return 1 })
	assert.Equal(t, []int{1, 1, 1, 1, 1}, slices.Collect(itertools.Take(is, 5)))

	i := -1
	is = itertools.WithFunc(func() int { i++; return i })
	assert.Equal(t, []int{0, 1, 2, 3, 4}, slices.Collect(itertools.Take(is, 5)))
}

func TestItertools_Repeat(t *testing.T) {
	ss := itertools.Repeat("a")
	assert.Equal(t, []string{"a", "a", "a", "a", "a"}, slices.Collect(itertools.Take(ss, 5)))
}

func TestItertools_RepeatN(t *testing.T) {
	ss := itertools.RepeatN("a", 5)
	assert.Equal(t, []string{"a", "a", "a", "a", "a"}, slices.Collect(ss))

	ss = itertools.RepeatN("a", 0)
	assert.Equal(t, []string(nil), slices.Collect(ss))
}

func TestItertools_Cycle(t *testing.T) {
	is := itertools.Cycle(IntRange(0, 2))
	assert.Equal(t, []int{0, 1, 0, 1, 0}, slices.Collect(itertools.Take(is, 5)))

	is = itertools.Cycle(Empty[int]())
	assert.Equal(t, []int(nil), slices.Collect(itertools.Take(is, 5)))
}

func TestItertools_Flatten(t *testing.T) {
	is := itertools.Flatten(itertools.Map(IntRange(0, 3), func(v int) iter.Seq[int] {
		return itertools.RepeatN(v, 2)
	}))
	assert.Equal(t, []int{0, 0, 1, 1, 2, 2}, slices.Collect(is))

	is = itertools.Flatten(Empty[iter.Seq[int]]())
	assert.Equal(t, []int(nil), slices.Collect(is))
}

func TestItertools_All(t *testing.T) {
	a := itertools.All(IntRange(0, 3), func(v int) bool { return v >= 0 })
	assert.Equal(t, true, a)

	a = itertools.All(IntRange(0, 3), func(v int) bool { return v > 1 })
	assert.Equal(t, false, a)

	a = itertools.All(IntRange(0, 3), func(v int) bool { return v < 2 })
	assert.Equal(t, false, a)
}

func TestItertools_Any(t *testing.T) {
	a := itertools.Any(IntRange(0, 3), func(v int) bool { return v >= 0 })
	assert.Equal(t, true, a)

	a = itertools.Any(IntRange(0, 3), func(v int) bool { return v > 2 })
	assert.Equal(t, false, a)

	a = itertools.Any(IntRange(0, 3), func(v int) bool { return v%2 == 1 })
	assert.Equal(t, true, a)
}

func TestItertools_None(t *testing.T) {
	a := itertools.None(IntRange(0, 3), func(v int) bool { return v < 0 })
	assert.Equal(t, true, a)

	a = itertools.None(IntRange(0, 3), func(v int) bool { return v < 2 })
	assert.Equal(t, false, a)

	a = itertools.None(IntRange(0, 3), func(v int) bool { return v%2 == 1 })
	assert.Equal(t, false, a)
}
