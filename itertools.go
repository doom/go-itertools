// Package itertools implements common utilities to work with Go 1.23's iterators.
package itertools

import (
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
