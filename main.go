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
	"time"
)

var yt youtube.Client

func processor(w http.ResponseWriter, r bunrouter.Request) error {
	parms := r.Params()
	url := parms.ByName("url")
	if url == "" {
		_, _ = fmt.Fprint(w, "error: Url field missing!")
		return nil
	}
	d := downloader.Downloader{Client: yt}
	v, err := d.GetVideo(url)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return nil
	}
	form := &youtube.Format{
		ItagNo:   18,
		Quality:  "medium",
		MimeType: "video/mp4",
	}
	out := fmt.Sprintf("./%s.mp4", v.Title)
	err = d.Download(context.Background(), v, form, out)
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
	router.GET("/download/:url", processor)
	port := os.Getenv("PORT")
	handler := http.Handler(router)
	if port == "" {
		port = "80"
	}
	server := &http.Server{
		Addr:        "0.0.0.0:" + port,
		ReadTimeout: 3 * time.Second,
		Handler:     handler,
	}
	yt = youtube.Client{}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	log.Println("Started!")
}
