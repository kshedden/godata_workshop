// This script calculates some basic summary statistics from
// one day of Reddit data.  The Reddit data are obtained from
// pushshift, whose top-level domain is:
//
// http://files.pushshift.io/reddit
//
// The Reddit data used here can be obtained from:
//
// http://files.pushshift.io/reddit/comments/daily
//
// To run this script, modify the 'datapath' variable below
// to point to the data file that you have downloaded.
package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

type record struct {
	Author                  string
	Author_created_utc      int
	Body                    string
	Created_utc             int
	Id                      string
	Subreddit               string
	Subreddit_id            string
	Link_id                 string
	Permalink               string
	Subreddit_name_prefixed string
}

const (
	datapath = "/nfs/kshedden/reddit/RC_2019-02-28.gz"
)

var (
	// All distinct users
	allUsers []string

	// Map from thread id to the user id's for that thread
	threadUsers map[string][]string

	// Map from user id to the threads that the user comments on
	userThreads map[string][]string
)

// unique takes a string slice, and returns a slice containing the
// unique values in x, sorted.  The returned slice is a slice of
// the provided slice.
func unique(x []string) []string {

	sort.StringSlice(x).Sort()

	var i, j int
	for j < len(x) {
		for ; j < len(x); j++ {
			if x[j] != x[i] {
				i++
				x[i] = x[j]
				break
			}
		}
	}

	return x[0 : i+1]
}

// Read the raw data and convert it to a structured form for analysis.
func getData() {

	fid, err := os.Open(datapath)
	if err != nil {
		panic(err)
	}
	defer fid.Close()

	gid, err := gzip.NewReader(fid)
	if err != nil {
		panic(err)
	}
	defer gid.Close()

	dec := json.NewDecoder(gid)

	ii := 0
	for {
		if ii%100000 == 0 {
			fmt.Printf("%d\n", ii)
		}
		ii++

		var r record
		if err := dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if r.Author == "[deleted]" {
			continue
		}

		// A unique id for the comment thread
		id := r.Subreddit_id + ":" + r.Link_id

		allUsers = append(allUsers, r.Author)
		threadUsers[id] = append(threadUsers[id], r.Author)
		userThreads[r.Author] = append(userThreads[id], id)
	}
}

// propConnectedUB returns an upper bound on the proportion of pairs of
// people who are connected (i.e. who comment on the same thread).
func propConnectedUB() float64 {

	// Number of users
	n := float64(len(allUsers))

	// Count the total edges between pairs of users (overcounting
	// when users are connected in multile threads).
	te := 0.0
	for _, th := range threadUsers {

		// Size of this thread
		m := float64(len(th))

		// Number of edges in this thread
		te += m * (m - 1)
	}

	// Total possible edges
	tpe := n * (n - 1)

	// An upper bound on the proportion of pairs of users who
	// are connected.
	ub := te / tpe

	return ub
}

// quantilesThreadSize returns estimates of the median and 90th percentiles
// of the thread sizes (number of comments on a thread).
func quantilesThreadSize() (float64, float64) {

	var ts []int

	for _, v := range threadUsers {
		ts = append(ts, len(v))
	}

	sort.IntSlice(ts).Sort()

	return float64(ts[len(ts)/2]), float64(ts[9*len(ts)/10])
}

func main() {

	threadUsers = make(map[string][]string)
	userThreads = make(map[string][]string)

	getData()

	allUsers = unique(allUsers)

	for k, v := range threadUsers {
		threadUsers[k] = unique(v)
	}

	for k, v := range userThreads {
		userThreads[k] = unique(v)
	}

	// Remove single-post threads
	xthreadusers := make(map[string][]string)
	for k, v := range threadUsers {
		if len(v) > 1 {
			xthreadusers[k] = v
		}
	}
	threadUsers = xthreadusers

	fmt.Printf("Number of distinct users:        %d\n", len(allUsers))
	fmt.Printf("Number of threads:               %d\n", len(threadUsers))

	q1, q2 := quantilesThreadSize()
	fmt.Printf("Median thread size:              %.2f\n", q1)
	fmt.Printf("90th percentile of thread size:  %.2f\n", q2)

	ub := propConnectedUB()
	fmt.Printf("Estimated proportion of connected users:  %v\n", ub)
}
