package models

type Base struct {
	URL      string `json:"href"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type Playlist struct {
	Collaborative bool              `json:"collaborative"`
	Description   string            `json:"description"`
	ExternalURLs  map[string]string `json:"external_urls"`
	URL           string            `json:"href"`
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	Public        bool              `json:"public"`
	SnapshotID    string            `json:"snapshot_id"`
	Tracks        PlaylistTracks    `json:"tracks"`
	URI           string            `json:"uri"`
}

type PlaylistTracks struct {
	URL   string `json:"href"`
	Total int    `json:"total"`
}

type PlaylistTracksPage struct {
	Base
	Tracks []PlaylistTrack `json:"items"`
}

type PlaylistResponse struct {
	Base
	Playlists []Playlist `json:"items"`
}
