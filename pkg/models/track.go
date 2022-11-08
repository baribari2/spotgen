package models

type FullTrack struct {
	Artists          []Artist          `json:"artists"`
	AvailableMarkets []string          `json:"available_markets"`
	DiscNumber       int               `json:"disc_number"`
	Duration         int               `json:"duration_ms"`
	Explicit         bool              `json:"explicit"`
	ExternalURLs     map[string]string `json:"external_urls"`
	Endpoint         string            `json:"href"`
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	PreviewURL       string            `json:"preview_url"`
	TrackNumber      int               `json:"track_number"`
	URI              string            `json:"uri"`
	Type             string            `json:"type"`
}

type PlaylistTrack struct {
	AddedAt string    `json:"added_at"`
	AddedBy string    `json:"added_by"`
	IsLocal bool      `json:"is_local"`
	Track   FullTrack `json:"track"`
}
