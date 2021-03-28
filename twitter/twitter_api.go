package twitter

import (
	"io/ioutil"
	//"log"
	"encoding/json"
	"fmt"
	"internet-thoughts/sentiment"
	"net/http"
	"sort"
)

type TwitterMain struct {
	Meta Meta   `json:"meta"`
	Data []Data `json:"data"`
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

	var retList []TwitterMain
	client := &http.Client{}
	nextToken := ""
	for numberOfTweets > 0 {
		curNum := 100
		if numberOfTweets < 100 {
			curNum = numberOfTweets
		}

		url := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%s&tweet.fields=geo,public_metrics&place.fields=country,name&max_results=%d", searchTerm, curNum)
		if nextToken != "" {
			url += "&next_token=" + nextToken
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAAJDDNwEAAAAAvR58LgY1daTewd2C4htoRLLBxO4%3D2gqj5PSlNfMOnrSVwX8La6FNsBj6Rc5vJRbH81tWqkoLS4Lcej")

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

func Twitter_Most(data_arr []Data, n int) []sentiment.Buzz {
	var x []sentiment.Buzz
	sort.Slice(data_arr, func(i, j int) bool {
		return data_arr[i].Public_Metrics.Retweet_Count > data_arr[j].Public_Metrics.Retweet_Count
	})
	for i := 0; i < n; i++ {
		var tmp sentiment.Buzz
		tmp.Id = data_arr[i].Id
		x = append(x, tmp)

	}
	return x
}
