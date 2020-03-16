// This is equivalent to location_stats.go, using a struct-encoded
// version of the data.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"

	"github.com/kshedden/godata_workshop/notable/notable"
)

const (
	// The data to analyze
	dataFile = "fb_struct.gob.gz"
)

// entropy returns the entropy of the frequency distribution of the
// values of the map num.
func entropy(num map[string]int) float64 {

	// Get the total
	tot := 0
	for _, v := range num {
		tot += v
	}

	// Accumulate the entropy
	e := float64(0)
	for _, v := range num {

		// The proportion of events at this location
		p := float64(v) / float64(tot)

		// Update the entropy
		e -= p * math.Log(p)
	}

	return e
}

// A collection of flags that indicate whether we are working with dates
// of birth or dates of death.
type birthOrDeath int

const (
	birth birthOrDeath = iota
	death
)

// getStats calculates summary statistics for either the birth
// locations of the death locations.
func getStats(bd birthOrDeath) float64 {

	fid, gid, dec := notable.GetGobDecoder(dataFile)

	// It would be a resource leak not to close these
	defer fid.Close()
	defer gid.Close()

	// Accumulate the sum of all birth years at each birth
	// location (later will be scaled to obtain the mean).
	year := make(map[string]float64)

	// Accumulate the number of people born at each birth
	// location.
	num := make(map[string]int)

	// Loop over the data records
	var nc int
	for ; ; nc++ {

		var person notable.Person
		if err := dec.Decode(&person); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Update the statistics
		switch bd {
		case birth:
			year[person.BLocLabel] += float64(person.BYear)
			num[person.BLocLabel]++
		case death:
			year[person.DLocLabel] += float64(person.DYear)
			num[person.DLocLabel]++
		default:
			panic("!!")
		}
	}

	// Divide the total by the number of values to get the mean
	for k := range year {
		year[k] /= float64(num[k])
	}

	// Extract the keys and sort them.
	var a []string
	for k := range year {
		a = append(a, k)
	}
	sort.StringSlice(a).Sort()

	// Save the results as a CSV file, sorted alphabetically by
	// location.
	fname := "birth_mean_by_year_structs.csv"
	if bd == death {
		fname = "death_mean_by_year_structs.csv"
	}

	// Create a file to hold the results
	out, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Prepare to write CSV format data
	cout := csv.NewWriter(out)
	defer cout.Flush()

	// Write one row of data
	var crow []string
	for _, k := range a {
		crow = crow[0:0]
		crow = append(crow, fmt.Sprintf("%s", k))
		crow = append(crow, fmt.Sprintf("%.0f", year[k]))
		crow = append(crow, fmt.Sprintf("%d", num[k]))
		if err := cout.Write(crow); err != nil {
			panic(err)
		}
	}

	return entropy(num)
}

func main() {

	e := getStats(birth)
	fmt.Printf("Birth entropy: %f\n", e)

	e = getStats(death)
	fmt.Printf("Death entropy: %f\n", e)
}
