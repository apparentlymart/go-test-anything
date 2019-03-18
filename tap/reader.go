package tap

import (
	"bufio"
	"io"
)

// Read is a convenience wrapper around constructing a Reader, reading all of
// its results, and constructing a report. A caller that doesn't need streaming
// access to results should use this for simplicity.
func Read(r io.Reader) (*RunReport, error) {
	tr := NewReader(r)
	return tr.ReadAll()
}

// Reader consumes TAP output from an io.Reader and provides a pull-based API
// to access that output.
type Reader struct {
	r  io.Reader
	sc *bufio.Scanner

	plan    *Plan
	nextNum int
	results map[int]*Report
}

// NewReader creates a new Reader that parses TAP output from the given
// io.Reader.
func NewReader(r io.Reader) *Reader {
	sc := bufio.NewScanner(r)
	return &Reader{
		r:  r,
		sc: sc,

		nextNum: 1,
	}
}

// Read will block until either a new test report is available or until there
// are no more reports to read (either due to successful end of file or via an
// error). The result is true if a new test report was found, or false if there
// are no more reports to read.
func (r *Reader) Read() bool {
	// TODO
	return false
}

// ReadAll is a convenience wrapper around calling Read in a loop for callers
// that don't need streaming TAP output. It will consume all of the results,
// update any other status, and then return the error from the reader if there
// is one.
func (r *Reader) ReadAll() (*RunReport, error) {
	for r.Read() {
	}
	return r.Report(), r.Err()
}

// Report creates and returns a RunReport object describing the overall outcome
// of a test run. The returned object will be incomplete if this method is called
// before the test run has finished.
func (r *Reader) Report() *RunReport {
	// TODO
	return nil
}

// Err returns an error that was encountered during reading, if any. Call this
// after Read stops returning true to learn if the reason was due to the end
// being reached (in which case Err returns nil) or some other problem.
func (r *Reader) Err() error {
	return r.sc.Err()
}
