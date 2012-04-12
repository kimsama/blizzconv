package main

import "github.com/0xC3/progress/barcli"
import "github.com/mewkiz/blizzconv/configs/min"
import "github.com/mewkiz/blizzconv/images/cel"
import "github.com/mewkiz/blizzconv/images/imgconf"
import "github.com/mewkiz/blizzconv/mpq"
import "github.com/mewkiz/pkg/pngutil"

import dbg "fmt"
import "flag"
import "fmt"
import "image"
import "log"
import "os"
import "path"
import "strings"

func init() {
   flag.Usage = usage
   flag.StringVar(&imgconf.IniPath, "celini", "cel.ini", "Path to an ini file containing image information.")
   flag.StringVar(&mpq.ExtractPath, "mpqdump", "mpqdump/", "Path to an extracted MPQ file.")
   flag.StringVar(&mpq.IniPath, "mpqini", "mpq.ini", "Path to an ini file containing relative path information.")
   flag.Parse()
   err := mpq.Init()
   if err != nil {
      log.Fatalln(err)
   }
   err = imgconf.Init()
   if err != nil {
      log.Fatalln(err)
   }
}

func usage() {
   fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [name.min]...\n", os.Args[0])
   flag.PrintDefaults()
}

func main() {
   for _, minName := range flag.Args() {
      err := minDump(minName)
      if err != nil {
         log.Fatalln(err)
      }
   }
}

// bar represents the progress bar.
var bar *barcli.Bar

// dumpPrefix is the name of the dump directory.
const dumpPrefix = "_dump_/"

// minDump creates a dump directory and dumps the MIN file's pillars using the
// frames from a CEL image level file, once for each image config (pal).
func minDump(minName string) (err error) {
   pillars, err := min.Parse(minName)
   if err != nil {
      return err
   }
   nameWithoutExt := minName[:len(minName)-len(path.Ext(minName))]
   imgName := nameWithoutExt + ".cel"
   relPalPaths := imgconf.GetRelPalPaths(imgName)
   for _, relPalPath := range relPalPaths {
      conf, err := cel.GetConf(imgName, relPalPath)
      if err != nil {
         return err
      }
      var palDir string
      if len(relPalPaths) > 1 {
         dbg.Println("using pal:", relPalPath)
         palDir = path.Base(relPalPath) + "/"
      }
      bar, err = barcli.New(len(pillars))
      if err != nil {
         return err
      }
      levelFrames, err := cel.DecodeAll(imgName, conf)
      if err != nil {
         return err
      }
      dumpDir := path.Clean(dumpPrefix+"_pillars_/"+nameWithoutExt) + "/" + palDir
      // prevent directory traversal
      if !strings.HasPrefix(dumpDir, dumpPrefix) {
         return fmt.Errorf("path (%s) contains no dump prefix (%s).", dumpDir, dumpPrefix)
      }
      err = os.MkdirAll(dumpDir, 0755)
      if err != nil {
         return err
      }
      err = dumpPillars(pillars, levelFrames, dumpDir)
      if err != nil {
         return err
      }
   }
   return nil
}

// dumpPillars stores each pillar as a new png image, using the frames from a
// CEL image level file.
func dumpPillars(pillars []min.Pillar, levelFrames []image.Image, dumpDir string) (err error) {
   for pillarNum, pillar := range pillars {
      pillarPath := dumpDir + fmt.Sprintf("pillar_%04d.png", pillarNum)
      bar.Inc()
      img := pillar.Image(levelFrames)
      err = pngutil.WriteFile(pillarPath, img)
      if err != nil {
         return err
      }
   }
   return nil
}
