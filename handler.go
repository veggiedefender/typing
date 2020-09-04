package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var keymap = map[string]rune{
	"q":         'q',
	"w":         'w',
	"e":         'e',
	"r":         'r',
	"t":         't',
	"y":         'y',
	"u":         'u',
	"i":         'i',
	"o":         'o',
	"p":         'p',
	"a":         'a',
	"s":         's',
	"d":         'd',
	"f":         'f',
	"g":         'g',
	"h":         'h',
	"j":         'j',
	"k":         'k',
	"l":         'l',
	"z":         'z',
	"x":         'x',
	"c":         'c',
	"v":         'v',
	"b":         'b',
	"n":         'n',
	"m":         'm',
	"0":         '0',
	"1":         '1',
	"2":         '2',
	"3":         '3',
	"4":         '4',
	"5":         '5',
	"6":         '6',
	"7":         '7',
	"8":         '8',
	"9":         '9',
	"backspace": '\b',
	"comma":     ',',
	"space":     ' ',
	"period":    '.',
	"enter":     '\n',
}

// RenderHandler renders the current screen
func RenderHandler(scrn *Screen) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		etag, err := scrn.Render(&buf)
		if err != nil {
			http.Error(w, "error rendering screen", 500)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", etag)

		w.Write(buf.Bytes())
	})
}

// TypeHandler types a character to the screen
func TypeHandler(scrn *Screen, repoURL string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch := mux.Vars(r)["character"]
		scrn.Add(keymap[ch])
		log.Printf("Pressed button: %q", ch)
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, repoURL, 302)
	})
}

func purgeGitHubCache(camoURL string) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				req, err := http.NewRequestWithContext(ctx, "PURGE", camoURL, nil)
				if err != nil {
					return
				}
				http.DefaultClient.Do(req)
			}()
		}
	}

}
