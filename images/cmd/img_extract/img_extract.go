package main

import "github.com/mewkiz/blizzconv/images/imgarchive"
import "github.com/mewkiz/blizzconv/images/imgconf"

import "flag"
import "fmt"
import "log"
import "os"
import "path"

func init() {
   flag.Usage = usage
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
   for _, imgName := range flag.Args() {
      err = imgarchive.Extract(imgName)
      if err != nil {
         log.Fatalln(err)
      }
   }
}
