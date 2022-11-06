package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/broothie/qst"
	"github.com/spotify-playlist-generator/lib/models"
)

// Generates a playlist using songs from the users featured section
func generateFeatured(length string, name string, public bool, collab bool, desc string, token *models.TokenResponse) (string, error) {
	if name == "" {
		return "", errors.New("Missing name parameter")
	}

	if length == "" {
		length = "50"
	}

	base := "https://api.spotify.com/v1"

	//GET request to obtain information about the current user
	res, err := qst.Get(
		base+"/me",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
	)
	if err != nil {
		return "", err
	}

	var currUser models.PUser
	err = json.NewDecoder(res.Body).Decode(&currUser)
	if err != nil {
		return "", err
	}

	// GET request to obtain featured playlists
	res, err = qst.Get(
		base+"/browse/featured-playlists",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.QueryValue("limit", length),
	)

	if err != nil {
		return "", err
	}

	var featured struct {
		Playlists models.PlaylistResponse `json:"playlists"`
		Message   string                  `json:"message"`
	}
	err = json.NewDecoder(res.Body).Decode(&featured)

	tracks := []string{}

	// Append track endpoint from every playlist to `tracks` slice
	for _, p := range featured.Playlists.Playlists {
		tracks = append(tracks, p.URL)
	}

	// POST request to create playlist
	res, err = qst.Post(
		base+"/users/"+currUser.Id+"/playlists",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.BodyJSON(
			map[string]interface{}{
				"name":          name,
				"public":        public,
				"collaborative": collab,
				"description":   desc,
			},
		),
	)

	if err != nil {
		return "", err
	}

	var p struct {
		models.Playlist
	}
	err = json.NewDecoder(res.Body).Decode(&p)

	// Parsing the track endoints to obtain track id to use with spotify uri
	uris := []string{}
	for _, t := range tracks {
		spl := strings.SplitAfter(t, "/")

		uris = append(uris, ("spotify:track:" + spl[len(spl)-1]))
	}

	uris = uris[:len(uris)-1]

	log.Printf("URIS: %v \n", uris)

	//POST request to add tracks
	ress, err := qst.Post(
		base+"/playlists/"+p.Id+"/tracks",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.QueryValue("position", "0"),
		qst.BodyJSON(
			map[string]interface{}{
				"uris": uris,
			},
		),
	)

	log.Printf("REQUEST: %v", ress)

	if err != nil {
		return "", nil
	}

	//log.Printf("RES: %v \n", res)

	return p.URL, nil
}
