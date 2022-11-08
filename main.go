package main

import (
	"flag"

	"log"
	"net/http"
	"os"
	"sync"

	"github.com/baribari2/spotify-playlist-generator/pkg/models"
)

/*
Display images on CLI?

To-Do:

	Switch generate functions to goroutines
	Use wrapper (Second thought: it's a big pkg)
*/
func main() {
	featured := flag.NewFlagSet("feat", flag.PanicOnError)
	flength := featured.String("len", "50", "Length of the playlist to be created")
	fname := featured.String("name", "", "Name of the playlist to be created")
	fdesc := featured.String("desc", "", "Description of the playlist to be created. May be left blank")
	fpublic := featured.Bool("pub", true, "Publicity of the playlist to be created")
	fcollab := featured.Bool("collab", false, "Collaboration capabilities of the playlist to be created")

	recommended := flag.NewFlagSet("rec", flag.PanicOnError)
	rlength := recommended.String("len", "50", "Length of the playlist to be created")
	rname := recommended.String("name", "", "Name of the playlist to be created")
	rdesc := recommended.String("desc", "", "Length of the playlist to be created")
	rartists := recommended.String("art", "", "Seed artists for playlist generation (Comma-separated list: a,b,c), total seed items must not exceed 5")
	rgenres := recommended.String("gen", "", "Seed artists for playlist generation (Comma-separated list: a,b,c), total seed items must not exceed 5")
	rpublic := recommended.Bool("pub", true, "Publicity of the playlist to be created")
	rcollab := recommended.Bool("collab", false, "Collaboration capabilities of the playlist to be created")
	//rtracks := recommended.String("len", "50", "Length of the playlist to be created")

	var (
		token  = &models.TokenResponse{}
		server = &http.Server{
			Addr: ":8888",
		}
		wg       = &sync.WaitGroup{}
		playlist string
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

		playlist, err = generateFeatured(*flength, *fname, *fdesc, *fpublic, *fcollab, token)
		if err != nil {
			log.Printf("Failed to generate `featured` playlist: %v", err.Error())
		}

	case "rec":

		err := recommended.Parse(os.Args[2:])
		if err != nil {
			log.Printf("Failed to parse OS arguments: %v", err.Error())
		}

		playlist, err = generateRecommended(*rlength, *rname, *rdesc, *rpublic, *rcollab, *rgenres, *rartists, token)
		if err != nil {
			log.Printf("Failed to generate `recommended` playlist: %v", err.Error())
		}
	default:
		log.Println("Expected a command (feat, word, rec)")
	}

	log.Printf(">>>    Generated playlist: %v    <<<", playlist)
}
