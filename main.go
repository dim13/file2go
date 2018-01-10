package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	var (
		inFile  = flag.String("in", "", "input file (required)")
		outFile = flag.String("out", "", "output file")
		pkgName = flag.String("pkg", "", "package name")
	)
	flag.Parse()

	if *inFile == "" {
		flag.Usage()
		return
	}
	if *outFile == "" {
		*outFile = fileName(*inFile)
	}
	if *pkgName == "" {
		pkg, err := packageName(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		*pkgName = pkg
	}

	r, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	w, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	banner := strings.Join(os.Args, " ")

	if err := generate(w, r, *inFile, *pkgName, banner); err != nil {
		log.Fatal(err)
	}
}

func generate(w io.Writer, r io.Reader, fname, pkg, banner string) error {
	vname := varName(fname)
	fmt.Fprintf(w, "// Code generated by %q; DO NOT EDIT.\n\n", banner)
	fmt.Fprintf(w, "package %s\n\n", pkg)
	fmt.Fprintf(w, "// %s holds content of %s\n", vname, fname)
	fmt.Fprintf(w, "var %s = []byte{", vname)
	scanner := bufio.NewScanner(r)
	scanner.Split(scanBytes)
	for scanner.Scan() {
		fmt.Fprint(w, "\n\t")
		b := scanner.Bytes()
		for i, v := range b {
			fmt.Fprintf(w, "%#02x,", v)
			if i != len(b)-1 {
				fmt.Fprint(w, " ")
			}
		}
	}
	fmt.Fprint(w, "\n}\n")
	return scanner.Err()
}

func scanBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const n = 12
	switch {
	case atEOF && len(data) == 0:
		return 0, nil, nil
	case !atEOF && len(data) < n:
		return 0, nil, nil
	case atEOF:
		return len(data), data, nil
	default:
		return n, data[:n], nil
	}
}

func varName(file string) string {
	return strings.Replace(strings.Title(path.Base(file)), ".", "", -1)
}

func fileName(file string) string {
	return strings.Replace(strings.ToLower(path.Base(file)), ".", "_", -1) + ".go"
}

func packageName(dir string) (string, error) {
	pkg, err := build.Default.ImportDir(path.Dir(dir), 0)
	if err != nil {
		return "", err
	}
	return pkg.Name, nil
}