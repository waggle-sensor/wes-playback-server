package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// makeStreamHandler creates a handler to serve static video content.
//
// For the purposes of video streaming, we assume that http.ServeFile will add the headers:
// Accept-Ranges: bytes
// Content-Type: *video mimetype*
func makeStreamHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}

// makeImageHandle creates a handler to serve static images in sequence from a directory.
// It will serve the n-th image mod number of files after n seconds.
func makeImageHandler(name string) http.HandlerFunc {
	match := filepath.Join(name, "*.jpg")
	files, err := filepath.Glob(match)
	log.Printf("found %d images matching %s", len(files), match)

	if err != nil || len(files) == 0 {
		return http.NotFound
	}

	makeTime := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		i := int(time.Since(makeTime).Seconds())
		imageName := files[i%len(files)]
		http.ServeFile(w, r, imageName)
	}
}

func main() {
	dataDir := flag.String("data", ".", "path to data directory")
	flag.Parse()

	log.Printf("serving data from %s", *dataDir)
	os.Chdir(*dataDir)

	// setup current ffserver interface
	http.HandleFunc("/live", makeStreamHandler("bottom/live.mp4"))

	// setup bottom camera interface
	http.HandleFunc("/bottom/live.mp4", makeStreamHandler("bottom/live.mp4"))
	http.HandleFunc("/bottom/image.jpg", makeImageHandler("bottom/images/"))

	// setup top camera interface
	http.HandleFunc("/top/live.mp4", makeStreamHandler("top/live.mp4"))
	http.HandleFunc("/top/image.jpg", makeImageHandler("top/images/"))

	log.Fatal(http.ListenAndServe(":8090", nil))
}
