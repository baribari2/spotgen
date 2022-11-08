package models

type Artist struct {
	Name         string            `json:"name"`
	Id           string            `json:"id"`
	URI          string            `json:"uri"`
	URL          string            `json:"href"`
	ExternalURLS map[string]string `json:"external_urls"`
}
