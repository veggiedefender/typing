package main

import (
	"crypto/sha1"
	"encoding/hex"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"sync"

	"github.com/fogleman/gg"
)

// Screen is a canvas that accepts input and can render an image
type Screen struct {
	bgImage  image.Image
	width    int
	height   int
	fontPath string
	bufMux   sync.RWMutex
	buf      []rune
}

// NewScreen creates a Screen with a background image
func NewScreen(bgImagePath, fontPath string) (*Screen, error) {
	bgImage, err := gg.LoadImage(bgImagePath)
	if err != nil {
		return nil, err
	}

	rect := bgImage.Bounds()
	return &Screen{
		bgImage:  bgImage,
		width:    rect.Max.X - rect.Min.X,
		height:   rect.Max.Y - rect.Min.Y,
		fontPath: fontPath,
	}, nil
}

// Add adds a character to the screen's buffer. Passing in '\b' will delete the last character
// if it exists.
func (s *Screen) Add(ch rune) {
	s.bufMux.Lock()
	defer s.bufMux.Unlock()

	if ch == '\b' {
		if len(s.buf) > 0 {
			s.buf = s.buf[:len(s.buf)-1]
		}
	} else {
		s.buf = append(s.buf, ch)
	}
}

// Render renders the screen's buffer
func (s *Screen) Render(w io.Writer) (string, error) {
	s.bufMux.RLock()
	message := string(s.buf)
	s.bufMux.RUnlock()

	frame1, err := s.renderString(message)
	if err != nil {
		return "", err
	}

	frame2, err := s.renderString(message + "|")
	if err != nil {
		return "", err
	}

	palettedImage1 := image.NewPaletted(frame1.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(palettedImage1, frame1.Bounds(), frame1, image.ZP)
	palettedImage2 := image.NewPaletted(frame2.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(palettedImage2, frame2.Bounds(), frame2, image.ZP)
	if err != nil {
		return "", err
	}

	gif.EncodeAll(w, &gif.GIF{
		Image: []*image.Paletted{
			palettedImage1,
			palettedImage2,
		},
		Delay: []int{50, 50},
	})
	return etag(message), nil
}

func (s *Screen) renderString(message string) (image.Image, error) {
	dc := gg.NewContext(s.width, s.height)
	dc.DrawImage(s.bgImage, 0, 0)

	if err := dc.LoadFontFace(s.fontPath, 18); err != nil {
		return nil, err
	}

	textRightMargin := 10.0
	textTopMargin := 30.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin

	dc.SetColor(color.Black)
	dc.DrawStringWrapped(message, x, y, 0, 0, maxWidth, 1.6, gg.AlignLeft)
	return dc.Image(), nil
}

func etag(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
