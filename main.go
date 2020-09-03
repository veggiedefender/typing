package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	var (
		imagesPath  = flag.String("static", "./images", "path to images")
		bgImagePath = flag.String("bg", "./images/main.png", "path to images")
		fontPath    = flag.String("font", "./Roboto-Regular.ttf", "path to font")
		camoURL     = flag.String("camoURL", "", "anonymized GitHub camo URL")
		enableHTTPS = flag.Bool("https", false, "enable HTTPS")
		certPath    = flag.String("certs", "./cert-cache", "path to letsencrypt autocert cache directory")
	)
	flag.Parse()

	if len(*camoURL) == 0 {
		log.Fatal("Pass in GitHub camo URL with -camoURL <url>")
	}

	scrn, err := NewScreen(*bgImagePath, *fontPath)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/k/").Handler(http.StripPrefix("/k/", http.FileServer(http.Dir(*imagesPath))))

	r.Handle("/type/{character:[a-z0-9]|backspace|comma|space|period|enter}", TypeHandler(scrn, *camoURL))
	r.Handle("/screen.gif", RenderHandler(scrn))

	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if *enableHTTPS {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(*certPath),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("kbd.jse.li"),
		}
		srv.Addr = ":https"
		srv.TLSConfig = m.TLSConfig()
		log.Fatal(srv.ListenAndServeTLS("", ""))
	} else {
		srv.Addr = ":8000"
		log.Fatal(srv.ListenAndServe())
	}
}
