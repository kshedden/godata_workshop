// Create a column-oriented version of the Freebase data.

package main

import (
	"io"

	"github.com/kshedden/godata_workshop/notable/notable"
)

func convert() {

	f1, g1, dec := notable.GetGobDecoder("fb_struct.gob.gz")

	// Close these to avoid a resource leak
	defer f1.Close()
	defer g1.Close()

	var people notable.People

	for {
		var person notable.Person
		if err := dec.Decode(&person); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Append all the attributes of the current person to
		// people.
		people.PrsLabel = append(people.PrsLabel, person.PrsLabel)
		people.BYear = append(people.BYear, person.BYear)
		people.BLocLabel = append(people.BLocLabel, person.BLocLabel)
		people.BLocLat = append(people.BLocLat, person.BLocLat)
		people.BLocLong = append(people.BLocLong, person.BLocLong)
		people.DYear = append(people.DYear, person.DYear)
		people.DLocLabel = append(people.DLocLabel, person.DLocLabel)
		people.DLocLat = append(people.DLocLat, person.DLocLat)
		people.DLocLong = append(people.DLocLong, person.DLocLong)
		people.Gender = append(people.Gender, person.Gender)
	}

	f2, g2, enc := notable.GetGobEncoder("fb_struct_cols.gob.gz")

	// Close these, or all the data may not be written to the file.
	defer f2.Close()
	defer g2.Close()

	if err := enc.Encode(&people); err != nil {
		panic(err)
	}
}

func main() {

	convert()
}
