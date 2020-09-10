// Copyright 2020 Brian E. Holland. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package intlist

import (
	"errors"
	"strconv"
	"strings"
)

// ErrDone is returned when iteration when completed.
var ErrDone = errors.New("no more items in iterator")

// Seq is used to denote both single integers and sequences of integers. A
// single integer is denoted by next == last.
type seq struct {
	next int // Next value to retrieve
	last int // Last value in sequence
	step int // Direction (I.e., +1 for increasing, -1 for decreasing)
}

// Iterator is the state for generating integers from an intlist description.
type Iterator struct {
	seqs []seq // Remaining sequences to handle
	err  error // Error in creating or ErrDone if Iterator finishes.
}

// NewIterator validates the specification and sets the state for iteration.
//
// The "spec" parameter is parsed as a string containing a comma-separated list
// of integers and integer sequences. Sequences are defined by two integers
// separated by an ellipsis (E.g., "3...100") and include both endpoints. See
// overall documentation for a more detailed definition of the format.
//
//   NewIterator("1,2,21,50...54,57...61") ->
//       [1 2 21 50 51 52 53 54 57 58 59 60 61]
//
// Potential errors set in state during creation of an Iterator:
//
//   strconv.ErrSyntax - Error parsing integer or sequence notation
//   strconv.ErrRange - Integer out of range
func NewIterator(spec string) *Iterator {
	var err error  // First error encountered, if any
	var seqs []seq // Sequences built during parsing
	const fnNewIterator = "NewIterator"
	items := strings.Split(spec, ",") // Break into comma-separated items
	if len(items) == 1 && items[0] == "" {
		seqs = []seq{} // Handle empty list case
	} else {
		// Handle non-empty list case
		for _, item := range items {
			var itemData seq
			parts := strings.Split(item, "...")
			switch len(parts) {
			// First error encountered will be handled after switch.
			case 1: // Single value (E.g., "265")
				// Treat as sequence of one to simplify iteration routine.
				itemData.next, err = strconv.Atoi(parts[0])
				itemData.last = itemData.next
			case 2: // Sequence
				itemData.next, err = strconv.Atoi(parts[0])
				if err != nil {
					break
				}
				itemData.last, err = strconv.Atoi(parts[1])
				if err != nil {
					break
				}
				if itemData.next < itemData.last {
					itemData.step = 1 // Increasing sequence
				} else {
					itemData.step = -1 // Decreasing sequence
				}
			default: // Multiple "..." in an item
				err = &strconv.NumError{
					Func: fnNewIterator,
					Num:  item,
					Err:  strconv.ErrSyntax,
				}
			}
			if err != nil {
				seqs = nil
				break
			}
			seqs = append(seqs, itemData)
		}
	}
	return &Iterator{
		seqs: seqs,
		err:  err,
	}
}

// Next returns the next integer if not done and an error to indicate if done.
//
// If ErrDone is returned the integer is not valid and there are no more items.
//
// It will panic for the following avoidable cases:
//   - Next called on invalid iterator.
//   - Next called after previous call to Next() returned ErrDone.
func (i *Iterator) Next() (int, error) {
	if i.err != nil {
		if i.err == ErrDone {
			// Caller was already informed that iterator was done.
			panic("Next() called again after returning ErrDone.")
		}
		panic("Next() called on invalid iterator.")
	}
	if len(i.seqs) == 0 {
		i.err = ErrDone
		return 0, ErrDone
	}
	item := &i.seqs[0] // Current sequence being handled
	val := item.next
	if val == item.last {
		// Done with this item. Remove handled expression.
		i.seqs = i.seqs[1:]
	} else {
		item.next += item.step // Move to next value in sequence.
	}
	return val, nil
}

// Err returns any error that occured when creating this Iterator or ErrDone
// if a previous Next call returned that to indicate that the end of the
// iteration occurred.
func (i *Iterator) Err() error {
	return i.err
}

// Parse will return an int slice represented by the passed specification.
//
// The "spec" parameter is parsed as containing a comma-separated list of
// integers and integer sequences. Sequences are defined by two integers
// separated by an ellipsis (E.g., "3...100") and include both endpoints. See
// overall documentation for a more detailed definition of the format.
//
// Parse("1,2,21,50...54,61..57") ->
//       [1 2 21 50 51 52 53 54 61 60 59 58 57], nil
//
// Potential errors returned:
//
//   strconv.ErrSyntax - Error parsing integer or sequence notation.
//   strconv.ErrRange - Integer out of range
func Parse(spec string) ([]int, error) {
	result := []int{}
	it := NewIterator(spec)
	if it.Err() != nil {
		return nil, it.Err()
	}
	for {
		val, err := it.Next()
		if err == ErrDone {
			break
		}
		result = append(result, val)
	}
	return result, nil
}
