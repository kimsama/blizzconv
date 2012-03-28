package main

import "github.com/mewkiz/blizzconv/images/cel"
///import "github.com/mewkiz/blizzconv/images/cl2"
import "github.com/mewkiz/blizzconv/images/imgarchive"
import "github.com/mewkiz/blizzconv/images/imgconf"

import "flag"
import "fmt"
import "image"
import "image/png"
import "log"
import "os"
import "path"
import "strings"

var flagAll bool

func init() {
   flag.Usage = usage
   flag.BoolVar(&flagAll, "a", false, "Dump all image files.")
   flag.StringVar(&imgconf.IniPath, "ini", "cel.ini", "Path to an ini file containing image information. Note: 'cl2.ini' will be used for files that have the '.cl2' extension.")
   flag.StringVar(&imgconf.MpqExtractPath, "mpqdump", "mpqdump/", "Path to an extracted MPQ file.")
   flag.Parse()
}

func usage() {
   fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [name.cel|name.cl2]...\n", os.Args[0])
   flag.PrintDefaults()
}

func main() {
   if flag.NArg() > 0 {
      if path.Ext(flag.Arg(0)) == ".cl2" {
         imgconf.IniPath = "cl2.ini"
      }
   }
   err := imgconf.Init()
   if err != nil {
      log.Fatalln(err)
   }
   if flagAll {
      // dump all images in the ini file.
      err := imgconf.AllFunc(Dump)
      if err != nil {
         log.Fatalln(err)
      }
      return
   }
   if flag.NArg() < 1 {
      flag.Usage()
      os.Exit(1)
   }
   for _, imgName := range flag.Args() {
      err := Dump(imgName)
      if err != nil {
         log.Fatalln(err)
      }
   }
}

// Dump extracts archived images if there are any, decodes image configs (pals)
// and dumps the image's frames, once for each image config.
func Dump(imgName string) (err error) {
   _, found := imgconf.GetImageCount(imgName)
   if found {
      // extract archived images
      err = imgarchive.Extract(imgName)
      if err != nil {
         return err
      }
      return nil
   }
   relPalPaths := imgconf.GetRelPalPaths(imgName)
   for palNum, relPalPath := range relPalPaths {
      /// ### todo ###
      ///   - add support for cl2
      ///   - maybe this should be handled with a wrapper library?
      /// ############
      conf, err := cel.GetConf(imgName, relPalPath)
      if err != nil {
         return err
      }
      var palDir string
      if len(relPalPaths) > 1 {
         palDir = fmt.Sprintf("pal_%04d/", palNum)
      }
      // dump the image's frames using conf (pal)
      err = dumpFrames(conf, palDir, imgName)
      if err != nil {
         return err
      }
   }
   return nil
}

// dumpFrames decodes an image's frames using a given image config (pal),
// creates a dump directory if there are more than one image and converts each
// frame to a new png image.
func dumpFrames(conf *cel.Config, palDir, imgName string) (err error) {
   // decode frames using the given image config (pal)
   imgs, err := cel.DecodeAll(imgName, conf)
   if err != nil {
      return err
   }
   // create dumpDir
   nameWithoutExt := imgName[:len(imgName) - len(path.Ext(imgName))]
   var frameDir, pngName string
   if len(imgs) > 1 {
      frameDir = nameWithoutExt + "/"
   } else {
      pngName = nameWithoutExt + ".png"
   }
   var dumpDir string
   /// ### todo ###
   ///   - should be len(imgs) > 1?
   /// ############
   if len(imgs) > 0 {
      dumpDir, err = createDumpDir(frameDir, palDir, imgName)
      if err != nil {
         return err
      }
   }
   for frameNum, img := range imgs {
      if len(imgs) > 1 {
         pngName = fmt.Sprintf("%s_%04d.png", nameWithoutExt, frameNum)
      }
      err := pngOutput(dumpDir + pngName, img)
      if err != nil {
         return err
      }
   }
   return nil
}

const dumpPrefix = "_dump_/"

// createDumpDir creates a dump directory for the image.
//
//    === [ dumpDir examples ] =================================================
//
//    --- [ one pal, one frame ] -----------------------------------------------
//
//       _dump_/imgDir/name.png
//
//    --- [ one pal, many frames ] ---------------------------------------------
//
//       _dump_/imgDir/name/name_0001.png
//       _dump_/imgDir/name/name_0002.png
//
//    --- [ many pals, one frame ] ---------------------------------------------
//
//       _dump_/imgDir/pal_0001/name.png
//       _dump_/imgDir/pal_0002/name.png
//
//    --- [ many pals, many frames ] -------------------------------------------
//
//       _dump_/imgDir/name/pal_0001/name_0001.png
//       _dump_/imgDir/name/pal_0001/name_0002.png
//       _dump_/imgDir/name/pal_0002/name_0001.png
//       _dump_/imgDir/name/pal_0002/name_0002.png
func createDumpDir(frameDir, palDir, imgName string) (dumpDir string, err error) {
   imgPath, err := imgconf.GetRelPath(imgName)
   if err != nil {
      return "", err
   }
   imgDir, _ := path.Split(imgPath)
   dumpDir = path.Clean(dumpPrefix + imgDir + frameDir + palDir) + "/"
   /**
    *    prevent directory traversal
    */
   if false == strings.HasPrefix(dumpDir, dumpPrefix) {
      return "", fmt.Errorf("path (%s) contains no dump prefix (%s).", dumpDir, dumpPrefix)
   }
   err = os.MkdirAll(dumpDir, 0755)
   if err != nil {
      return "", err
   }
   return dumpDir, nil
}

// pngOutput creates a new png image at pngPath from img.
func pngOutput(pngPath string, img image.Image) (err error) {
   f, err := os.Create(pngPath)
   if err != nil {
      return err
   }
   defer f.Close()
   err = png.Encode(f, img)
   if err != nil {
      return err
   }
   return nil
}
