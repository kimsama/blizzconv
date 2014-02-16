// Package cl2 implements a CL2 image decoder.
//
// The CL2 format is the second version of the CEL format. It uses run-length
// encoding to decrease the size of images. Other than this addition, the format
// itself is identical to the CEL format.
package cl2

import (
	"image"
	"path"

	"github.com/mewrnd/blizzconv/images/cel"
)

// DecodeAll returns the sequential frames of a CEL or CL2 image based on a
// given conf.
func DecodeAll(imgName string, conf *cel.Config) (imgs []image.Image, err error) {
	// Decode CEL version 1 images using the cel package.
	if path.Ext(imgName) == ".cel" {
		return cel.DecodeAll(imgName, conf)
	}

	// Get frame contents.
	frames, err := cel.GetFrames(imgName)
	if err != nil {
		return nil, err
	}

	// Decode frames.
	for frameNum, frame := range frames {
		width, ok := conf.FrameWidth[frameNum]
		if !ok {
			// Use default frame width.
			width = conf.Width
		}
		height, ok := conf.FrameHeight[frameNum]
		if !ok {
			// Use default frame height.
			height = conf.Height
		}

		// Decode frame.
		img := DecodeFrameType6(frame, width, height, conf.Pal)
		imgs = append(imgs, img)
	}

	return imgs, nil
}
