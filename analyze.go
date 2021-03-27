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
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"

	// [START imports]
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	// [END imports]
	"internet-thoughts/reddit"
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
	Neutral
)

func (s sentimentClass) String() string {
	return [...]string{"Positive", "Negative", "Mixed", "Neutral"}[s]
}

func main() {
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

	// filename := os.Args[2]
	// if filename == "" {
	// 	usage("Missing text.")
	// }

	// docs := readFile(filename)
	docsAndUpvotes := reddit.GetReddit("bts")
	docs := make([]string, 0)
	for _, doc := range docsAndUpvotes {
		docs = append(docs, doc.Body)
	}

	analyzeMultipleSentiments(ctx, client, docs)

	// switch os.Args[1] {
	// case "entities":
	// 	printResp(analyzeEntities(ctx, client, text))
	// case "sentiment":
	// 	printResp(analyzeSentiment(ctx, client, text))
	// case "syntax":
	// 	printResp(analyzeSyntax(ctx, client, text))
	// // case "entitysentiment":
	// // 	printResp(analyzeEntitySentiment(ctx, betaClient(), text))
	// case "classify":
	// 	printResp(classifyText(ctx, client, text))
	// default:
	// 	usage("Unknown command.")
	// }
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

func analyzeMultipleSentiments(ctx context.Context, client *language.Client, docs []string) {
	var totScore float32

	// var
	sm := map[string]int{}
	lock := sync.RWMutex{}
	c := make(chan float32)
	count := 0
	for _, doc := range docs {
		go func(doc string) {
			score := float32(0.0)
			defer func() { c <- score }()
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
			lock.Unlock()
			// tot_score += score
		}(doc)
	}
	// totScore = float32(0.0)
	for i := 0; i < len(docs); i++ {
		totScore += <-c
	}
	fmt.Printf("%+v", sm)
	fmt.Printf("%.2f", totScore/float32(len(docs)))
}

func getSentiment(sentiment *languagepb.Sentiment) (score float32, class string) {
	score = sentiment.Score
	if score > SENTIMENT_THRESHOLD {
		class = Positive.String()
	} else if score < -SENTIMENT_THRESHOLD {
		class = Negative.String()
	} else if sentiment.Magnitude > NEUTRAL_THRESHOLD {
		class = Mixed.String()
	} else {
		class = Neutral.String()
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
