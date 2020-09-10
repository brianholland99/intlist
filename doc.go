// Copyright 2020 Brian E. Holland. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package intlist supports a string notation specifying a series of integers.
//
// This was written to support a data-driven text file entered by humans that
// contained a mix of integers and sequences. This format made it easy to enter
// and to visually recognize a sequence of consecutive integers.
//
// Format:
//   - Comma-separated expressions of integers or integer sequences.
//   - An empty string indicates an empty list.
//   - Sequences are consecutive integers notated by two endpoints
//     separated by an ellipsis and includes both endpoints.
//   - Both increasing and decreasing sequences are supported.
//
// Examples:
//   spec = "4,6,10...15" --> [4, 6, 10, 11, 12, 13, 14, 15]
//   spec = "4,12...8,-3" --> [4, 12, 11, 10, 9, 8, -3]
//
// There are two supported use cases; creating an int slice and an Iterator to
// produce the ints as needed.
//
// "Parse" will parse a string and return a integer slice. This is useful when
// a slice is wanted and the size of the result is not too large.
//
// "NewIterator" / "Next" / "Err" functions - provide the functionality
// necessary to iterate through the list of integers. This may be especially
// useful when the resulting list is too huge or when it is possible to stop
// before using the whole list.
//
// Example of iterator usage:
//
//    it := intlist.NewIterator("1...1000,1030...1014,2000")
//    if it.Err() != nil {
//        // Handle error. Don't fall through to the for loop as the Next call
//        // will panic when the Iterator has an error.
//    }
//    for {
//        val, err := it.Next()
//        if err == intlist.ErrDone() {
//            break // Invalid val. Don't use it and stop iterating.
//        }
//        fmt.Println(val) // Or whatever processing is to be done.
//    }
package intlist
