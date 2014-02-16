package cel

import (
	"fmt"
	"image/color"
	"io/ioutil"

	"github.com/mewrnd/blizzconv/mpq"
)

// GetPal parses the provided PAL file and returns it as a color.Palette. Below
// is a description of the PAL format.
//
// PAL format:
//    c [256]Color
//
// Color format:
//    r byte   // red
//    g byte   // green
//    b byte   // blue
//
// Note: The absolute path of relPalPath is relative to mpq.ExtractPath.
func GetPal(relPalPath string) (pal color.Palette, err error) {
	palPath := mpq.AbsPath(relPalPath)
	buf, err := ioutil.ReadFile(palPath)
	if err != nil {
		return nil, err
	}
	if len(buf) != 256*3 {
		return nil, fmt.Errorf("cel.GetPal: invalid pal size (%d) for %q", len(buf), relPalPath)
	}
	pal = make(color.Palette, 256)
	for i := range pal {
		c := color.RGBA{
			R: buf[3*i],
			G: buf[3*i+1],
			B: buf[3*i+2],
			A: 0xFF,
		}
		pal[i] = c
	}
	return pal, nil
}
