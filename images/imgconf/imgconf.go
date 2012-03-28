// Package imgconf implements functions for retrieving relevant CEL and CL2
// image information.
//
// This information is stored in an ini file, since neither the CEL nor the CL2
// image format contains the relevant image information.
package imgconf

import ini "github.com/glacjay/goini"

import "fmt"
import "sort"
import "strconv"
import "strings"

// MpqExtractPath is the path to an extracted MPQ file.
var MpqExtractPath string

// IniPath is the path to the 'cel.ini' or 'cl2.ini' file which provides CEL and
// CL2 image information.
var IniPath string

var dict ini.Dict

// Init loads an ini file which provides CEL and CL2 image information.
func Init() (err error) {
   dict, err = ini.Load(IniPath)
   if err != nil {
      return err
   }
   return nil
}

// AllFunc calls the function f with the parameter imgName once for each image
// in the ini file.
func AllFunc(f func(string) error) (err error) {
   var imgNames []string
   for imgName, _ := range dict {
      if imgName == "" {
         continue
      }
      imgNames = append(imgNames, imgName)
   }
   sort.Strings(imgNames)
   for _, imgName := range imgNames {
      err = f(imgName)
      if err != nil {
         return err
      }
   }
   return nil
}

// GetRelPath returns the relative path to the image.
func GetRelPath(imgName string) (relPath string, err error) {
   relPath, found := dict.GetString(imgName, "path")
   if !found {
      return "", fmt.Errorf("path not found for '%s'.", imgName)
   }
   return relPath, nil
}

// GetPath returns the full path to the image.
func GetPath(imgName string) (path string, err error) {
   path, err = GetRelPath(imgName)
   if err != nil {
      return "", err
   }
   return MpqExtractPath + path, nil
}

// GetWidth returns the image width.
func GetWidth(imgName string) (width int, err error) {
   width, found := dict.GetInt(imgName, "width")
   if !found {
      return 0, fmt.Errorf("width not found for '%s'.", imgName)
   }
   return width, nil
}

// GetHeight returns the image height.
func GetHeight(imgName string) (height int, err error) {
   height, found := dict.GetInt(imgName, "height")
   if !found {
      return 0, fmt.Errorf("height not found for '%s'.", imgName)
   }
   return height, nil
}

// GetRelPalPaths returns the relative paths to the image palettes.
func GetRelPalPaths(imgName string) (relPalPaths []string) {
   rawRelPalPaths, found := dict.GetString(imgName, "pals")
   if !found {
      // Default pal path:
      //    'levels/towndata/town.pal'
      return []string{"levels/towndata/town.pal"}
   }
   return strings.Split(rawRelPalPaths, ",")
}

// GetHeaderSize returns the header size of the image.
func GetHeaderSize(imgName string) (headerSize int) {
   headerSize, found := dict.GetInt(imgName, "header_size")
   if !found {
      return 0
   }
   return headerSize
}

// GetImageCount returns the number of archived images within the archive.
func GetImageCount(imgName string) (imageCount int, found bool) {
   imageCount, found = dict.GetInt(imgName, "image_count")
   if !found {
      return 0, false
   }
   return imageCount, true
}

// GetFrameWidth returns the width of the image's frames as a map from frameNum
// (key) to frameWidth (val).
func GetFrameWidth(imgName string) (frameWidth map[int]int, err error) {
   rawFrameWidths, found := dict.GetString(imgName, "frame_widths")
   if !found {
      return nil, nil
   }
   return getFrameDimension(rawFrameWidths)
}

// GetFrameHeight returns the height of the image's frames as a map from
// frameNum (key) to frameHeight (val).
func GetFrameHeight(imgName string) (frameHeight map[int]int, err error) {
   rawFrameHeights, found := dict.GetString(imgName, "frame_heights")
   if !found {
      return nil, nil
   }
   return getFrameDimension(rawFrameHeights)
}

// getFrameDimension parses frame widths and heights into a map from frameNum
// (key) to frameDimension (val). Below is an example frame_widths entry:
//    frame_widths=\
//          0:33,\
//        1-9:32,\
//         10:23,\
//      11-85:28,\
//     86-110:56
func getFrameDimension(rawFramesDimensions string) (frameDimension map[int]int, err error) {
   frameDimension = make(map[int]int)
   for _, rawFramesDimension := range strings.Split(rawFramesDimensions, ",") {
      rawFramesDimension = strings.TrimSpace(rawFramesDimension)
      posDelim := strings.LastIndex(rawFramesDimension, ":")
      if posDelim == -1 {
         return nil, fmt.Errorf("no delim ':' found for '%s'.", rawFramesDimension)
      }
      rawFrameNums := rawFramesDimension[:posDelim]
      rawDimension := rawFramesDimension[posDelim + 1:]
      dimension, err := strconv.Atoi(rawDimension)
      if err != nil {
         return nil, err
      }
      posDash := strings.LastIndex(rawFrameNums, "-")
      if posDash == -1 {
         frameNum, err := strconv.Atoi(rawFrameNums)
         if err != nil {
            return nil, err
         }
         frameDimension[frameNum] = dimension
      } else {
         frameNumStart, err := strconv.Atoi(rawFrameNums[:posDash])
         if err != nil {
            return nil, err
         }
         frameNumEnd, err := strconv.Atoi(rawFrameNums[posDash + 1:])
         if err != nil {
            return nil, err
         }
         for frameNum := frameNumStart; frameNum <= frameNumEnd; frameNum ++ {
            frameDimension[frameNum] = dimension
         }
      }
   }
   return frameDimension, nil
}
