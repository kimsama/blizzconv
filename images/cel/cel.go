// Package cel implements a CEL image decoder.
//
// There are many similarities between CEL and GIF images. Both can contain
// multiple frames and use palettes. Below is a description of the CEL image
// format. All integers are stored in little endian.
//
// CEL format:
//    // frameCount specifies the number of frames contained within the image.
//    frameCount   uint32
//    // frameOffsets contains the offsets to each frame.
//    frameOffsets [frameCount + 1]uint32
//    // frames contains the header and data of each frame.
//    //    start: frameOffsets[frameNum]
//    //    end:   frameOffsets[frameNum + 1]
//    frames       [frameCount][]byte
//
// CEL frame format:
//    // header is optional.
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

// A Config specifies the frame dimensions and color palette of an image.
type Config struct {
	// The width of each frame in pixels.
	Width int
	// The height of each frame in pixels.
	Height int
	// A map from frameNum to frameWidth. It's used to override the default frame
	// width for specific frames.
	FrameWidth map[int]int
	// A map from frameNum to frameHeight. It's used to override the default
	// frame height for specific frames.
	FrameHeight map[int]int
	// The palette used for decoding.
	Pal color.Palette
}

// DecodeAll returns the sequential frames of a CEL image based on a given conf.
//
// Note: The absolute path of celName is resolved using mpq.GetPath.
func DecodeAll(celName string, conf *Config) (imgs []image.Image, err error) {
	// Get frame contents.
	frames, err := GetFrames(celName)
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
		decodeFrame := GetFrameDecoder(celName, frame, frameNum)
		img := decodeFrame(frame, width, height, conf.Pal)
		imgs = append(imgs, img)
	}

	return imgs, nil
}

// GetFrames returns a slice of frames, whose content has been retrieved based
// on the CEL format described above.
//
// Note: The absolute path of celName is resolved using mpq.GetPath.
func GetFrames(celName string) (frames [][]byte, err error) {
	// Open CEL file.
	celPath, err := mpq.GetPath(celName)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(celPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read frame count.
	var frameCount uint32
	err = binary.Read(f, binary.LittleEndian, &frameCount)
	if err != nil {
		return nil, fmt.Errorf("cel.GetFrames: unable to read frame count for %q: %v", celName, err)
	}

	// Read frame offsets.
	frameOffsets := make([]uint32, frameCount+1)
	err = binary.Read(f, binary.LittleEndian, frameOffsets)
	if err != nil {
		return nil, fmt.Errorf("cel.GetFrames: unable to read frame offsets for %q: %v", celName, err)
	}

	// Read frame contents.
	frames = make([][]byte, frameCount)
	for frameNum := range frames {
		// Ignore frame header.
		headerSize := imgconf.GetHeaderSize(celName)
		frameStart := int64(frameOffsets[frameNum]) + int64(headerSize)

		// Read frame content.
		frameEnd := int64(frameOffsets[frameNum+1])
		frameSize := frameEnd - frameStart
		frame := make([]byte, frameSize)
		_, err = f.ReadAt(frame, frameStart)
		if err != nil {
			return nil, fmt.Errorf("cel.GetFrames: unable to read frame content for %q: %v", celName, err)
		}
		frames[frameNum] = frame
	}

	return frames, nil
}

// GetConf returns a conf containing the relevant image information.
//
// Note: The absolute path of celName is resolved using mpq.GetPath and
// relPalPath is relative to mpq.ExtractPath.
func GetConf(celName, relPalPath string) (conf *Config, err error) {
	width, err := imgconf.GetWidth(celName)
	if err != nil {
		return nil, err
	}
	height, err := imgconf.GetHeight(celName)
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
	pal, err := GetPal(relPalPath)
	if err != nil {
		return nil, err
	}
	conf = &Config{
		Width:       width,
		Height:      height,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Pal:         pal,
	}
	return conf, nil
}
