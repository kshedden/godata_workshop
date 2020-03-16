// This script creates a version of the Freebase data in which the records are stored
// as structs.  This is a counterpart to the "convert.go" script in which the
// data are stored as arrays of strings.  By storing the data as structs,
// each data field is stored as a typed value, which makes it easier to use the data
// without further conversion.
package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kshedden/godata_workshop/notable/notable"
)

const (
	// Read the data from a file in which each data record is stored
	// as a slice of strings.
	dataFile = "fb.gob.gz"

	// Write the data as a stream of struct values to this location.
	outFile = "fb_struct.gob.gz"
)

func convert() {

	var row []string

	f1, g1, dec := notable.GetGobDecoder(dataFile)
	f2, g2, enc := notable.GetGobEncoder(outFile)

	// It would be a resource leak not to close these
	defer f1.Close()
	defer g1.Close()
	defer f2.Close()
	defer g2.Close()

	// Loop over the data records
	var nc int
	for ; ; nc++ {

		// Get one row of data
		if err := dec.Decode(&row); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Skip the header row
		if nc == 0 {
			continue
		}

		bloclat, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			panic(err)
		}

		bloclon, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			panic(err)
		}

		dloclat, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			panic(err)
		}

		dloclon, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			panic(err)
		}

		byear, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}

		dyear, err := strconv.ParseInt(row[7], 10, 64)
		if err != nil {
			panic(err)
		}

		// Create a struct holding the data
		person := notable.Person{
			PrsLabel:  row[0],
			BYear:     int(byear),
			BLocLabel: row[3],
			BLocLat:   bloclat,
			BLocLong:  bloclon,
			DYear:     int(dyear),
			DLocLabel: row[8],
			DLocLat:   dloclat,
			DLocLong:  dloclon,
			Gender:    row[12],
		}

		if err := enc.Encode(&person); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Processed %d records\n", nc)
}

func main() {

	convert()
}
