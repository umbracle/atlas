package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var file = "atlas"

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Debug ngrok server`))
}

func atlasDownload(w http.ResponseWriter, r *http.Request) {
	// read file every time
	downloadBytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", http.DetectContentType(downloadBytes))
	w.Header().Set("Content-Disposition", "attachment; filename="+file+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(len(string(downloadBytes))))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/atlas", atlasDownload)

	log.Print("Ngrok downstream server started at localhost:8080")
	http.ListenAndServe(":8080", nil)
}
