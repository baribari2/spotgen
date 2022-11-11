package main

import (
	"context"
	b64 "encoding/base64"
	"log"
	"net/http"
	"sync"

	"github.com/baribari2/spotgen/pkg/models"
	"github.com/broothie/qst"
)

func initAuth(server *http.Server, token *models.TokenResponse, wg *sync.WaitGroup) {
	var (
		accounts = "https://accounts.spotify.com"
		code     string
		state    = "12345"
	)

	log.Println(">>>   Starting Spotify Playlist Generator   <<<")

	// Assign handlers to routes
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		get, err := qst.NewGet(
			accounts+"/authorize",
			qst.QueryValue("client_id", CLIENT_ID),
			qst.QueryValue("response_type", "code"),
			qst.QueryValue("redirect_uri", "http://localhost:8888/"),
			qst.QueryValue("scope", "playlist-read-private playlist-read-collaborative playlist-modify-private playlist-modify-public user-read-playback-position user-top-read user-read-recently-played user-library-modify user-library-read"),
			qst.QueryValue("state", "12345"),
		)

		http.Redirect(w, req, get.URL.String(), http.StatusSeeOther)

		if err != nil {
			log.Printf("Failed to make initial auth request: %v", err.Error())
			return
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		//log.Printf("Callback Request: %v \n \n", req)

		if c := req.URL.Query().Has("code"); !c {
			log.Printf("Request does not contain 'code'")
			return
		} else {
			code = req.URL.Query().Get("code")
		}

		if s := req.URL.Query().Has("state"); !s {
			log.Printf("Request does not contain 'state'")
			return
		} else if st := req.URL.Query().Get("state"); st != state {
			log.Printf("State mistatch: %v", st)
			return
		}

		r, err := qst.Post(
			accounts+"/api/token",
			qst.Header("Content-Type", "application/x-www-form-urlencoded"),
			qst.Header("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(CLIENT_ID+":"+CLIENT_SECRET))),
			qst.QueryValue("grant_type", "authorization_code"),
			qst.QueryValue("code", code),
			qst.QueryValue("redirect_uri", "http://localhost:8888/"),
		)

		if err != nil {
			log.Printf("Failed to make initial token request: %v", err.Error())
			return
		}

		err = DecodeJSONResponse(r, token)
		if err != nil {
			log.Printf("Failed to decode JSON response body: %v", err.Error())
			return
		}

		wg.Done()

		log.Println(">>>   Authentication completed!   <<<")
	})

	log.Println(">>>   Please navigate to http://localhost:8888/login to iniate authentication   <<<")

	// Start server
	wg.Add(1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen and serve on port 8888: %v", err.Error())
			return
		}

		defer server.Shutdown(context.Background())
	}()

}
