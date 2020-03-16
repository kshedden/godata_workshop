package notable

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
)

// A struct holding information about a notable person
type Person struct {

	// The person's name
	PrsLabel string

	// The person's year of birth
	BYear int

	// The person's birth location
	BLocLabel string

	// The latitude of the person's birth location
	BLocLat float64

	// The longitude of the person's birth location
	BLocLong float64

	// The year of the person's birth
	DYear int

	// The location where the person died
	DLocLabel string

	// The latitude of the location where the person died
	DLocLat float64

	// The longitude of the location where the person died
	DLocLong float64

	// The person's gender
	Gender string
}

// A struct holding information about a collection of notable people.
type People struct {

	// The person's name
	PrsLabel []string

	// The person's year of birth
	BYear []int

	// The person's birth location
	BLocLabel []string

	// The latitude of the person's birth location
	BLocLat []float64

	// The longitude of the person's birth location
	BLocLong []float64

	// The year of the person's birth
	DYear []int

	// The location where the person died
	DLocLabel []string

	// The latitude of the location where the person died
	DLocLat []float64

	// The longitude of the location where the person died
	DLocLong []float64

	// The person's gender
	Gender []string
}

// GetCSVWriter returns two Closer's and a csv.Writer for writing
// csv formatted data to the given file.
func GetCSVWriter(fname string) (io.Closer, io.Closer, *csv.Writer) {

	// Open a file for writing
	fid, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	// Write compressed data to the file
	out := gzip.NewWriter(fid)

	// Write CSV-formatted data to the file
	cout := csv.NewWriter(out)

	return fid, out, cout
}

// GetJSONEncoder returns two io.Closer's and a json encoder for writing to
// the given file.
func GetJSONEncoder(fname string) (io.Closer, io.Closer, *json.Encoder) {

	// Open a file for writing
	fid, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	// Write compressed data to the file
	out := gzip.NewWriter(fid)

	// Write the data in json format
	enc := json.NewEncoder(out)

	return fid, out, enc
}

// GetGobDecoder returns a decoder for reading from a gob-encoded data source.
// It also returns two system resources that should be closed
// after the decoder is no longer needed.
func GetGobDecoder(fname string) (io.Closer, io.Closer, *gob.Decoder) {

	// Open a reader for the file
	fid, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	// Decompress the stream on-the-fly
	gid, err := gzip.NewReader(fid)
	if err != nil {
		panic(err)
	}

	// Use this to decode
	dec := gob.NewDecoder(gid)

	return fid, gid, dec
}

// GetGobEncoder returns an encoder for writing gob-encoded data to a stream.
// It also returns two system resources that should be closed
// after the encoder is no longer needed.
func GetGobEncoder(fname string) (io.Closer, io.Closer, *gob.Encoder) {

	// Open a reader for the file
	fid, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	// Decompress the stream on-the-fly
	gid := gzip.NewWriter(fid)

	// Use this to encode
	enc := gob.NewEncoder(gid)

	return fid, gid, enc
}
