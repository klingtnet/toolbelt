package toolbelt

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type BufferedReader struct {
	bufSize int64
	bufA    []byte
	bufB    []byte
}

func NewBufferedReader(bufSize int64) *BufferedReader {
	if bufSize < 1 {
		panic("bufSize must be > 0")
	}
	return &BufferedReader{
		bufSize: bufSize,
		bufA:    make([]byte, bufSize),
		bufB:    make([]byte, bufSize),
	}
}

const DefaultReadBufferSize = 4 * 1024

var DefaultBufferedReader = NewBufferedReader(DefaultReadBufferSize)

// FilesEqual returns true if both files are equal, i.e. if both store the same bytes.
func FilesEqual(fileA, fileB string) (bool, error) {
	a, err := os.Open(fileA)
	if err != nil {
		return false, err
	}
	defer a.Close()
	b, err := os.Open(fileB)
	if err != nil {
		return false, err
	}
	defer b.Close()

	aInfo, err := a.Stat()
	if err != nil {
		return false, err
	}
	bInfo, err := b.Stat()
	if err != nil {
		return false, err
	}

	// check inode equality
	if os.SameFile(aInfo, bInfo) {
		return true, nil
	}
	if aInfo.Size() != bInfo.Size() {
		return false, nil
	}
	return DefaultBufferedReader.ReaderEqual(a, b)
}

// ReaderEqual returns true if both readers returned the same bytes.
func (br *BufferedReader) ReaderEqual(a, b io.Reader) (bool, error) {
	if a == nil || b == nil {
		return false, errors.New("either reader a or b is nil")
	}

	for {
		nA, err := a.Read(br.bufA)
		if err != nil && err != io.EOF {
			return false, err
		}
		nB, err := b.Read(br.bufB)
		if err != nil && err != io.EOF {
			return false, err
		}
		if nA != nB {
			return false, nil
		}
		if nA == 0 && nB == 0 {
			return true, nil
		}
		for idx := 0; idx < nA; idx++ {
			if br.bufA[idx] != br.bufB[idx] {
				return false, nil
			}
		}
	}
}

// ReplaceFile overwrites a file atomically by first creating a temporary file with the given data
// and afterwards replacing the destination with temporary file.
// File permissions of the destination are kept.
// The destination is unchanged if an intermediate error occurs.
func ReplaceFile(data io.Reader, dest string) error {
	destInfo, err := os.Stat(dest)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "toolbelt")
	if err != nil {
		return err
	}
	defer f.Close()

	err = os.Chmod(f.Name(), destInfo.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, data)
	if err != nil {
		return err
	}

	return os.Rename(f.Name(), dest)
}

// CloseAndRemove first closes the given file and secondly deletes it.
func CloseAndRemove(f *os.File) error {
	err := f.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.Name())
}
