package twitter

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"internet-thoughts/sentiment"
	"net/http"
	"sort"
	"net/url"
)

type TwitterMain struct {
	Meta Meta   `json:"meta"`
	Data []Data `json:"data"`
}

type TwitterCreds struct {
	Auth string `json:"Authorization"`
}

type Meta struct {
	Next_token string `json:"next_token"`
}
type Data struct {
	Id             string         `json:"id"`
	Text           string         `json:"text"`
	Public_Metrics Public_Metrics `json:"public_metrics"`
}

type Public_Metrics struct {
	Retweet_Count int `json:"retweet_count"`
	Reply_Count   int `json:"reply_count"`
	Like_Count    int `json:"like_count"`
	Quote_Count   int `json:"quote_count"`
}

func TwitterCall(searchTerm string, location string, numberOfTweets int) []TwitterMain {
	creds, _ := getTwitterCreds()

	var retList []TwitterMain
	client := &http.Client{}
	nextToken := ""
	for numberOfTweets > 0 {
		curNum := 100
		if numberOfTweets < 100 {
			curNum = numberOfTweets
		}
		
		searchTerm = url.QueryEscape(searchTerm)
		searchUrl := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s&tweet.fields=geo,public_metrics&place.fields=country,name&max_results=%d", searchTerm, curNum)
		if nextToken != "" {
			searchUrl += "&next_token=" + nextToken
		}
		req, _ := http.NewRequest("GET", searchUrl, nil)

		req.Header.Set("Authorization", creds.Auth)

		resp, _ := client.Do(req)
		responseData, _ := ioutil.ReadAll(resp.Body)
		var twitterResponse TwitterMain
		json.Unmarshal(responseData, &twitterResponse)
		nextToken = twitterResponse.Meta.Next_token
		numberOfTweets -= curNum
		retList = append(retList, twitterResponse)

	}
	return retList
}

func getTwitterCreds() (*TwitterCreds, error) {
	data, err := ioutil.ReadFile("twitter.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var creds TwitterCreds
	err = json.Unmarshal(data, &creds)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &creds, nil
}

func Twitter_Most(data_arr []Data, n int) []sentiment.Buzz {
	var x []sentiment.Buzz
	sort.Slice(data_arr, func(i, j int) bool {
		return data_arr[i].Public_Metrics.Retweet_Count > data_arr[j].Public_Metrics.Retweet_Count
	})

	keys := make(map[string]bool)

	for i := 0; i < n; i++ {
		if _, value := keys[data_arr[i].Id]; !value{

		keys[data_arr[i].Id] = true
		var tmp sentiment.Buzz
		tmp.Id = data_arr[i].Id
		x = append(x, tmp)
	} else{
		n+= 1
	}

	}
	return x
}
