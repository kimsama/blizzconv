package cel

import "fmt"
import "image/color"
import "io/ioutil"

import "github.com/mewrnd/blizzconv/mpq"

// GetPal returns a color.Palette created from relPalPath. Below is a
// description of the PAL format.
//
// PAL format:
//    c [256]Color
//
// Color format:
//    r byte   // red
//    g byte   // green
//    b byte   // blue
func GetPal(relPalPath string) (pal color.Palette, err error) {
	palPath := mpq.ExtractPath + relPalPath
	buf, err := ioutil.ReadFile(palPath)
	if err != nil {
		return nil, err
	}
	if len(buf) != 256*3 {
		return nil, fmt.Errorf("invalid pal size (%d) for %q.", len(buf), relPalPath)
	}
	pal = make(color.Palette, 256)
	for i := range pal {
		pal[i] = color.RGBA{buf[3*i], buf[3*i+1], buf[3*i+2], 0xFF}
	}
	return pal, nil
}
