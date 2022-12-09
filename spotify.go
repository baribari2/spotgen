package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/baribari2/spotgen/pkg/models"
	"github.com/broothie/qst"
)

// Generates a playlist using songs from the users featured section
func generateFeatured(length, name, desc string, public, collab bool, currUser *models.PUser, token *models.TokenResponse, wg *sync.WaitGroup) (string, error) {
	log.Printf(">>>   Generating featured playlist    <<<")

	if name == "" {
		return "", errors.New("Missing name parameter")
	}

	if length == "" {
		length = "50"
	}

	base := "https://api.spotify.com/v1"

	// GET request to obtain featured playlists
	res, err := qst.Get(
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
	track := make(chan string, 100)

	// GET requests to obtain playlist(s) items
	for _, p := range featured.Playlists.Playlists[len(featured.Playlists.Playlists)/2 : (len(featured.Playlists.Playlists)/2)+1] {
		wg.Add(1)

		go func(p models.Playlist) {

			res, err := qst.Get(
				base+"/playlists/"+p.Id+"/tracks",
				qst.Header("Authorization", "Bearer "+token.AccessToken),
				qst.Header("Content-Type", "application/json"),
			)

			if err != nil {
				return
			}

			var result models.PlaylistTracksPage
			err = json.NewDecoder(res.Body).Decode(&result)
			if err != nil {
				return
			}

			defer wg.Done()

			// Append track ids from every playlist to `tracks` slice
			for _, t := range result.Tracks {
				track <- t.Track.Id
			}
		}(p)
	}

	wg.Wait()

	close(track)
	for t := range track {
		if len(tracks) > 99 {
			continue
		} else {
			tracks = append(tracks, t)
		}
	}

	// Formatting Spotify track URIs
	uris := []string{}
	for _, id := range tracks {
		uris = append(uris, ("spotify:track:" + id))
	}

	pl, err := createPlaylist(name, length, desc, public, collab, currUser, token)
	if err != nil {
		return "", err
	}
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

// Generates a playlist using seed values (artists & genres)
func generateRecommended(length, name, desc string, public, collab bool, gen, art string, currUser *models.PUser, token *models.TokenResponse, wg *sync.WaitGroup) (string, error) {
	log.Printf(">>>   Generating recommended playlist    <<<")

	if name == "" {
		return "", errors.New("Missing name parameter")
	}

	if length == "" {
		length = "50"
	}

	artists := strings.SplitAfter(art, ",")
	genres := strings.SplitAfter(gen, ",")
	searchQuery := []string{}

	if (len(artists)-1)+(len(genres)-1) > 5 {
		return "", errors.New("Too many seed items")
	}

	for _, v := range artists {
		t := strings.SplitAfter(v, " ")

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

	var uriQuery string
	artistC := make(chan *models.Artist, len(artists))

	// GET requests to obtain artist id's
	for i := range artists {
		wg.Add(1)

		go func(i int) {
			res, err := qst.Get(
				base+"/search",
				qst.Header("Authorization", "Bearer "+token.AccessToken),
				qst.Header("Content-Type", "application/json"),
				qst.QueryValue("type", "artist"),
				qst.QueryValue("q", searchQuery[i]),
			)

			if err != nil {
				log.Printf("Failed to make search GET request: %v", err.Error())
				return
			}

			var artist struct {
				models.ArtistResponse `json:"artists"`
			}

			err = json.NewDecoder(res.Body).Decode(&artist)
			if err != nil {
				log.Printf("Failed to make search GET request: %v", err.Error())
				return
			}

			defer wg.Done()
			if len(artist.Artists) > 0 {
				artistC <- &artist.Artists[0]
			}
		}(i)
	}

	wg.Wait()

	close(artistC)
	for a := range artistC {
		uriQuery += a.Id + ","
	}

	//GET recommendations
	res, err := qst.Get(
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

	p, err := createPlaylist(name, length, desc, public, collab, currUser, token)
	if err != nil {
		return "", err
	}

	err = addToPlaylist(p, uris, token)
	if err != nil {
		return "", err
	}

	return p.URL, nil
}

// POST request to create playlist
func createPlaylist(name, length, desc string, public, collab bool, currUser *models.PUser, token *models.TokenResponse) (*models.Playlist, error) {
	base := "https://api.spotify.com/v1"

	res, err := qst.Post(
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
		return nil, err
	}

	var p struct {
		models.Playlist
	}
	err = json.NewDecoder(res.Body).Decode(&p)

	return &p.Playlist, nil
}

// POST add to playlist
func addToPlaylist(playlist *models.Playlist, URIS []string, token *models.TokenResponse) error {
	base := "https://api.spotify.com/v1"

	_, err := qst.Post(
		base+"/playlists/"+playlist.Id+"/tracks",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
		qst.Header("Content-Type", "application/json"),
		qst.BodyJSON(
			map[string]interface{}{
				"uris":     URIS,
				"position": 0,
			},
		),
	)

	if err != nil {
		return err
	}

	return nil
}

// GET request to obtain information about the current user
func getCurrentUser(token *models.TokenResponse) (*models.PUser, error) {
	base := "https://api.spotify.com/v1"

	res, err := qst.Get(
		base+"/me",
		qst.Header("Authorization", "Bearer "+token.AccessToken),
	)

	if err != nil {
		return nil, err
	}

	var currUser models.PUser
	err = json.NewDecoder(res.Body).Decode(&currUser)
	if err != nil {
		return nil, err
	}

	return &currUser, nil
}
