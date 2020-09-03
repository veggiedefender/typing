package screen

import (
	"crypto/sha1"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
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

	dc := gg.NewContext(s.width, s.height)
	dc.DrawImage(s.bgImage, 0, 0)

	if err := dc.LoadFontFace(s.fontPath, 50); err != nil {
		return "", err
	}

	textRightMargin := 60.0
	textTopMargin := 60.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin

	dc.SetColor(color.Black)
	dc.DrawStringWrapped(message, x, y, 0, 0, maxWidth, 1.6, gg.AlignLeft)

	// dc.MeasureMultilineString

	png.Encode(w, dc.Image())
	log.Printf("Rendered: %q", message)
	return etag(message), nil
}

func etag(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
