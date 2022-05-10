/*
Serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Navigating to http://localhost:3004 will display the index.html or directory
listing file.
*/
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("p", "3004", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	os.Mkdir("img", os.ModePerm)
	for i := uint8(0); i < 255; i++ {
		makeImg(i)
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func makeImg(i uint8) {
	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	for x := 0; x < 100; x++ {
		for y := 0; y < 50; y++ {
			img.Set(x, y, color.White)
			if x > 40 && x < 60 {
				img.Set(x, y, color.RGBA{255 - i, i, 0, 255})
			}
		}
	}

	// Save to .png
	fname := fmt.Sprintf("img/%v.png", i)
	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
