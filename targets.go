package lynd

import (
	"io/ioutil"
	"os"
)

// tempPrefix is prepended to any temporary filename.
const tempPrefix = "lynd-"

// WithoutIO implements an io.ReadWriteCloser, that does nothing.
type WithoutIO struct{}

// Close is a noop.
func (t *WithoutIO) Close() error { return nil }

// Read is not implemented.
func (t *WithoutIO) Read(p []byte) (n int, err error) {
	return 0, errNotImplemented
}

// Write is not implemented.
func (t *WithoutIO) Write(p []byte) (n int, err error) {
	return 0, errNotImplemented
}

// Done is a target, that justs always exists. It does not support any IO.
type Done struct {
	WithoutIO
}

// Name, to satisfy the interface, not really useful.
func (t *Done) Name() string { return "<done>" }

// Exists always returns true.
func (t *Done) Exists() bool { return true }

// Failed is a target, that never exists.It does not support any IO.
type Failed struct {
	WithoutIO
}

func (t *Failed) Name() string { return "<failed>" }
func (t *Failed) Exists() bool { return false }

// File is a local file target. Most common type of target. It encapsulates
// basic atomicity by writing to a temporary file first.
type File struct {
	Path  string
	f     *os.File
	isTmp bool
}

// Read from io.Reader.
// TODO(miku): Check-lock-check.
func (t *File) Read(p []byte) (n int, err error) {
	if t.f == nil {
		f, err := os.Open(t.Path)
		if err != nil {
			return 0, err
		}
		t.f = f
	}
	return t.f.Read(p)
}

// Write from io.Writer.
// TODO(miku): Check-lock-check.
func (t *File) Write(p []byte) (n int, err error) {
	if t.f == nil {
		f, err := ioutil.TempFile("", tempPrefix)
		if err != nil {
			return 0, err
		}
		t.f = f
		t.isTmp = true
	}
	return t.f.Write(p)
}

// Close from io.Closer.
// TODO(miku): Check-lock-check.
func (t *File) Close() error {
	if t.f != nil {
		err := t.f.Close()
		if err != nil {
			return err
		}
		if t.isTmp {
			err = rename(t.f.Name(), t.Path)
			if err != nil {
				return err
			}
			t.isTmp = false
		}
		t.f = nil
	}
	return nil
}

// Name should returns the full path of the file.
func (t *File) Name() string { return t.Path }

// Exists returns true, if we could actually stat the file.
// TODO(miku): be more strict.
func (t *File) Exists() bool {
	_, err := os.Stat(t.Name())
	return err == nil
}

// String reports the path.
func (t *File) String() string {
	return t.Path
}
