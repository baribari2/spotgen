package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/baribari2/spotify-playlist-generator/pkg/models"
	"github.com/broothie/qst"
)

// Generates a playlist using songs from the users featured section
func generateFeatured(length, name, desc string, public, collab bool, token *models.TokenResponse) (string, error) {
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
			if len(tracks) > 99 {
				continue
			}

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
		qst.Header("Content-Type", "application/json"),
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

	var pl struct {
		models.Playlist
	}
	err = json.NewDecoder(res.Body).Decode(&pl)

	// POST request to add tracks
	res, err = qst.Post(
		base+"/playlists/"+pl.Id+"/tracks",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.Header("Content-Type", "application/json"),
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

	return pl.URL, nil
}

func generateRecommended(length, name, desc string, public, collab bool, gen, art string, token *models.TokenResponse) (string, error) {
	if name == "" {
		return "", errors.New("Missing name parameter")
	}

	if length == "" {
		length = "50"
	}

	artists := strings.SplitAfter(art, ",")
	genres := strings.SplitAfter(gen, ",")
	searchQuery := []string{}

	if len(artists)+len(genres) >= 5 {
		return "", errors.New("Too many seed items")
	}

	for _, v := range artists {
		t := strings.SplitAfter(v, " ")

		// log.Printf("Artists after split: %v", t)

		if len(t) > 1 {
			var te string
			for i, tm := range t {
				if i == len(t)-1 {
					te += tm
				} else {
					te += tm + "%"
				}
			}

			searchQuery = append(searchQuery, te)

		} else {
			searchQuery = append(searchQuery, t[0])
		}
	}

	base := "https://api.spotify.com/v1"

	//GET /me (move function call to main to keep modular)
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

	var uriQuery string
	var artist struct {
		models.ArtistResponse `json:"artists"`
	}

	// GET requests to obtain artist id's
	for i := range artists {
		res, err = qst.Get(
			base+"/search",
			qst.Header("Authorization", "Bearer "+token.AccessToken),
			qst.Header("Content-Type", "application/json"),
			qst.QueryValue("type", "artist"),
			qst.QueryValue("q", searchQuery[i]),
		)

		if err != nil {
			return "", err
		}

		err = json.NewDecoder(res.Body).Decode(&artist)
		if err != nil {
			return "", err
		}

		uriQuery += artist.Artists[0].Id + ","
	}

	//GET recommendations
	res, err = qst.Get(
		base+"/recommendations",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.Header("Content-Type", "application/json"),
		qst.QueryValue("seed_artists", uriQuery),
		qst.QueryValue("seed_genres", gen),
		qst.QueryValue("limit", length),
	)

	if err != nil {
		return "", err
	}

	var recs models.Recommendations
	err = json.NewDecoder(res.Body).Decode(&recs)

	if err != nil {
		return "", err
	}

	uris := []string{}
	for _, r := range recs.Tracks {
		uris = append(uris, "spotify:track:"+r.Id)
	}

	//POST Create playlist
	res, err = qst.Post(
		base+"/users/"+currUser.Id+"/playlists",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.Header("Content-Type", "application/json"),
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

	//POST add to playlist
	res, err = qst.Post(
		base+"/playlists/"+p.Id+"/tracks",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.Header("Content-Type", "application/json"),
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

	return p.URL, nil
}
