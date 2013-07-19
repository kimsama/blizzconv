// Package cel implements a CEL image decoder.
//
// There are many similarities between CEL and GIF images. Both can contain
// multiple frames and use palettes. Below is a description of the CEL image
// format.
//
// CEL format:
//    // (little endian)
//    frameCount   uint32
//    // frameOffsets contains the offsets to each frame. (little endian)
//    frameOffsets [frameCount + 1]uint32
//    // frames contains the header and data of each frame.
//    //    start: frameOffsets[frameNum]
//    //    end:   frameOffsets[frameNum + 1]
//    frames       [frameCount][]byte
//
// CEL frame format:
//    // header is optional
//    header []byte
//    // data contains the frame pixel content.
//    //
//    // ref: DecodeFrameType1
//    data   []byte
package cel

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/mewrnd/blizzconv/images/imgconf"
	"github.com/mewrnd/blizzconv/mpq"
)

// Config holds an image's palette and dimensions.
type Config struct {
	Width       int
	Height      int
	FrameWidth  map[int]int
	FrameHeight map[int]int
	Pal         color.Palette
}

// DecodeAll returns the sequential frames of a CEL image based on a given conf.
func DecodeAll(celName string, conf *Config) (imgs []image.Image, err error) {
	frames, err := GetFrames(celName)
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
		decodeFrame := GetFrameDecoder(celName, frame, frameNum)
		img := decodeFrame(frame, width, height, conf.Pal)
		imgs = append(imgs, img)
	}
	return imgs, nil
}

// GetFrames returns a slice of frames, whose content has been retrieved based
// on the CEL format described above.
func GetFrames(celName string) (frames [][]byte, err error) {
	celPath, err := mpq.GetPath(celName)
	if err != nil {
		return nil, err
	}
	fr, err := os.Open(celPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	var frameCount uint32
	err = binary.Read(fr, binary.LittleEndian, &frameCount)
	if err != nil {
		return nil, fmt.Errorf("cel.GetFrames: error while reading frame count for %q: %s.", celName, err)
	}
	frameOffsets := make([]uint32, frameCount+1)
	err = binary.Read(fr, binary.LittleEndian, frameOffsets)
	if err != nil {
		return nil, fmt.Errorf("cel.GetFrames: error while reading frame offsets for %q: %s.", celName, err)
	}
	for frameNum := uint32(0); frameNum < frameCount; frameNum++ {
		headerSize := imgconf.GetHeaderSize(celName)
		frameStart := int64(frameOffsets[frameNum]) + int64(headerSize)
		frameEnd := int64(frameOffsets[frameNum+1])
		frameSize := frameEnd - frameStart
		frame := make([]byte, frameSize)
		_, err := fr.ReadAt(frame, frameStart)
		if err != nil {
			return nil, fmt.Errorf("cel.GetFrames: error while reading frame content for %q: %s.", celName, err)
		}
		frames = append(frames, frame)
	}
	return frames, nil
}

// GetConf returns a conf containing the relevant image information.
func GetConf(celName, relPalPath string) (conf *Config, err error) {
	width, err := imgconf.GetWidth(celName)
	if err != nil {
		return nil, err
	}
	height, err := imgconf.GetHeight(celName)
	if err != nil {
		return nil, err
	}
	pal, err := GetPal(relPalPath)
	if err != nil {
		return nil, err
	}
	frameWidth, err := imgconf.GetFrameWidth(celName)
	if err != nil {
		return nil, err
	}
	frameHeight, err := imgconf.GetFrameHeight(celName)
	if err != nil {
		return nil, err
	}
	conf = &Config{
		Width:       width,
		Height:      height,
		Pal:         pal,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
	}
	return conf, nil
}
