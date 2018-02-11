package main

import (
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"os"
	"os/user"
	"regexp"
	"strings"
)

// Credentials
const (
	ConsumerKey    string = ""
	ConsumerSecret string = ""
	AccessToken    string = ""
	AccessSecret   string = ""
)

const (
	// File is the location in the user's home directory and filename to write to
	File string = "Downloads/tweets_%s.txt"
	// MaxTweets is the maximum number of tweets to download
	MaxTweets int = 32000
)

var (
	reEllipse = regexp.MustCompile(`\.\.\.`)
	reLink    = regexp.MustCompile(`https?\:\/\/\S*`)
	reAt      = regexp.MustCompile(`\@[a-zA-Z0-9_]*`)
)

func main() {
	var doClean bool
	var username string

	set := flag.NewFlagSet("", flag.ExitOnError)
	set.BoolVar(&doClean, "clean", false, "Clean tweets of mentions, hashtags, and links.")

	if len(os.Args) == 3 {
		set.Parse(os.Args[2:])
		username = os.Args[1]
	} else if len(os.Args) == 2 {
		username = os.Args[1]
	} else {
		panic("No username provided")
	}

	if username == "" || strings.HasPrefix(username, "-clean") || strings.HasPrefix(username, "--clean") {
		panic("No username provided")
	}

	api := getAPI()
	tweets := getTweets(api, username)
	if tweets != nil {
		file := fmt.Sprintf(File, username)

		writeTweets(tweets, file, doClean)
	}
}

func getAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)
	api := anaconda.NewTwitterApi(AccessToken, AccessSecret)
	return api
}

func getTweets(api *anaconda.TwitterApi, username string) *[]anaconda.Tweet {
	v := url.Values{}
	v.Set("screen_name", username)
	v.Set("include_rts", "false")
	v.Set("exclude_replies", "true")
	v.Set("count", "200")
	v.Set("tweet_mode", "extended")

	counter := 0
	var lastTweet *anaconda.Tweet

	var tweets []anaconda.Tweet

	for counter < MaxTweets {
		batch, err := api.GetUserTimeline(v)

		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		tweets = append(tweets, batch...)

		// Set up next iteration
		counter += 200
		lastTweet = &tweets[len(tweets)-1]
		v.Set("max_id", lastTweet.IdStr)
	}

	return &tweets
}

func writeTweets(tweets *[]anaconda.Tweet, filename string, doClean bool) {
	myself, err := user.Current()

	if err != nil {
		panic(err)
	}

	homedir := myself.HomeDir
	location := homedir + "/" + filename

	f, err := os.Create(location)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, tweet := range *tweets {
		text := tweet.FullText
		if doClean {
			text = clean(text)
		}

		_, err := f.WriteString(text + "\n")
		if err != nil {
			panic(err)
		}
	}

	f.Sync()
}

func clean(text string) string {
	dirty := reEllipse.ReplaceAllString(text, ` `)
	dirty = reLink.ReplaceAllString(dirty, ` `)
	clean := reAt.ReplaceAllString(dirty, ` `)

	return clean
}
