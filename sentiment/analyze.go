// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command analyze performs sentiment, entity, entity sentiment, and syntax analysis
// on a string of text via the Cloud Natural Language API.
package sentiment

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"internet-thoughts/utils"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/groovili/gogtrends"

	// [START imports]
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	// [END imports]
)

var SENTIMENT_THRESHOLD = float32(0.25) // todo: tune this param
var NEUTRAL_THRESHOLD = float32(1.0)

type sentimentClassification struct {
}

type sentimentClass int

const (
	Positive sentimentClass = iota
	Negative
	Mixed
)

func (s sentimentClass) String() string {
	return [...]string{"Positive", "Negative", "Mixed"}[s]
}

func CheckSentiment(docs []string) (Sentimeter, []Sentiment) {
	// if len(os.Args) < 2 {
	// 	usage("Missing command.")
	// }

	// [START init]
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END init]

	sentimeter, sentimentMap := analyzeMultipleSentiments(ctx, client, docs)

	sentiments := make([]Sentiment, 0)
	for sentiClass, score := range sentimentMap {
		sentiments = append(sentiments, Sentiment{sentiClass, score})
	}
	return sentimeter, sentiments
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "usage: analyze [entities|sentiment|syntax|entitysentiment|classify] <text>")
	os.Exit(2)
}

// [START language_entities_text]

func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// [END language_entities_text]

// [START language_sentiment_text]

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_sentiment_text]

// [START language_syntax_text]

func analyzeSyntax(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// [END language_syntax_text]

// [START language_classify_text]

func classifyText(ctx context.Context, client *language.Client, text string) (*languagepb.ClassifyTextResponse, error) {
	return client.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_classify_text]

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}

func analyzeMultipleSentiments(ctx context.Context, client *language.Client, docs []string) (Sentimeter, map[string]int) {
	totScore := float32(0.0)
	totSucceed := 0

	// var
	sm := map[string]int{}
	lock := sync.RWMutex{}
	c := make(chan int)
	count := 0
	for _, doc := range docs {
		go func(doc string) {
			succeed := 0
			defer func() { c <- succeed }()
			count += 1
			senti, err := analyzeSentiment(ctx, client, doc)
			// print(senti.DocumentSentiment.Score)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			score, class := getSentiment(senti.DocumentSentiment)
			lock.Lock()
			sm[class] += 1
			totScore += utils.Scale(score)
			lock.Unlock()
			succeed = 1
			// tot_score += score
		}(doc)
	}
	// totScore = float32(0.0)
	for i := 0; i < len(docs); i++ {
		totSucceed += <-c
	}
	// fmt.Printf("%+v", sm)
	// fmt.Printf("%.2f", totScore/float32(len(docs)))
	log.Printf("%f\n", totScore)
	return Sentimeter{totScore, totSucceed}, sm
}

func getSentiment(sentiment *languagepb.Sentiment) (score float32, class string) {
	score = sentiment.Score
	if score > SENTIMENT_THRESHOLD {
		class = Positive.String()
	} else if score < -SENTIMENT_THRESHOLD {
		class = Negative.String()
	} else {
		class = Mixed.String()
	}
	return
}

func scanBR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil

}

func readFile(name string) []string {
	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Split(scanBR)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func Related(searchTerm string) []string {
	res := make([]string, 0)
	ctx := context.Background()
	explore, err := gogtrends.Explore(ctx,
		&gogtrends.ExploreRequest{
			ComparisonItems: []*gogtrends.ComparisonItem{
				{
					Keyword: searchTerm,
					Geo:     "US",
					Time:    "today 12-m",
				},
			},
			Category: 0, // all categories
			Property: "",
		}, "EN")
	if err != nil {
		log.Fatal(err)
	}

	relT, err := gogtrends.Related(ctx, explore[2], "EN")
	if err != nil {
		log.Fatal(err)
	}

	length := 9
	if len(relT) < 9 {
		length = len(relT)
	}
	for i := 0; i < length; i++ {
		if strings.ToUpper(relT[i].Topic.Title) != strings.ToUpper(searchTerm) {
			res = append(res, relT[i].Topic.Title)
		}
	}
	if len(res) > 8 {
		res = res[:len(res)-1]
	}
	log.Println(res)

	return res
}
