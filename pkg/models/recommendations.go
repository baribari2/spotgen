package models

type Recommendations struct {
	Seeds  []Seed      `json:"seeds"`
	Tracks []FullTrack `json:"tracks"`
}

type Seed struct {
	AfterFilteringSize int    `json:"afterFilteringSize"`
	AfterRelinkingSize int    `json:"afterRelinkingSize"`
	Endpoint           string `json:"href"`
	Id                 string `json:"id"`
	InitialPoolSize    int    `json:"initialPoolSize"`
	Type               string `json:"type"`
}

type ArtistResponse struct {
	Base
	Artists []Artist `json:"items"`
}
