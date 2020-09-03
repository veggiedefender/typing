package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/veggiedefender/typing/screen"
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

// TypeHandler types a character to the screen
func TypeHandler(scr *screen.Screen) http.Handler {
	return appHandler(func(w http.ResponseWriter, r *http.Request) *appError {
		vars := mux.Vars(r)
		log.Printf("Pressed button: %q", vars["character"])
		scr.Add(keymap[vars["character"]])
		http.Redirect(w, r, "https://github.com/veggiedefender/keyboard", 302)
		return nil
	})
}

// RenderHandler renders the current screen
func RenderHandler(scr *screen.Screen) http.Handler {
	return appHandler(func(w http.ResponseWriter, r *http.Request) *appError {
		var buf bytes.Buffer
		etag, err := scr.Render(&buf)
		if err != nil {
			return &appError{Error: err, Message: "error rendering image", Code: 500}
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", etag)

		w.Write(buf.Bytes())
		return nil
	})
}
