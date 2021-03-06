package main

import (
	"encoding/json"
	"fmt"
	"internet-thoughts/reddit"
	"internet-thoughts/sentiment"
	"internet-thoughts/twitter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	wordcloud "internet-thoughts/python"
	//"internet-thoughts/reddit"
	"strconv"

	"github.com/gorilla/mux"
)

type Sentiment_API struct {
	Search_Term   string `json:"search_term"`
	Popular_Count int    `json:"popular_count"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TBD")
}

type SentimentResponse struct {
	Sentimeter          sentiment.Sentimeter   `json:"sentimeter"`
	Sentiment_Breakdown []sentiment.Sentiment  `json:"sentiment_breakdown"`
	Emotions            sentiment.Emotions     `json:"emotions"`
	Word_Cloud          []wordcloud.Word_Cloud `json:"word_cloud"`
	Buzz_List           []sentiment.Buzz       `json:"buzz_list"`
}

type TrendingResponse struct {
	Trends []string `json:"trends"`
}

//Sentiment + emotion
func sentiment_search_twitter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Access-Control-Allow-Origin")

	log.Printf("%+v", w.Header())
	if r.Method == "OPTIONS" {
		return
	}

	var s Sentiment_API

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(responseData, &s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	search_term := s.Search_Term
	if search_term == "" {
		http.Error(w, "No search term key or empty search term!", 400)
		return
	}

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
	//redditResponses := reddit.GetReddit(search_term)

	//for _, r := range redditResponses {
	//	docs = append(docs, r.Body)
	//}

	wg := &sync.WaitGroup{}
	// Send to sentiment
	wg.Add(1)
	go func() {
		defer wg.Done()
		ourResponse.Sentimeter, ourResponse.Sentiment_Breakdown = sentiment.CheckSentiment(docs)
	}()

	// Send to emotions
	wg.Add(1)
	go func() {
		defer wg.Done()
		ourResponse.Emotions = sentiment.CheckEmotion(docs)
	}()

	//wordCloud
	fullString := strings.Join(docs, "")
	wg.Add(1)
	go func() {
		defer wg.Done()
		cloudWords := wordcloud.WordCloud(fullString)
		for i := 0; i < len(cloudWords); i++ {
			var tmp wordcloud.Word_Cloud
			tmp.Word = cloudWords[i][0]
			tmp.Count, _ = strconv.Atoi(cloudWords[i][1])
			ourResponse.Word_Cloud = append(ourResponse.Word_Cloud, tmp)
		}
	}()

	wg.Wait()
	//HotBuzz
	ourResponse.Buzz_List = twitter.Twitter_Most(twitterData, s.Popular_Count)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ourResponse)

}

func sentiment_search_reddit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Access-Control-Allow-Origin")

	log.Printf("%+v", w.Header())
	if r.Method == "OPTIONS" {
		return
	}

	var s Sentiment_API

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(responseData, &s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	search_term := s.Search_Term
	if search_term == "" {
		http.Error(w, "No search term key or empty search term!", 400)
		return
	}

	var ourResponse SentimentResponse

	// Twitter call

	// print(len(twitterData))

	docs := make([]string, 0)

	// Reddit call
	redditResponses := reddit.GetReddit(search_term)

	for _, r := range redditResponses {
		docs = append(docs, r.Body)
	}

	wg := &sync.WaitGroup{}
	// Send to sentiment
	wg.Add(1)
	go func() {
		defer wg.Done()
		ourResponse.Sentimeter, ourResponse.Sentiment_Breakdown = sentiment.CheckSentiment(docs)
	}()

	// Send to emotions
	wg.Add(1)
	go func() {
		defer wg.Done()
		ourResponse.Emotions = sentiment.CheckEmotion(docs)
	}()

	//wordCloud
	fullString := strings.Join(docs, "")
	wg.Add(1)
	go func() {
		defer wg.Done()
		cloudWords := wordcloud.WordCloud(fullString)
		for i := 0; i < len(cloudWords); i++ {
			var tmp wordcloud.Word_Cloud
			tmp.Word = cloudWords[i][0]
			tmp.Count, _ = strconv.Atoi(cloudWords[i][1])
			ourResponse.Word_Cloud = append(ourResponse.Word_Cloud, tmp)
		}
	}()

	wg.Wait()
	//HotBuzz
	ourResponse.Buzz_List = reddit.Reddit_Most(redditResponses, s.Popular_Count)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ourResponse)

}

func trending_search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Access-Control-Allow-Origin")

	log.Printf("%+v", w.Header())
	if r.Method == "OPTIONS" {
		return
	}

	var s Sentiment_API
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(responseData, &s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	search_term := s.Search_Term
	if search_term == "" {
		http.Error(w, "No search term key or empty search term!", 400)
		return
	}

	related_terms := sentiment.Related(search_term)

	var ourResponse TrendingResponse
	ourResponse.Trends = related_terms

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ourResponse)

}

//TODO related searches

// returns variable number of words + their values

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/sentiment_search/twitter", sentiment_search_twitter).Methods("POST", "OPTIONS")
	router.HandleFunc("/sentiment_search/reddit", sentiment_search_reddit).Methods("POST", "OPTIONS")
	router.HandleFunc("/trending", trending_search).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", router))
}
