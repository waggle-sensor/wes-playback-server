package main

import (
	"bytes"
	"flag"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
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
	if _, err := os.Stat(name); os.IsNotExist(err) {
		log.Printf("No video \"%s\" - falling back to blank video.", name)
		return blankMJPEGHandler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}

// makeImageHandle creates a handler to serve static images in sequence from a directory.
// It will serve the n-th image mod number of files after n seconds.
func makeImageHandler(name string) http.HandlerFunc {
	match := filepath.Join(name, "*.jpg")
	files, err := filepath.Glob(match)

	if err != nil || len(files) == 0 {
		log.Printf("No images found in \"%s\" - falling back to blank image.", name)
		return blankJPEGHandler
	}

	log.Printf("Found %d images in \"%s.\"", len(files), match)

	makeTime := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		i := int(time.Since(makeTime).Seconds())
		imageName := files[i%len(files)]
		http.ServeFile(w, r, imageName)
	}
}

// makeJPEGImageBuffer encodes a blank JPEG image to a buffer
func makeJPEGImageBuffer() *bytes.Buffer {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	jpeg.Encode(&buf, img, nil)
	return &buf
}

// blankImageBuffer is a cached blank JPEG image for both image and video endpoints.
var blankImageBuffer = makeJPEGImageBuffer()

// blankJPEGHandler will serve a black JPEG image.
func blankJPEGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(blankImageBuffer.Bytes())
}

// blankMJPEGHandler will serve a black MJPEG video
func blankMJPEGHandler(w http.ResponseWriter, r *http.Request) {
	mimeWriter := multipart.NewWriter(w)
	defer mimeWriter.Close()

	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mimeWriter.Boundary())

	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	for {
		pw, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			return
		}
		if _, err := pw.Write(blankImageBuffer.Bytes()); err != nil {
			return
		}
		time.Sleep(33 * time.Millisecond)
	}
}

func main() {
	dataDir := flag.String("data", ".", "Path to data directory.")
	flag.Parse()

	log.Printf("Serving data from %s.", *dataDir)
	os.Chdir(*dataDir)

	// setup bottom camera interface
	http.HandleFunc("/bottom/live.mp4", makeStreamHandler("bottom/live.mp4"))
	http.HandleFunc("/bottom/image.jpg", makeImageHandler("bottom/images/"))

	// setup top camera interface
	http.HandleFunc("/top/live.mp4", makeStreamHandler("top/live.mp4"))
	http.HandleFunc("/top/image.jpg", makeImageHandler("top/images/"))

	log.Fatal(http.ListenAndServe(":8090", nil))
}
