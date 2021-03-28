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

// Package automl contains samples for Google Cloud AutoML API v1.
package sentiment

// [START automl_language_text_classification_predict]
import (
	"context"
	"fmt"
	"log"
	"sync"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// languageTextClassificationPredict does a prediction for text classification.
func languageTextClassificationPredict(ctx context.Context, client *automl.PredictionClient, content string) (map[string]float32, error) {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "TCN123456789..."
	// content := "text to classify"

	req := &automlpb.PredictRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", "331721144440", "us-central1", "TCN910635321333383168"),
		Payload: &automlpb.ExamplePayload{
			Payload: &automlpb.ExamplePayload_TextSnippet{
				TextSnippet: &automlpb.TextSnippet{
					Content:  content,
					MimeType: "text/plain", // Types: "text/plain", "text/html"
				},
			},
		},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Predict: %v", err)
	}

	emotions := make(map[string]float32, 0)

	for _, payload := range resp.GetPayload() {
		emotions[payload.GetDisplayName()] = payload.GetClassification().GetScore()
	}

	return emotions, nil
}

func predictEmotions(ctx context.Context, client *automl.PredictionClient, docs []string) (int, map[string]float64) {
	defer client.Close()
	allEmotions := make(map[string]float64)
	lock := sync.Mutex{}
	c := make(chan int)
	for _, doc := range docs {
		go func(doc string) {
			succeed := 0
			defer func() { c <- succeed }()
			emotions, err := languageTextClassificationPredict(ctx, client, doc)
			// print(senti.DocumentSentiment.Score)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			// score, class := getSentiment(senti.DocumentSentiment)
			lock.Lock()
			for emotion, score := range emotions {
				allEmotions[emotion] += float64(score)
			}
			lock.Unlock()
			succeed += 1
			// tot_score += score
		}(doc)
	}
	totalSucceed := 0
	for i := 0; i < len(docs); i++ {
		totalSucceed += <-c
	}
	return totalSucceed, allEmotions
}

func CheckEmotion(docs []string) Emotions {
	// if len(os.Args) < 2 {
	// 	usage("Missing command.")
	// }

	// [START init]
	ctx := context.Background()
	client, err := automl.NewPredictionClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	totalScore, breakdown := predictEmotions(ctx, client, docs)
	emotions := make([]Emotion, 0)
	for feeling, score := range breakdown {
		emotion := Emotion{feeling, score}
		emotions = append(emotions, emotion)
	}
	return Emotions{totalScore, emotions}
}

// [END automl_language_text_classification_predict]
