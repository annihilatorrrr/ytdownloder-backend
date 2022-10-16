package main

import (
	"context"
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
	out := fmt.Sprintf("./%s.mp4", v.Title)
	err = d.Download(context.Background(), v, &v.Formats.WithAudioChannels()[1], out)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return nil
	}
	http.ServeFile(w, r.Request, out)
	w.WriteHeader(http.StatusOK)
	return nil
}

func main() {
	log.Println("Starting ...")
	router := bunrouter.New()
	router.GET("/download/:id", processor)
	port := os.Getenv("PORT")
	handler := http.Handler(router)
	if port == "" {
		port = "80"
	}
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		ReadTimeout:  120,
		Handler:      handler,
		WriteTimeout: 120,
	}
	yt = youtube.Client{}
	log.Println("Started!")
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	log.Println("Bye!")
}
