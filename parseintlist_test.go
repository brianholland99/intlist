// Copyright 2020 Brian E. Holland. All rights reserved.
// The use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package intlist

import (
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type parseTest struct {
	in  string
	out []int
	err error
}

var parseTests = []parseTest{
	// Good cases
	{"1234", []int{1234}, nil},                              // Single Int
	{"6...9", []int{6, 7, 8, 9}, nil},                       // Single seq
	{"-1...2,6...4", []int{-1, 0, 1, 2, 6, 5, 4}, nil},      // Two seq
	{"", []int{}, nil},                                      // Empty list
	{"1...3,7,5...3,9", []int{1, 2, 3, 7, 5, 4, 3, 9}, nil}, // Ints and Seqs
	// Error cases
	{"   12, 4, 9...6", nil, strconv.ErrSyntax}, // Whitespace
	{"-2...-4...-6,12", nil, strconv.ErrSyntax}, // Multiple ... in one item
	{"3.5,12", nil, strconv.ErrSyntax},          // Non-integer
	{"3.9...5", nil, strconv.ErrSyntax},         // Seq. start - non-integer
	{"2...5.4", nil, strconv.ErrSyntax},         // Seq. end - non-integer
}

// This tests Parse and indirectly tests most of the Iterator code.
func TestParse(t *testing.T) {
	for _, test := range parseTests {
		out, err := Parse(test.in)
		if !cmp.Equal(out, test.out) || !errors.Is(err, test.err) {
			t.Errorf("Parse(%q) = (%v), (%v) -- wanted (%v), (%v)",
				test.in, out, err, test.out, test.err)
		}
	}
}

// The remaining tests check proper response to misuse of Iterator functions
// by callers and also the proper return of ErrDone by Next().

func TestUseOfBadIterator(t *testing.T) {
	it := NewIterator("2.3") // Non-int error so Iterator is invalid.
	// Ignore error to verify that value
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Next() of invalid Iterator did not panic.")
		} else {
			// Make sure correct panic message is given to developer.
			exp := "Next() called on invalid iterator."
			if err != exp {
				t.Errorf("Wrong panic message - Should be \"" + exp + "\"")
			}
		}
	}()
	_, _ = it.Next()
}

func TestUseOfNextWithErrDone(t *testing.T) {
	it := NewIterator("") // Empty list
	// First check that ErrDone is returned
	_, err := it.Next()
	if err != ErrDone {
		t.Errorf("ErrDone not returned on empty list")
	}
	// Ignore error to see if next call to Next() panics as expected.
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Next() of bad Iterator did not panic.")
		} else {
			// Make sure correct panic message is given to developer.
			exp := "Next() called again after returning ErrDone."
			if err != exp {
				t.Errorf("Wrong panic message - Should be \"" + exp + "\"")
			}
		}
	}()
	_, _ = it.Next()
}
