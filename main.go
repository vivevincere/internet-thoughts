package main

import(
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"twitter"

)

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Println("TBD")
}

func sentiment_search(w http.ResponseWriter, r *http.Request){



}


func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8080", router))
}