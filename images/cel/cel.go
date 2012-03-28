// Package cel implements a CEL image decoder.
//
// There are many similarities between CEL and GIF images. Both can contain
// multiple frames and use palettes. Below is a description of the CEL image
// format.
//
// CEL format:
//    frameCount     uint32                  // (little endian)
//    frameOffsets   [frameCount + 1]uint32  // (little endian)
//    frames         [frameCount][]byte      // The content of frameNum starts at frameOffsets[frameNum] and ends at frameOffsets[frameNum + 1].
//
// CEL frame format:
//    header   []byte   // Optional
//    data     []byte   // Frame pixel content. ref: DecodeType1
package cel

import "github.com/mewkiz/blizzconv/images/imgconf"

import dbg "fmt"
import "encoding/binary"
import "image"
import "image/color"
import "log"
import "os"

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
      var img image.Image
      switch celName {
      case "town.cel", "l1.cel", "l2.cel", "l3.cel", "l4.cel":
         log.Printf("decode not implemented for '%s'.", celName)
         return nil, nil
      default:
         width, ok := conf.FrameWidth[frameNum]
         if !ok {
            width = conf.Width
         }
         height, ok := conf.FrameHeight[frameNum]
         if !ok {
            height = conf.Height
         }
         img = DecodeType1(frame, width, height, conf.Pal)
      }
      imgs = append(imgs, img)
   }
   return imgs, nil
}

// GetFrames returns a slice of frames, whose content has been retrieved based
// on the CEL format described above.
func GetFrames(celName string) (frames [][]byte, err error) {
   celPath, err := imgconf.GetPath(celName)
   if err != nil {
      return nil, err
   }
   dbg.Println("cel:", celPath)
   f, err := os.Open(celPath)
   if err != nil {
      return nil, err
   }
   defer f.Close()
   var frameCount uint32
   err = binary.Read(f, binary.LittleEndian, &frameCount)
   if err != nil {
      return nil, err
   }
   dbg.Println("frame count:", frameCount)
   frameOffsets := make([]uint32, frameCount + 1)
   err = binary.Read(f, binary.LittleEndian, frameOffsets)
   if err != nil {
      return nil, err
   }
   for frameNum := uint32(0); frameNum < frameCount; frameNum++ {
      headerSize := imgconf.GetHeaderSize(celName)
      frameStart := int64(frameOffsets[frameNum]) + int64(headerSize)
      frameEnd := int64(frameOffsets[frameNum + 1])
      frameSize := frameEnd - frameStart
      frame := make([]byte, frameSize)
      _, err := f.ReadAt(frame, frameStart)
      if err != nil {
         return nil, err
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
   conf = &Config {
      Width:         width,
      Height:        height,
      Pal:           pal,
      FrameWidth:    frameWidth,
      FrameHeight:   frameHeight,
   }
   return conf, nil
}
