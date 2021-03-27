package main

import(
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/vivevincere/internet-thoughts/twitter"
	"encoding/json"
	"io/ioutil"

)

type Sentiment_API struct{
	Search_Term string `json:"search_term"`
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Println("TBD")
}

func sentiment_search(w http.ResponseWriter, r *http.Request){
	var s Sentiment_API
	responseData, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(responseData, &s)
	search_term := s.Search_Term

	//Twitter call
	twitterResponses := twitter.TwitterCall(search_term,"", 500)
	var twitterData []twitter.Data
	for i := 0; i < len(twitterResponses) ;i++{
		twitterData = append(twitterData, twitterResponses[i].Data...)
	}


	//TODO Reddit call



	//TODO Send to sentiment


	//TODO Word Cloud

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(twitterData)

}


func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/sentiment_search", sentiment_search)
	log.Fatal(http.ListenAndServe(":8080", router))
}