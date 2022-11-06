package main

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONResponse(res *http.Response, target interface{}) error {
	return json.NewDecoder(res.Body).Decode(target)
}
