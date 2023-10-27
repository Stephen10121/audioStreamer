package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dhowden/tag"
)

type audioData struct {
	Title   string      `json:"title"`
	Artist  string      `json:"artist"`
	Album   string      `json:"album"`
	Picture tag.Picture `json:"picture"`
}

func audioInfo(fileSource string) (audioData, error) {
	f, err := os.Open("server/songs/" + fileSource)
	if err != nil {
		return audioData{}, err
	}
	defer f.Close()
	m, err := tag.ReadFrom(f)
	if err != nil {
		return audioData{}, err
	}
	log.Println(m)
	log.Print(m.Format()) // The detected format.
	log.Print(m.Title())  // The title of the track (see Metadata interface for more details).
	log.Print(m.Artist())
	log.Print(m.Album())
	log.Print(m.Picture())
	return audioData{
		Title:   m.Title(),
		Artist:  m.Artist(),
		Album:   m.Album(),
		Picture: *m.Picture(),
	}, nil
}

func main() {
	// configure the songs directory name and port
	const songsDir = "server"
	const port = 3001

	http.HandleFunc("/songData/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var slug string
		if strings.HasPrefix(r.URL.Path, "/songData/") {
			slug = r.URL.Path[len("/songData/"):]
		}
		info, err := audioInfo(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - FileData not found"))
			return
		}
		jsonData, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something went wrong"))
			return
		}
		fmt.Fprintf(w, fmt.Sprintf(string(jsonData), jsonData))
	})
	// add a handler for the song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(songsDir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
