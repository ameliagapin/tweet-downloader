# tweet-downloader

# Usage

Download project:
```
go get github.com/entirelyamelia/tweet-downloader
```

Set your credentials in `main.go`:
```

const (
	ConsumerKey    string = "<your_consumer_key>"
	ConsumerSecret string = "<your_consumer_secret>"
	AccessToken    string = "<your_access_token>"
	AccessSecret   string = "<your_access_secret"
)
```

Install:
```
go install
```

Execute:
```
tweet-downloader <twitter_handle>
```

To clean user mentions, hashtags, and links, use the optional argument `-clean=true`

The script will download up to 3200 of your last tweets and write them to a file called `tweets_<twitter_handle>.txt` in your downloads folder. Native retweets and replies will be excluded.
