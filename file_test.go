package toolbelt

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func createFile(t *testing.T, content string) string {
	f, err := ioutil.TempFile("", "toolbelt-test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	n, err := f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	if n != len(content) {
		t.Fatalf("failed to write full content into test file %q", f.Name())
	}
	return f.Name()
}

func createRandomFileWithSize(t *testing.T, size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, size)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != size {
		t.Fatalf("expected to generate %d random bytes but was %d", size, n)
	}
	return createFile(t, string(buf))
}

func TestFilesEqual(t *testing.T) {
	largeRandomFile := createRandomFileWithSize(t, DefaultReadBufferSize*2)
	tCases := map[string]struct {
		fileA, fileB string
		equal        bool
		errMsg       string
	}{
		"equal": {
			createFile(t, "equal"),
			createFile(t, "equal"),
			true,
			"",
		},
		"equal-larger-than-buffer": {
			largeRandomFile,
			largeRandomFile,
			true,
			"",
		},
		"unequal": {
			createFile(t, "something"),
			createFile(t, "different"),
			false,
			"",
		},
		"unequal-random-content": {
			createRandomFileWithSize(t, DefaultReadBufferSize),
			createRandomFileWithSize(t, DefaultReadBufferSize),
			false,
			"",
		},
		"does-not-exist": {
			"file-does-not-exist",
			"file-does-not-exist",
			false,
			"no such file or directory",
		},
		"start-differs": {
			createFile(t, "has some prefix but this identical"),
			createFile(t, "this is identical"),
			false,
			"",
		},
		"end-differs": {
			createFile(t, "this is identical"),
			createFile(t, "this is identical except the end"),
			false,
			"",
		},
	}
	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				os.Remove(tCase.fileA)
				os.Remove(tCase.fileB)
			}()
			equal, err := FilesEqual(tCase.fileA, tCase.fileB)
			if err != nil {
				if tCase.errMsg == "" {
					t.Fatalf("expected to succeed but failed: %s", err)
				}
				if !strings.Contains(err.Error(), tCase.errMsg) {
					t.Fatalf("expected error message to contain %q but was %q", tCase.errMsg, err.Error())
				}
			}
			if equal != tCase.equal {
				t.Fatalf("expected equal %t but was %t", tCase.equal, equal)
			}
		})
	}
}

func TestReplaceFile(t *testing.T) {
	dest, err := ioutil.TempFile("", "toolbelt-test")
	if err != nil {
		t.Fatal(err)
	}
	destFilepath := dest.Name()
	defer os.Remove(destFilepath)
	_, err = io.Copy(dest, bytes.NewReader([]byte("original content")))
	if err != nil {
		t.Fatal(err)
	}
	destMode := os.FileMode(0600)
	err = dest.Chmod(destMode)
	if err != nil {
		t.Fatal(err)
	}
	err = dest.Close()
	if err != nil {
		t.Fatal(err)
	}
	newData, err := os.Open("README.md")
	if err != nil {
		t.Fatal(err)
	}
	err = ReplaceFile(newData, destFilepath)
	if err != nil {
		t.Fatal(err)
	}
	equal, err := FilesEqual("README.md", destFilepath)
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Fatal("contents differ")
	}
	destInfo, err := os.Stat(destFilepath)
	if err != nil {
		t.Fatal(err)
	}
	if destInfo.Mode() != destMode {
		t.Fatalf("expected file mode to be %v but was %v", destMode, destInfo.Mode())
	}
}

func TestReplaceFileIfDifferent(t *testing.T) {
	// cases:
	// files are different -> change
	// files are equal -> no change

	tCases := map[string]struct {
		src, dest    string
		shouldChange bool
	}{
		"different": {
			createFile(t, "this is new"),
			createFile(t, "this is old"),
			true,
		},
		"different-larger-than-buffer": {
			createRandomFileWithSize(t, DefaultReadBufferSize*2),
			createRandomFileWithSize(t, DefaultReadBufferSize*2),
			true,
		},
		"equal": {
			createFile(t, "equal content"),
			createFile(t, "equal content"),
			false,
		},
	}
	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				os.Remove(tCase.src)
				os.Remove(tCase.dest)
			}()
			src, err := os.Open(tCase.src)
			if err != nil {
				t.Fatal(err)
			}
			changed, err := ReplaceFileIfDifferent(src, tCase.dest)
			if err != nil {
				t.Fatal(err)
			}
			if tCase.shouldChange {
				if !changed {
					t.Fatal("expected destination to change but did not")
				}
				equal, err := FilesEqual(tCase.src, tCase.dest)
				if err != nil {
					t.Fatal(err)
				}
				if !equal {
					t.Fatal("destination was updated but is not equal to source")
				}
			}
			if !tCase.shouldChange {
				if changed {
					t.Fatal("expected not to change but changed")
				}
			}
		})
	}
}
