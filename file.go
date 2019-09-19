package toolbelt

import (
	"errors"
	"io"
	"os"
)

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
	return ReaderEqual(a, b, DefaultReadBufferSize)
}

const DefaultReadBufferSize = 4 * 1024

// ReaderEqual returns true if both readers returned the same bytes.
func ReaderEqual(a, b io.Reader, bufSize int64) (bool, error) {
	if a == nil || b == nil {
		return false, errors.New("either reader a or b is nil")
	}
	if bufSize < 1 {
		return false, errors.New("bufSize must be > 0")
	}

	aBuf := make([]byte, bufSize)
	bBuf := make([]byte, bufSize)

	for {
		nA, err := a.Read(aBuf)
		if err != nil && err != io.EOF {
			return false, err
		}
		nB, err := b.Read(bBuf)
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
			if aBuf[idx] != bBuf[idx] {
				return false, nil
			}
		}
	}
}
