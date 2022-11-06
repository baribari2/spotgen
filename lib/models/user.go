package models

type User struct {
	DisplayName  string            `json:"display_name"`
	ExternalURLs map[string]string `json:"external_urls"`
	URL          string            `json:"href"`
	Id           string            `json:"id"`
	URI          string            `json:"uri"`
}

type PUser struct {
	User
	Country   string `json:"country"`
	Email     string `json:"email"`
	Product   string `json:"product"`
	Birthdate string `json:"birthdate"`
}
