package reddit

import (
	"encoding/json"
	"fmt"
	"internet-thoughts/sentiment"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var LIMIT = 500

type Auth struct {
	AccessToken string `json:"access_token"`
}

type RedditData struct {
	Body      string
	Id        string
	Ups       int
	Author    string
	Permalink string
	Replies   interface{} `json:"replies"`
	Children  []RedditEntity
}

type RedditReplies struct {
	Count int
}

func (rr *RedditReplies) UnmarshalJson(b []byte) error {
	var reply interface{}
	json.Unmarshal(b, &reply)

	rep, ok := reply.(map[string]interface{})
	if !ok {
		return nil
	}
	data, ok := rep["data"]
	if !ok {
		return nil
	}
	dat, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	children, ok := dat["children"]
	if !ok {
		return nil
	}

	arr, ok := children.([]interface{})
	if !ok {
		return nil
	}

	child, ok := arr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	data, ok = child["data"]
	if !ok {
		log.Println("3")
		return nil
	}

	dc, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	float, ok := dc["count"].(float64)
	if !ok {
		return nil
	}
	rr.Count = int(float)
	return nil
}

type RedditEntity struct {
	Kind string
	Data RedditData
}

type Comment struct {
	Ups          int
	Body         string
	RepliesCount int
	UserHandle   string
	CommentLink  string
}

type RedditCreds struct {
	Username     string
	Password     string
	UserAgent    string `json:"user_agent"`
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
}

func main() {
	comms := GetReddit("bts")
	print(comms)
}

// TODO: tune the number of submissions searched and number of comments searched for each subreddit
func GetReddit(query string) []Comment {
	creds, err := getRedditCreds()
	if err != nil {
		return nil
	}
	auth := getAuth(creds)
	augmentReq := setUpCreds(auth, creds)
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://oauth.reddit.com/search?q=%s", query), nil)
	augmentReq(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error : %s", err)
	}

	submissions := getSubmissions(resp)
	lock := sync.Mutex{}
	allComments := make([]Comment, 0)
	c := make(chan int)

	for _, submission := range submissions {
		go func(submission string) {
			client_new := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("https://oauth.reddit.com/comments/%s", submission), nil)
			if err != nil {
				return
			}
			augmentReq(req)
			resp, err := client_new.Do(req)
			if err != nil {
				return
			}
			comments := getCommentsFromSubmission(resp)
			lock.Lock()
			allComments = append(allComments, comments...)
			lock.Unlock()
			c <- len(comments)
		}(submission)
	}
	timer1 := time.NewTimer(7 * time.Second)

	total := 0
loop:
	for {
		select {
		case num := <-c:
			total += num
			if total > LIMIT {
				break loop
			}
		case <-timer1.C:
			print("times up")
			break loop
		}
	}

	print(len(allComments))
	return allComments
}

func getRedditCreds() (*RedditCreds, error) {
	data, err := ioutil.ReadFile("reddit.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var creds RedditCreds
	err = json.Unmarshal(data, &creds)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &creds, nil
}

func getAuth(cred *RedditCreds) *Auth {
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", cred.Username)
	data.Set("password", cred.Password)

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	req.SetBasicAuth(cred.ClientId, cred.ClientSecret)
	// req.SetBasicAuth(client_id, client_secret)
	req.Header.Set("User-Agent", cred.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error : %s", err)
	}
	defer resp.Body.Close()
	var auth Auth
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v\n", string(bodyBytes))
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bodyBytes, &auth)
	} else {
		log.Fatalf("Bad status code %d received\n", resp.StatusCode)
	}
	return &auth
}

func setUpCreds(auth *Auth, cred *RedditCreds) func(*http.Request) {
	augmentRequest := func(req *http.Request) {
		req.Header.Set("User-Agent", cred.UserAgent)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", auth.AccessToken))
	}
	return augmentRequest
}

func getSubmissions(resp *http.Response) []string {
	defer resp.Body.Close()
	submissions := make([]string, 0)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	entity := RedditEntity{}
	err = json.Unmarshal(bodyBytes, &entity)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if entity.Kind != "Listing" {
		log.Fatalf("Listing expected for submissions")
	}
	for _, child := range entity.Data.Children {
		if child.Kind == "t3" { // t3 represents submissions
			submissions = append(submissions, child.Data.Id)
		}
	}
	return submissions
}

func getCommentsFromSubmission(resp *http.Response) []Comment {
	defer resp.Body.Close()
	comments := make([]Comment, 0)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	entities := []RedditEntity{}
	err = json.Unmarshal(bodyBytes, &entities)
	if len(entities) != 2 {
		log.Fatal("Error: length of response to getting comments is not 2 as expected\n")
	}

	for _, child := range entities[1].Data.Children {
		if child.Kind != "t1" { // t1 represents comments
			continue
		}
		// count := 0
		replyBytes, _ := json.Marshal(child.Data.Replies)
		reply := &RedditReplies{}
		err = reply.UnmarshalJson(replyBytes)
		if err != nil {
			fmt.Println(err)
		}
		count := reply.Count
		link := fmt.Sprintf("https://reddit.com%s", child.Data.Permalink)

		comments = append(comments, Comment{child.Data.Ups, child.Data.Body, count, child.Data.Author, link})
	}
	return comments
}

func Reddit_Most(comments []Comment, n int) []sentiment.Buzz {
	var x []sentiment.Buzz
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Ups > comments[j].Ups
	})
	for i := 0; i < n; i++ {
		var tmp sentiment.Buzz
		// tmp.Text = comments[i].Body
		// tmp.Comment_Count = comments[i].RepliesCount
		// tmp.Retweet_Count = comments[i].Public_Metrics.Retweet_Count
		// tmp.Upvote_Count = comments[i].Ups
		tmp.Id = comments[i].CommentLink
		// tmp.User = comments[i].UserHandle
		x = append(x, tmp)

	}
	return x
}
