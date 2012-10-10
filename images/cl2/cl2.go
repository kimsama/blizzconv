// Package cl2 implements a CL2 image decoder.
//
// The CL2 format is the second version of the CEL format. It uses run length
// encoding to decrease the images' size. Other than this addition, the format
// itself is identical to the CEL format.
package cl2

import "image"
import "path"

import "github.com/mewrnd/blizzconv/images/cel"

// DecodeAll returns the sequential frames of a CEL image based on a given conf.
func DecodeAll(imgName string, conf *cel.Config) (imgs []image.Image, err error) {
	if path.Ext(imgName) == ".cel" {
		return cel.DecodeAll(imgName, conf)
	}
	frames, err := cel.GetFrames(imgName)
	if err != nil {
		return nil, err
	}
	for frameNum, frame := range frames {
		width, ok := conf.FrameWidth[frameNum]
		if !ok {
			width = conf.Width
		}
		height, ok := conf.FrameHeight[frameNum]
		if !ok {
			height = conf.Height
		}
		img := DecodeFrameType6(frame, width, height, conf.Pal)
		imgs = append(imgs, img)
	}
	return imgs, nil
}
