package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	golden, err := ioutil.ReadFile("testdata/empty_gz.go")
	if err != nil {
		t.Fatal(err)
	}
	bin, err := os.Open("testdata/empty.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer bin.Close()
	buf := new(bytes.Buffer)
	err = generate(buf, bin, "empty.gz", "testdata", "file2go -in empty.gz -pkg testdata")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), golden) {
		t.Error("not equal")
	}
}
