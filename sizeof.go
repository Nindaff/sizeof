package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const help = `
Usage:
	size [options...] [data]
	echo "some data" | size
	size -f "some.file"
 	size -u B "some data"

Options:
	-f Filename
	-u Units for output (b, B, Kb, KB, Mb, MB, Gb, GB, Tb, TB)
`

// Map of units
var units = map[string]int64{
	"b":  1,
	"B":  1 << 3,
	"Kb": 1 << 10,
	"KB": 1 << 13,
	"Mb": 1 << 20,
	"MB": 1 << 23,
	"Gb": 1 << 30,
	"GB": 1 << 33,
	"Tb": 1 << 40,
	"TB": 1 << 43,
}

// Flags
var (
	fileFlag  = flag.String("f", "", "")
	unitsFlag = flag.String("u", "B", "")
)

// Parse units, default "B" (bytes).
func parseUnits(u string) int64 {
	unit, ok := units[u]
	if !ok {
		return units["B"]
	}
	return unit
}

// Get byte size of a reader.
func readerSize(r io.Reader) (n int64, err error) {
	buf, err := ioutil.ReadAll(r)
	n = int64(len(bytes.TrimRight(buf, "\n")))
	return
}

// Get byte size of a file via stat.
func fileSize(filename string) (n int64, err error) {
	n = 0
	f, err := os.Open(filename)
	defer f.Close()
	info, err := f.Stat()
	if err == nil {
		n = info.Size()
	}
	return
}

// String size
func stringSize(s string) (int64, error) {
	return int64(len([]byte(s))), nil
}

// Convert bytes into units.
func convert(n int64, u int64) float64 {
	return float64(n*8) / float64(u)
}

// Usage func for flag.
func usage() {
	fmt.Fprintf(os.Stderr, help)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var n int64
	var err error
	u := parseUnits(*unitsFlag)
	switch {
	case *fileFlag != "":
		n, err = fileSize(*fileFlag)
	case flag.NArg() > 0:
		n, err = stringSize(flag.Args()[0])
	default:
		n, err = readerSize(os.Stdin)
	}

	if err != nil {
		fmt.Println("Error")
	} else {
		fmt.Println(convert(n, u))
	}
}
