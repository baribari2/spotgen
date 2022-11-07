package main

import (
	"encoding/json"
	"errors"
	"log"

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

	// To-Do: switch to goroutines

	//Make GET requests to obtain playlist(s) items
	for _, p := range featured.Playlists.Playlists {
		if len(tracks) > 99 {
			continue
		}

		res, err = qst.Get(
			base+"/playlists/"+p.Id+"/tracks",
			qst.Header("Authorization", "Bearer "+token.AccessToken),
			qst.Header("Content-Type", "application/json"),
		)

		if err != nil {
			return "", err
		}

		var result models.PlaylistTracksPage
		err = json.NewDecoder(res.Body).Decode(&result)

		// Append track endpoints from every playlist to `tracks` slice
		for _, t := range result.Tracks {
			tracks = append(tracks, t.Track.Id)
		}
	}

	// Parsing the track endpoints to obtain track id to use with spotify uri
	uris := []string{}
	for _, id := range tracks {
		uris = append(uris, ("spotify:track:" + id))
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

	// POST request to add tracks
	res, err = qst.Post(
		base+"/playlists/"+p.Id+"/tracks",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.BodyJSON(
			map[string]interface{}{
				"uris":     uris,
				"position": 0,
			},
		),
	)

	if err != nil {
		return "", nil
	}

	log.Printf("RESPONSE: %v", res)

	return p.URL, nil
}
