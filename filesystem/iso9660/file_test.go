package iso9660_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

func TestFileRead(t *testing.T) {
	f, content := iso9660.GetTestFile(t)

	b, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("received unexpected error when reading: %v", err)
	}
	bString := string(b)
	if bString != content {
		t.Errorf("Mismatched content:\nActual: '%s'\nExpected: '%s'", bString, content)
	}
}

func TestPassesStdTestReader(t *testing.T) {
	file, content := iso9660.GetTestFile(t)

	if err := iotest.TestReader(file, []byte(content)); err != nil {
		t.Errorf("Failed iotest.TestReader: %v", err)
	}

	largeFile, size := iso9660.GetLargeTestFile(t)
	bytes := make([]byte, size, size)
	if err := iotest.TestReader(largeFile, bytes); err != nil {
		// only print the first line because the error contains the entire file content
		t.Errorf("large file failed iotest.TestReader: %s", strings.Split(err.Error(), "\n")[0])
	}
}

func TestLargeFileCopy(t *testing.T) {
	f, size := iso9660.GetLargeTestFile(t)

	copied, err := io.Copy(ioutil.Discard, f)
	if err != nil {
		t.Errorf("received unexpected error when copying: %v", err)
	}
	if copied != int64(size) {
		t.Errorf("copied %d bytes instead of expected %d", copied, size)
	}
}

func TestFileWrite(t *testing.T) {
	// pretty simple: never should be able to write as it is a read-only filesystem
	f := &iso9660.File{}
	b := make([]byte, 8, 8)
	written, err := f.Write(b)
	if err == nil {
		t.Errorf("received no error when should have been prevented from writing")
	}
	if written != 0 {
		t.Errorf("wrote %d bytes instead of expected %d", written, 0)
	}
}
