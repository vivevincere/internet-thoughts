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
