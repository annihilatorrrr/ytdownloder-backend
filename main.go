package main

import (
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
	"os"
)

var yt youtube.Client

func processor(w http.ResponseWriter, r bunrouter.Request) error {
	id := r.Params().ByName("id")
	quality := r.Params().ByName("quality")
	if id == "" {
		_, _ = fmt.Fprint(w, "error: Id field missing!")
		return nil
	}
	d := downloader.Downloader{Client: yt}
	v, err := d.GetVideo(id)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return nil
	}
	q := v.Formats.WithAudioChannels()[1].URL
	if quality == "720" {
		q = v.Formats.WithAudioChannels()[0].URL
	}
	http.Redirect(w, r.Request, q, http.StatusFound)
	return nil
}

func welcm(w http.ResponseWriter, _ bunrouter.Request) error {
	_, _ = fmt.Fprint(w, "Welcome!")
	return nil
}

func okand(w http.ResponseWriter, _ bunrouter.Request) error {
	_, _ = fmt.Fprint(w, "error: Follow this format: download/video_id/quality ( by default 480p. )")
	return nil
}

func main() {
	log.Println("Starting ...")
	router := bunrouter.New()
	router.GET("/", welcm)
	router.GET("/download", okand)
	router.GET("/download/:id/:quality", processor)
	router.GET("/download/:id", processor)
	port := os.Getenv("PORT")
	handler := http.Handler(router)
	if port == "" {
		port = "80"
	}
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handler,
	}
	yt = youtube.Client{}
	log.Println("Started!")
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	log.Println("Bye!")
}
