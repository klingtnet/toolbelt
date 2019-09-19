package toolbelt

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
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

func TestFilesEqual(t *testing.T) {
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
		"unequal": {
			createFile(t, "something"),
			createFile(t, "different"),
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
			createFile(t, "this is indentical except the end"),
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
