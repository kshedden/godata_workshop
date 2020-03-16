// This script calculates some simple statistical summaries using data on
// the births and deaths of notable people.
//
// See the comments at the top of convert.go for more information about
// the data.
//
// This script takes the geographical latitude and longitude coordinates
// for each person's birth and death location, and calculates the distance
// in km between these two points.  It then prints a sequence of quantiles
// of the distribution of these distances to stdout.
package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/kshedden/godata_workshop/notable/notable"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

var (
	// Raw data, map from person's name to birth and death
	// locations
	rdata map[string]*loct
)

// Location information for one person
type loct struct {
	BirthLoc       orb.Point // Birth location
	DeathLoc       orb.Point // Death location
	BirthDeathDist float64   // Distance from birth location to death location
}

// Make a map from column name to column index
func makeColix(head []string) map[string]int {
	colix := make(map[string]int)
	for i, v := range head {
		colix[v] = i
	}
	return colix
}

// readData reads the raw data file and creates a map from the
// person's name to an instance of the rec_t struct containing birth
// and death location information.  The BDDist field is not filled in
// here.
func readData(first, last int) {

	f, g, dec := notable.GetGobDecoder("fb_struct_cols.gob.gz")
	defer f.Close()
	defer g.Close()

	var people notable.People
	if err := dec.Decode(&people); err != nil {
		panic(err)
	}

	// Populate rdata
	rdata = make(map[string]*loct)
	for i := 0; i < len(people.PrsLabel); i++ {

		if people.BYear[i] < first || people.BYear[i] > last {
			continue
		}

		// Convert the coordinates to Point objects
		birthloc := orb.Point{people.BLocLat[i], people.BLocLong[i]}
		deathloc := orb.Point{people.DLocLat[i], people.DLocLong[i]}

		rdata[people.PrsLabel[i]] = &loct{BirthLoc: birthloc, DeathLoc: deathloc}
	}
}

// getDistances calculates the distance in km between birth and death
// locations for each person.
func getDistances() {
	for _, v := range rdata {
		di := geo.Distance(v.BirthLoc, v.DeathLoc)
		v.BirthDeathDist = di / 1000 // Convert from meters to km
	}
}

// sumaries prints some statistical summaries of the data.  The
// summaries are quantiles of the distribution of distances between
// the birth and death location of a person.
func summaries() {

	// Extract the distances into an array
	dx := make([]float64, len(rdata))
	i := 0
	for _, v := range rdata {
		dx[i] = v.BirthDeathDist
		i++
	}

	// Sort the distances
	sort.Float64Slice(dx).Sort()

	// The quantiles to display
	qtl := []float64{0.1, 0.25, 0.5, 0.75, 0.9}

	// Calculate and display the quantiles
	fmt.Printf("Probability   Quantile\n")
	for _, q := range qtl {
		pos := int(q * float64(len(dx)-1))
		fmt.Printf("%5.2f       %9.2f\n", q, dx[pos])
	}
}

func main() {

	var first, last int
	flag.IntVar(&first, "first", -100000, "First year of data selection")
	flag.IntVar(&last, "last", 100000, "Last year of data selection")
	flag.Parse()

	readData(first, last)
	getDistances()
	summaries()
}
