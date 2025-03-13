package itertools_test

import (
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/doom/go-itertools"
)

func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

func Empty2[V, W any]() iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {}
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

func TestItertools_MapFromSeq2(t *testing.T) {
	is := itertools.MapFromSeq2(itertools.FromMap(map[int]int{0: 1, 1: 2, 2: 3, 3: 4}), func(a, b int) int { return a + b })
	assert.ElementsMatch(t, []int{0 + 1, 1 + 2, 2 + 3, 3 + 4}, slices.Collect(is))

	is = itertools.MapFromSeq2(Empty2[int, int](), func(a, b int) int { return a + b })
	assert.ElementsMatch(t, []int{}, slices.Collect(is))
}

func TestItertools_MapToSeq2(t *testing.T) {
	is := itertools.MapToSeq2(IntRange(0, 5), func(v int) (string, int) { return strconv.Itoa(v), v })
	assert.Equal(t, map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4}, maps.Collect(is))

	is = itertools.MapToSeq2(Empty[int](), func(v int) (string, int) { return strconv.Itoa(v), v })
	assert.Equal(t, map[string]int{}, maps.Collect(is))
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

func TestItertools_Reduce(t *testing.T) {
	n := itertools.Reduce(IntRange(0, 5), func(a, b int) int {
		return a + b
	}, 0)
	assert.Equal(t, 0+1+2+3+4, n)

	n = itertools.Reduce(IntRange(0, 5), func(a, b int) int {
		return a + b
	}, 123)
	assert.Equal(t, 123+0+1+2+3+4, n)

	n = itertools.Reduce(Empty[int](), func(a, b int) int {
		return a + b
	}, 123)
	assert.Equal(t, 123, n)
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

func TestItertools_Min(t *testing.T) {
	a, ok := itertools.Min(IntRange(0, 3))
	assert.Equal(t, true, ok)
	assert.Equal(t, 0, a)

	a, ok = itertools.Min(itertools.FromSlice([]int{4, 3, 2, -1, 0}))
	assert.Equal(t, true, ok)
	assert.Equal(t, -1, a)

	_, ok = itertools.Min(Empty[int]())
	assert.Equal(t, false, ok)
}

func TestItertools_MinFunc(t *testing.T) {
	a, ok := itertools.MinFunc(itertools.FromSlice([]string{"ghi", "abc", "def"}), strings.Compare)
	assert.Equal(t, true, ok)
	assert.Equal(t, "abc", a)

	_, ok = itertools.MinFunc(Empty[string](), strings.Compare)
	assert.Equal(t, false, ok)
}

func TestItertools_Max(t *testing.T) {
	a, ok := itertools.Max(IntRange(0, 3))
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, a)

	a, ok = itertools.Max(itertools.FromSlice([]int{4, 3, 2, 5, 0}))
	assert.Equal(t, true, ok)
	assert.Equal(t, 5, a)

	_, ok = itertools.Max(Empty[int]())
	assert.Equal(t, false, ok)
}

func TestItertools_MaxFunc(t *testing.T) {
	a, ok := itertools.MaxFunc(itertools.FromSlice([]string{"abc", "ghi", "def"}), strings.Compare)
	assert.Equal(t, true, ok)
	assert.Equal(t, "ghi", a)

	_, ok = itertools.MaxFunc(Empty[string](), strings.Compare)
	assert.Equal(t, false, ok)
}

func TestItertools_InterleaveShortest(t *testing.T) {
	ss := itertools.InterleaveShortest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, []string{"abc", "def", "ghi", "jkl"}, slices.Collect(ss))

	ss = itertools.InterleaveShortest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def"}),
	)
	assert.Equal(t, []string{"abc", "def", "ghi"}, slices.Collect(ss))

	ss = itertools.InterleaveShortest(
		itertools.FromSlice([]string{"abc"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, []string{"abc", "def"}, slices.Collect(ss))
}

func TestItertools_InterleaveLongest(t *testing.T) {
	ss := itertools.InterleaveLongest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, []string{"abc", "def", "ghi", "jkl"}, slices.Collect(ss))

	ss = itertools.InterleaveLongest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def"}),
	)
	assert.Equal(t, []string{"abc", "def", "ghi"}, slices.Collect(ss))

	ss = itertools.InterleaveLongest(
		itertools.FromSlice([]string{"abc"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, []string{"abc", "def", "jkl"}, slices.Collect(ss))

	ss = itertools.InterleaveLongest(
		itertools.FromSlice([]string{"abc", "ghi", "jkl"}),
		itertools.FromSlice([]string{"def"}),
	)
	assert.Equal(t, []string{"abc", "def", "ghi", "jkl"}, slices.Collect(ss))
}

func TestItertools_ZipShortest(t *testing.T) {
	ss := itertools.ZipShortest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, map[string]string{"abc": "def", "ghi": "jkl"}, maps.Collect(ss))

	ss = itertools.ZipShortest(
		itertools.FromSlice([]string{"abc", "ghi"}),
		itertools.FromSlice([]string{"def"}),
	)
	assert.Equal(t, map[string]string{"abc": "def"}, maps.Collect(ss))

	ss = itertools.ZipShortest(
		itertools.FromSlice([]string{"abc"}),
		itertools.FromSlice([]string{"def", "jkl"}),
	)
	assert.Equal(t, map[string]string{"abc": "def"}, maps.Collect(ss))

	ss = itertools.ZipShortest(
		itertools.FromSlice([]string{}),
		itertools.FromSlice([]string{"abc", "ghi", "jkl"}),
	)
	assert.Equal(t, map[string]string{}, maps.Collect(ss))
}

func TestItertools_ChunkBy(t *testing.T) {
	iss := itertools.ChunkBy(IntRange(-2, 2), func(i int) bool {
		return i < 0
	})
	collected := slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 2, len(collected))
	require.Equal(t, []int{-2, -1}, collected[0])
	require.Equal(t, []int{0, 1}, collected[1])

	iss = itertools.ChunkBy(Empty[int](), func(i int) bool {
		return i < 0
	})
	collected = slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 0, len(collected))

	iss = itertools.ChunkBy(IntRange(-2, 2), func(i int) bool {
		return true
	})
	collected = slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 1, len(collected))
	require.Equal(t, []int{-2, -1, 0, 1}, collected[0])

	iss = itertools.ChunkBy(IntRange(0, 5), func(i int) bool {
		return i < 0
	})
	collected = slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 1, len(collected))
	require.Equal(t, []int{0, 1, 2, 3, 4}, collected[0])

	iss = itertools.ChunkBy(IntRange(0, 5), func(i int) int {
		return i
	})
	collected = slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 5, len(collected))
	for i := range 5 {
		require.Equal(t, []int{i}, collected[i])
	}

	iss = itertools.ChunkBy(IntRange(0, 5), func(i int) int {
		return i % 2
	})
	collected = slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 5, len(collected))
	for i := range 5 {
		require.Equal(t, []int{i}, collected[i])
	}
}

func TestItertools_Chunks(t *testing.T) {
	iss := itertools.Chunks(IntRange(0, 10), 2)
	collected := slices.Collect(itertools.Map(iss, slices.Collect))
	require.Equal(t, 5, len(collected))
	for i := range 5 {
		require.Equal(t, []int{i * 2, i*2 + 1}, collected[i])
	}
}
