package sentiment

type Sentiment struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Emotion struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Sentimeter struct {
	TotalScore float32 `json:"total_score"`
	ValidCount int     `json:"valid_count"`
}

type Emotions struct {
	ValidCount        int       `json:"valid_count"`
	Emotion_Breakdown []Emotion `json:"breakdown"`
}

type Buzz struct {
	Text          string `json:"text"`
	Comment_Count int    `json:"comment_count"`
	Retweet_Count int    `json:"retweet_count"`
	Upvote_Count  int    `json:"upvote_count"`
	Url           string `json:"url"`
}
