# Package intlist

## Intro

This package contains a couple of functions for parsing strings containing a
representation of an integer list to produce a slice or an iterator to iterate
through the integers. I couldn't find any public project providing this
functionality at that time. This was for a personal fringe use case, but may be
of use to others.

I was also using this to play with variations of testing code and iterators.

## Thoughts on Iterators and why I implemented the one I did

I searched for various styles of Iterators in "go". My first search only found
goroutines/channels, closures, callbacks, and stateful iterators using
New/Next/Value. I first implemented the New/Next/Value style. The samples I ran
across did not address certain sequences of calls that a developer might
erroneously make. (E.g., calling Value before the first call to Next to set the
value) Sure, the developer should know better, but it seems best to let them
know they are doing something wrong. I added a panic for that. I did allow
Value to be called multiple times in a row as someone might do that in a loop
while I expect that the normal usage would be to call it once and save the
result.

I really wanted something that combined both the Next and Value call. I first
modified my code to have a GetNext function that returned both a value and a
bool. The bool would indicate no more items and that the value was invalid.
This made one less call to implement and got rid of the case of a developer
calling Value before Next.

I was going to switch back to the first implementation since I was hesitant to
add yet another style but then ran across this Google Cloud API link
https://github.com/googleapis/google-cloud-go/wiki/Iterator-Guidelines. The
Next function there was doing close to what I was doing with GetNext. I changed
my GetNext to Next and changed the bool to an error to match up with that
usage. The use of ErrDone rather than Done was made since the code checking I
was using in VS Code gave a suggestion to make names of error be in the form of
ErrMyErrorName.

I added panics for the two combinations of calls that should be avoided. Those
are using an Iterator that had an error when parsing the specification and
calling Next after Next already returned ErrDone.

I added tests to make sure that the panics had the proper info to make it clear
to a developer using this code what the issue was.

See the embedded documentation in the "go" files.