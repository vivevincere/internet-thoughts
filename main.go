package main

import (
	"encoding/json"
	"fmt"
	"internet-thoughts/sentiment"
	"internet-thoughts/twitter"
	"io/ioutil"
	"log"
	"net/http"

	"internet-thoughts/reddit"

	"github.com/gorilla/mux"
)

type Sentiment_API struct {
	Search_Term string `json:"search_term"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TBD")
}

//Sentiment + emotion
func sentiment_search(w http.ResponseWriter, r *http.Request) {
	var s Sentiment_API
	responseData, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(responseData, &s)
	search_term := s.Search_Term

	// Twitter call
	twitterResponses := twitter.TwitterCall(search_term, "", 500)
	var twitterData []twitter.Data
	for i := 0; i < len(twitterResponses); i++ {
		twitterData = append(twitterData, twitterResponses[i].Data...)
	}

	// print(len(twitterData))

	docs := make([]string, 0)
	for _, t := range twitterData {
		docs = append(docs, t.Text)
	}

	// Reddit call
	redditResponses := reddit.GetReddit(search_term)

	for _, r := range redditResponses {
		docs = append(docs, r.Body)
	}

	// Send to sentiment
	sentiment.CheckSentiment(docs)

	//TODO Word Cloud

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(redditResponses)

}

//TODO related searches

//TODO wordCloud
// returns variable number of words + their values

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/sentiment_search", sentiment_search)
	log.Fatal(http.ListenAndServe(":8080", router))
}
