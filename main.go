package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/veggiedefender/typing/screen"

	"github.com/gorilla/mux"
)

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Println(e.Error.Error())
		http.Error(w, e.Message, e.Code)
	}
}

func main() {
	var (
		imagesPath  = flag.String("static", "./images", "path to images")
		bgImagePath = flag.String("bg", "./images/main.png", "path to images")
		fontPath    = flag.String("font", "./Roboto-Regular.ttf", "path to font")
	)
	flag.Parse()

	scr, err := screen.NewScreen(*bgImagePath, *fontPath)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/k/").Handler(http.StripPrefix("/k/", http.FileServer(http.Dir(*imagesPath))))

	r.Handle("/type/{character:[a-z0-9]|backspace|comma|space|period|enter}", TypeHandler(scr))
	r.Handle("/text.png", RenderHandler(scr))

	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	srv.Addr = ":8000"
	log.Fatal(srv.ListenAndServe())
}
