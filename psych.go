package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	APIURL         = "https://api.twitter.com/1/statuses/user_timeline.json?count=200&screen_name="
	fetchNewQuotes = flag.Bool("r", false, "Fetch new quotes")
)

type Quote struct {
	Text string
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	quotes := loadQuotes(*fetchNewQuotes)
	rnum := rand.Intn(len(quotes) - 1)
	quote := quotes[rnum]
	fmt.Printf("%s\n", quote.Text)
}

func getTweets(name string) []byte {
	client := &http.Client{}
	res, err := client.Get(APIURL + name)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Error fetching tweets: %v", err)
	}
	json, rerr := ioutil.ReadAll(res.Body)
	if rerr != nil {
		log.Printf("Error reading body of res: %v", rerr)
	}
	return json
}

func getQuotes(js []byte) []Quote {
	quotes := []Quote{}
	err := json.Unmarshal(js, &quotes)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
	}
	return quotes
}

func loadQuotes(fetchNewQuotes bool) []Quote {
	content, err := ioutil.ReadFile("/tmp/quotes.json")
	if err != nil || fetchNewQuotes {
		tweets := getTweets("QuotablePsych")
		quotes := getQuotes(tweets)
		saveQuotes(quotes)
		return quotes
	}
	quotes := getQuotes(content)
	return quotes
}

func saveQuotes(quotes []Quote) {
	b, err := json.Marshal(quotes)
	if err != nil {
		log.Printf("Error encoding json: %v", err)
	}
	ioutil.WriteFile("/tmp/quotes.json", b, 0666)
}
