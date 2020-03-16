// This script demonstrates calculating summary statistics within
// subgroups of a data set.
//
// The data set contains locations and dates of births and deaths for
// notable people.
//
// We calculate the mean year of birth and death in each location.
// Then we sort these alphabetically by the location name, and save
// the results as a CSV file.
//
// We also calculate the frequency distribution of births and deaths
// by location, and calculate the entropy of each distribution.  A
// distribution with more entropy is more diffuse, and it turns out
// that the the birth locations have more entropy than the death
// locations.
//
// See the convert.go script to prepare the data needed by this
// script.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/kshedden/godata_workshop/notable/notable"
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

const (
	// Location of the source data
	dataFile = "fb.gob.gz"
)

// A collection of flags that indicate whether we are working with dates
// of birth or dates of death.
type birthOrDeath int

const (
	birth birthOrDeath = iota
	death
)

func getStats(bd birthOrDeath) float64 {

	// Storage for the current record of data being processed.
	var row []string

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

	// The column where the relevant location information is stored.
	locix := 3
	yearix := 2
	if bd == death {
		locix = 8
		yearix = 7
	}

	// Loop over the data records
	var nc int
	for ; ; nc++ {

		// Read one record
		if err := dec.Decode(&row); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Skip the first row (it contains column labels)
		if nc == 0 {
			continue
		}

		// Convert year to a number
		y, err := strconv.ParseFloat(row[yearix], 64)
		if err != nil {
			continue
		}

		// Update the statistics
		loc := row[locix]
		year[loc] += y
		num[loc]++
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
	fname := "birth_mean_by_year.csv"
	if bd == death {
		fname = "death_mean_by_year.csv"
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
