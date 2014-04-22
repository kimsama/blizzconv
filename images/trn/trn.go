// Package trn implements parsing of TRN files.
//
// TRN files contains information about color transitions and may be thought of
// as a palette for palettes. Each TRN file contains 256 color transitions, one
// for each palette index. Below is a description of the TRN file format.
//
// TRN format:
//    // index maps from original to new palette indicies.
//    index [256]uint8
package trn

import (
	"fmt"
	"image/color"
	"io/ioutil"

	"github.com/mewrnd/blizzconv/mpq"
)

// ConvertPal converts the src palette based on the provided TRN file and
// returns it as a color.Palette.
//
// Note: The absolute path of relTrnPath is relative to mpq.ExtractPath.
func ConvertPal(src color.Palette, relTrnPath string) (dst color.Palette, err error) {
	trnPath := mpq.AbsPath(relTrnPath)
	trn, err := ioutil.ReadFile(trnPath)
	if err != nil {
		return nil, err
	}
	if len(trn) != 256 {
		return nil, fmt.Errorf("trn.ConvertPal: invalid TRN size (%d) for %q", len(trn), relTrnPath)
	}

	// ref: 46567D
	dst = make(color.Palette, 256)
	for i := range dst {
		dst[i] = src[trn[i]]
	}

	return dst, nil
}
