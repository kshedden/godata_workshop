// This script collects Tweet data from Twitter using its API, then
// calculates some simple summaries of it.  The Tweets collected
// are from the current (real-time) stream of Tweets.  This script
// also saves some information about each tweet to a file.
//
// Running this script requires the go-twitter and oauth1 libraries.
// You will need to run the following to get them:
//
//   go get github.com/dghubble/go-twitter
//   go get github.com/dghubble/oauth1
//
// To use this script, you must first apply for a Twitter developer
// account at this link:
//
//     https://developer.twitter.com/en/apps.
//
// Your request must be manually reviewed and approved by Twitter,
// which could take a few days.
//
// After you get your account, you will receive a consumer key and
// secret that you can use to access the API.  Once you have these
// keys and tokens, you can run this script using:
//
//   go run streaming.go -consumer-key=## -consumer-secret=## -access-token=## -access-secret=##
//
// replacing ## as appropriate.
//
// Alternatively, see the streaming_public.sh file (a bash shell
// script), where you can enter your key and secret information, and
// then run the bash script (which in turn runs this Go program with
// the required information).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/kshedden/godata_workshop/twitter/twit"
	"gonum.org/v1/gonum/floats"
)

const (
	// We will collect this many tweets.
	nrec int = 2000
)

var (
	// Retain Tweets written in these languages
	langs = []string{"en", "es", "fr"}

	// Use this channel to communicate between the routine the collects
	// the tweets, and the main program.
	rc chan twit.Record
)

// getDemixAndFilter returns a demixer and a stream of tweets.  The demuxer
// is what handles the individual tweets as they are received from the API.
func getDemuxAndFilter(client *twitter.Client) (twitter.SwitchDemux, *twitter.Stream) {

	// Demultiplex stream messages
	demux := twitter.NewSwitchDemux()

	// Format string matching the way that dates are formatted in the Twitter
	// api output.
	s := "Mon Jan 2 15:04:05 -0700 2006"

	// This function gets called every time a tweet is received.
	demux.Tweet = func(tweet *twitter.Tweet) {
		user := tweet.User
		ca, err := time.Parse(s, user.CreatedAt)
		if err != nil {
			panic(err)
		}

		// Create a record containing the information that
		// we want from this tweet, and send it through the
		// channel we have created.
		rc <- twit.Record{ca, tweet.Lang, user.FollowersCount, user.ScreenName}
	}

	fmt.Println("Starting Stream...")

	// Filter, see below for information about how to configure this.
	// https://developer.twitter.com/en/docs/tweets/filter-realtime/guides/basic-stream-parameters
	filterParams := &twitter.StreamFilterParams{

		// Search terms
		Track: []string{"covid", "virus", "coronavirus"},

		// Restrict to this language
		Language: langs,

		// Warn if something seems wrong
		StallWarnings: twitter.Bool(true),
	}

	// Create a stream
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	return demux, stream
}

// Create some summary statistics from the tweet data.
func analyze(results []twit.Record) {

	var created []float64
	var nfollowers []float64

	// Create summary statistics for each language.
	for _, lang := range langs {

		// Reset the slices to empty
		created = created[0:0]
		nfollowers = nfollowers[0:0]
		num := 0

		// Create arrays holding only the results for this language.
		for _, r := range results {
			if r.Language == lang {
				created = append(created, float64(r.CreatedAt.Unix()))
				nfollowers = append(nfollowers, float64(r.FollowersCount))
				num++
			}
		}

		// Sort so we can get the medians.
		sort.Float64Slice(created).Sort()
		sort.Float64Slice(nfollowers).Sort()

		fmt.Printf("Tweets in language %s\n", lang)
		fmt.Printf("Total tweets: %v\n", len(created))

		m := floats.Sum(created) / float64(len(created))
		fmt.Printf("Average account creation date: %v\n", time.Unix(int64(m), 0))

		m = created[len(created)/2]
		fmt.Printf("Median account creation date:  %v\n", time.Unix(int64(m), 0))

		m = floats.Sum(nfollowers) / float64(len(nfollowers))
		fmt.Printf("Average number of followers: %10.2f\n", m)

		m = nfollowers[len(nfollowers)/2]
		fmt.Printf("Median number of followers:  %10.2f\n\n", m)
	}
}

func main() {

	rc = make(chan twit.Record)

	client := twit.GetClient()
	demux, stream := getDemuxAndFilter(client)

	// Receive messages until stopped or stream quits.  This
	// is run in the background using "go".
	go demux.HandleChan(stream.Messages)

	// Save the screen names to this file
	out, err := os.Create("screen_names.json")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	enc := json.NewEncoder(out)

	// Read the desired number of records.
	var results []twit.Record
	for len(results) < nrec {

		// Read a tweet from the channel. This is a blocking operation,
		// control flow sits here and waits for a tweet to arrive.
		r := <-rc

		// Save the results for subsequent analysis
		results = append(results, r)

		if err := enc.Encode(&r); err != nil {
			panic(err)
		}
	}

	analyze(results)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
