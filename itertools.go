// Package itertools implements common utilities to work with Go 1.23's iterators.
package itertools

import (
	"cmp"
	"iter"
)

// FromSlice returns an iterator yielding all the values from vs.
func FromSlice[V any](vs []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vs {
			if !yield(v) {
				return
			}
		}
	}
}

// Map returns an iterator that will yield values from seq after transforming them using f.
func Map[V any, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Filter returns an iterator that will yield values from seq only if they pass p.
func Filter[V any](seq iter.Seq[V], p func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if p(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// TakeWhile returns an iterator that will yield values from seq as long as they pass p.
// The iterator stops when it encounters a value that does not pass p.
func TakeWhile[V any](seq iter.Seq[V], p func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !p(v) || !yield(v) {
				return
			}
		}
	}
}

// Take returns an iterator that will yield the n first values from seq.
func Take[V any](seq iter.Seq[V], n uint) iter.Seq[V] {
	return TakeWhile(seq, func(_ V) bool {
		if n == 0 {
			return false
		}
		n--
		return true
	})
}

// DropWhile returns an iterator that will drop values from seq as long as they pass p.
// The iterator yields the remaining values when it encounters the first value that does not pass p.
func DropWhile[V any](seq iter.Seq[V], p func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		for v, ok := next(); ok; v, ok = next() {
			if p(v) {
				continue
			}

			if !yield(v) {
				return
			}
			break
		}

		for v, ok := next(); ok; v, ok = next() {
			if !yield(v) {
				return
			}
		}
	}
}

// Drop returns an iterator that will drop the n first values from seq.
func Drop[V any](seq iter.Seq[V], n uint) iter.Seq[V] {
	return DropWhile(seq, func(_ V) bool {
		if n == 0 {
			return false
		}
		n--
		return true
	})
}

// Chain returns an iterator that will first yield all the values from seq1, then all the values from seq2.
func Chain[V any](seq1, seq2 iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for s := range seq1 {
			if !yield(s) {
				return
			}
		}
		for s := range seq2 {
			if !yield(s) {
				return
			}
		}
	}
}

// WithFunc returns an iterator yielding values obtained by indefinitely calling f.
func WithFunc[V any](f func() V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			if !yield(f()) {
				return
			}
		}
	}
}

// Repeat returns an iterator that will indefinitely yield v.
func Repeat[V any](v V) iter.Seq[V] {
	return WithFunc(func() V { return v })
}

// RepeatN works like Repeat, but returns an iterator that stops after yielding n values.
func RepeatN[V any](v V, n uint) iter.Seq[V] {
	return Take(Repeat(v), n)
}

// Cycle returns an iterator that cycles through seq indefinitely.
// Values from seq are progressively accumulated into a slice during the first cycle,
// and reused for the next cycles.
func Cycle[V any](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		var vs []V

		for v := range seq {
			if !yield(v) {
				return
			}
			vs = append(vs, v)
		}

		for len(vs) > 0 {
			for _, v := range vs {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Flatten returns an iterator that yields each value from a nested iterator.
func Flatten[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for s := range seq {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// All reports whether all values yielded by seq pass p.
// All is short-circuiting, i.e. it will stop when it reaches a value that does not pass p.
func All[V any](seq iter.Seq[V], p func(V) bool) bool {
	for v := range seq {
		if !p(v) {
			return false
		}
	}
	return true
}

// Any reports whether any value yielded by seq passes p.
// Any is short-circuiting, i.e. it will stop when it reaches a value that passes p.
func Any[V any](seq iter.Seq[V], p func(V) bool) bool {
	for v := range seq {
		if p(v) {
			return true
		}
	}
	return false
}

// None reports whether no value yielded by seq passes p.
// None is short-circuiting, i.e. it will stop when it reaches a value that passes p.
func None[V any](seq iter.Seq[V], p func(V) bool) bool {
	return !Any(seq, p)
}

// MinFunc returns the minimum value yielded by seq, comparing values using cmp.
// If no values are yielded by seq, a zero-value is returned and the second return value is false.
// If there is more than one minimal element according to the cmp function, MinFunc returns the first one.
func MinFunc[V any](seq iter.Seq[V], cmp func(V, V) int) (V, bool) {
	next, stop := iter.Pull(seq)
	defer stop()

	minV, ok := next()
	if !ok {
		return minV, false
	}

	for v, ok := next(); ok; v, ok = next() {
		if cmp(v, minV) < 0 {
			minV = v
		}
	}

	return minV, true
}

// Min returns the minimum value yielded by seq.
// If no values are yielded by seq, a zero-value is returned and the second return value is false.
// If there is more than one minimal element according to the cmp function, Min returns the first one.
func Min[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	return MinFunc(seq, cmp.Compare)
}

// MaxFunc returns the minimum value yielded by seq, comparing values using cmp.
// If no values are yielded by seq, a zero-value is returned and the second return value is false.
// If there is more than one minimal element according to the cmp function, MaxFunc returns the first one.
func MaxFunc[V any](seq iter.Seq[V], cmp func(V, V) int) (V, bool) {
	next, stop := iter.Pull(seq)
	defer stop()

	maxV, ok := next()
	if !ok {
		return maxV, false
	}

	for v, ok := next(); ok; v, ok = next() {
		if cmp(v, maxV) > 0 {
			maxV = v
		}
	}

	return maxV, true
}

// Max returns the minimum value yielded by seq.
// If no values are yielded by seq, a zero-value is returned and the second return value is false.
// If there is more than one minimal element according to the cmp function, Max returns the first one.
func Max[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	return MaxFunc(seq, cmp.Compare)
}
