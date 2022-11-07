package main

import (
	"context"
	b64 "encoding/base64"
	"flag"

	//"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/broothie/qst"
	"github.com/spotify-playlist-generator/lib/models"
)

func initAuth(server *http.Server, token *models.TokenResponse, wg *sync.WaitGroup) {
	var (
		accounts = "https://accounts.spotify.com"
		code     string
		state    = "12345"
	)

	log.Println(">>>   Starting Spotify Playlist Generator   <<<")

	// if ACCESS_TOKEN != "" && REFRESH_TOKEN != "" {
	// 	token.AccessToken = ACCESS_TOKEN
	// 	token.RefreshToken = REFRESH_TOKEN

	// 	return
	// }

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

		// _ = os.Setenv("ACCESS_TOKEN", token.AccessToken)
		// _ = os.Setenv("REFRESH_TOKEN", token.RefreshToken)

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

/*
Change to 127.0.0.1?
Display images on CLI?

To-Do:

	Switch generate functions to goroutines
	Increase personalization (Spotify API requests)
*/
func main() {
	featured := flag.NewFlagSet("feat", flag.PanicOnError)
	flength := featured.String("len", "50", "Length of the album to be created")
	fname := featured.String("name", "", "Name of the playlist to be created")
	fpublic := featured.Bool("pub", true, "Publicity of the playlist to be created")
	fcollab := featured.Bool("collab", false, "Collaboration capabilities of the playlist to be created")
	fdesc := featured.String("desc", "", "Description of the playlist to be created. May be left blank")

	var (
		token  = &models.TokenResponse{}
		server = &http.Server{
			Addr: ":8888",
		}
		wg = &sync.WaitGroup{}
	)

	if len(os.Args) < 2 {
		log.Printf("Expected a command (feat, word, rec)")
		os.Exit(1)
	}

	initAuth(server, token, wg)

	wg.Wait()

	switch os.Args[1] {
	case "feat":

		err := featured.Parse(os.Args[2:])
		if err != nil {
			log.Printf("Failed to parse OS arguments: %v", err.Error())
		}

		playlist, err := generateFeatured(*flength, *fname, *fpublic, *fcollab, *fdesc, token)
		if err != nil {
			log.Printf("Failed to generate `feat` playlist: %v", err.Error())
		}

		log.Printf(">>>    Generated playlist: %v    <<<", playlist)
	default:
		log.Println("Expected a command (feat, word, rec)")
	}
}
