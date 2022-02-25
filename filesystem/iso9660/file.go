package iso9660

import (
	"fmt"
	"io"
	"os"
)

// File represents a single file in an iso9660 filesystem
//  it is NOT used when working in a workspace, where we just use the underlying OS
type File struct {
	*directoryEntry
	isReadWrite bool
	isAppend    bool
	offset      int64
	closed      bool
}

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF
// reads from the last known offset in the file from last read or write
// use Seek() to set at a particular point
func (fl *File) Read(b []byte) (int, error) {
	if fl == nil || fl.closed {
		return 0, os.ErrClosed
	}

	locationBytes := int64(fl.location) * fl.filesystem.blocksize
	isoOffset := locationBytes + fl.offset
	isoLimit := locationBytes + int64(fl.size)

	if isoOffset >= isoLimit {
		return 0, io.EOF
	}
	if max := isoLimit - isoOffset; int64(len(b)) > max {
		b = b[0:max]
	}
	n, err := fl.filesystem.file.ReadAt(b, isoOffset)
	fl.offset += int64(n)
	return n, err
}

// Write writes len(b) bytes to the File.
//  you cannot write to an iso, so this returns an error
func (fl *File) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("Cannot write to a read-only iso filesystem")
}

// Seek set the offset to a particular point in the file
func (fl *File) Seek(offset int64, whence int) (int64, error) {
	if fl == nil || fl.closed {
		return 0, os.ErrClosed
	}
	newOffset := int64(0)
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekEnd:
		newOffset = int64(fl.size) + offset
	case io.SeekCurrent:
		newOffset = fl.offset + offset
	}
	if newOffset < 0 {
		return fl.offset, fmt.Errorf("Cannot set offset %d before start of file", offset)
	}
	fl.offset = newOffset
	return fl.offset, nil
}

func (fl *File) Location() uint32 {
	return fl.location
}

// Close close the file
func (fl *File) Close() error {
	fl.closed = true
	return nil
}
