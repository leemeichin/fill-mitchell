package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/go-chi/chi/v5"
)

func index(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("../assets/templates/index.html")

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write(html)
}

func randomPlaceholder(w http.ResponseWriter, r *http.Request) {
	imagePath := fmt.Sprintf("../assets/placeholders/%d.jpg", rand.Intn(9)+1)
	imageBytes, err := ioutil.ReadFile(imagePath)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)
}

func randomPlaceholderResized(w http.ResponseWriter, r *http.Request) {
	width, _ := strconv.Atoi(chi.URLParam(r, "width"))
	height, _ := strconv.Atoi(chi.URLParam(r, "height"))

	imagePath := fmt.Sprintf("../assets/placeholders/%d.jpg", rand.Intn(9)+1)
	img, err := vips.NewImageFromFile(imagePath)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = img.Thumbnail(width, height, vips.InterestingCentre)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	exportParams := vips.NewDefaultJPEGExportParams()
	imageBytes, _, err := img.Export(exportParams)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)
}

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	r := chi.NewRouter()

	r.Get("/r", randomPlaceholder)
	r.Get("/r/{width:[0-9]+}/{height:[0-9]+}", randomPlaceholderResized)
	r.Get("/", index)

	http.ListenAndServe("0.0.0.0:3003", r)
}
