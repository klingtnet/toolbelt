package toolbelt

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkReaderEqualMemory(b *testing.B) {
	size := 10 * 1024 * 1024
	r := make([]byte, size)
	n, err := rand.Read(r)
	if err != nil {
		b.Fatal(err)
	}
	if n != size {
		b.Fatalf("expected random source to have %d size but had size %d", size, n)
	}
	bufA := bytes.NewReader(r)
	bufB := bytes.NewReader(r)

	for _, bufSize := range []int64{1, 1024, DefaultReadBufferSize, 1024 * 1024} {
		br := NewBufferedReader(bufSize)

		b.Run(fmt.Sprintf("buffer-size-%d", bufSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				equal, err := br.ReaderEqual(bufA, bufB)
				b.StopTimer()

				if err != nil {
					b.Fatal(err)
				}
				if !equal {
					b.Fatal("must be equal but was not")
				}
				bufA.Reset(r)
				bufB.Reset(r)
			}
		})
	}
}

func BenchmarkReaderEqualDisk(b *testing.B) {
	size := 10 * 1024 * 1024
	r := make([]byte, size)
	n, err := rand.Read(r)
	if err != nil {
		b.Fatal(err)
	}
	if n != size {
		b.Fatalf("expected random source to have %d size but had size %d", size, n)
	}
	fileA, err := ioutil.TempFile("", "testbelt-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer fileA.Close()
	fileB, err := ioutil.TempFile("", "testbelt-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer fileB.Close()

	defer func() {
		os.Remove(fileA.Name())
		os.Remove(fileB.Name())
	}()

	n, err = fileA.Write(r)
	if err != nil {
		b.Fatal(err)
	}
	if n != len(r) {
		b.Fatalf("failed to write full benchmark data, expected %d but only wrote %d bytes", len(r), n)
	}
	n, err = fileB.Write(r)
	if err != nil {
		b.Fatal(err)
	}
	if n != len(r) {
		b.Fatalf("failed to write full benchmark data, expected %d but only wrote %d bytes", len(r), n)
	}

	_, err = fileA.Seek(0, 0)
	if err != nil {
		b.Fatal(err)
	}
	_, err = fileB.Seek(0, 0)
	if err != nil {
		b.Fatal(err)
	}

	for _, bufSize := range []int64{1, 1024, DefaultReadBufferSize, 1024 * 1024} {
		br := NewBufferedReader(bufSize)

		b.Run(fmt.Sprintf("buffer-size-%d", bufSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				equal, err := br.ReaderEqual(fileA, fileB)
				b.StopTimer()

				if err != nil {
					b.Fatal(err)
				}
				if !equal {
					b.Fatal("must be equal but was not")
				}
				_, err = fileA.Seek(0, 0)
				if err != nil {
					b.Fatal(err)
				}
				_, err = fileB.Seek(0, 0)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
