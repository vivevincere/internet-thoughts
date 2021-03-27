package main

import (
	"encoding/json"
	"fmt"
	"internet-thoughts/sentiment"
	"internet-thoughts/twitter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"internet-thoughts/python"
	"internet-thoughts/reddit"
	"strconv"
)

type Sentiment_API struct {
	Search_Term string `json:"search_term"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TBD")
}

type SentimentResponse struct {
	Sentimeter          Sentimeter   `json:"sentimeter"`
	Sentiment_Breakdown []Sentiment  `json:"sentiment_breakdown"`
	Emotions            Emotions     `json:"emotions"`
	Word_Cloud          []Word_Cloud `json:"word_cloud"`
	Buzz_List           []Buzz       `json:"buzz_list"`
}

type Sentimeter struct {
	TotalScore int `json:"total_score"`
	ValidCount int `json:"valid_count"`
}

type Sentiment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Emotions struct {
	ValidCount        int         `json:"valid_count"`
	Emotion_Breakdown []Sentiment `json:"breakdown"`
}

type Word_Cloud struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Buzz struct {
	Text          string `json:"text"`
	Comment_Count int    `json:"comment_count"`
	Retweet_Count int    `json:"retweet_count"`
	Upvote_Count  int    `json:"upvote_count"`
	Url           string `json:"url"`
}

type TrendingResponse struct {
	Trends []string `json:"trends"`
}

//Sentiment + emotion
func sentiment_search(w http.ResponseWriter, r *http.Request) {
	var s Sentiment_API

	responseData, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(responseData, &s)
	search_term := s.Search_Term

	var ourResponse SentimentResponse

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
	//sentiment.CheckSentiment(docs)

	//wordCloud
	fullString := strings.Join(docs, "")

	cloudWords := wordcloud.WordCloud(fullString)
	for i := 0; i < len(cloudWords); i++ {
		var tmp Word_Cloud
		tmp.Word = cloudWords[i][0]
		tmp.Count, _ = strconv.Atoi(cloudWords[i][1])
		ourResponse.Word_Cloud = append(ourResponse.Word_Cloud, tmp)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ourResponse)

}

func trending_search(w http.ResponseWriter, r *http.Request) {
	var s Sentiment_API
	responseData, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(responseData, &s)
	search_term := s.Search_Term

	related_terms := sentiment.Related(search_term)

	var ourResponse TrendingResponse
	ourResponse.Trends = related_terms

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ourResponse)

}

//TODO related searches

//TODO wordCloud
// returns variable number of words + their values

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/sentiment_search", sentiment_search)
	router.HandleFunc("/trending", trending_search)
	log.Fatal(http.ListenAndServe(":8080", router))
}
