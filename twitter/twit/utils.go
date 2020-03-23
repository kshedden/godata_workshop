package twit

import (
	"flag"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// record contains the information from one tweet that we will retain.
type Record struct {

	// The date when the account was created
	CreatedAt time.Time

	// The language in which the tweet was written (using 2-letter ISO
	// language codes)
	Language string

	// The number of followers of the person making this Tweet.
	FollowersCount int

	// The name (account handle) of the person making this Tweet.
	ScreenName string
}

// GetClient returns a client that can be used to create a stream of twitter
// data from the Twitter API.
func GetClient() *twitter.Client {

	// Get the user's account information from command-line flags
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flag.String("access-token", "", "Twitter Access Token")
	accessSecret := flag.String("access-secret", "", "Twitter Access Secret")
	flag.Parse()

	// If any of the flags are missing, exit with an error message.
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	// Use these to authorize to the Twitter server
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	return client
}
