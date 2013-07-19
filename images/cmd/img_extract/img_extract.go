// img_extract is a tool for extracting CEL and CL2 archives.
//
// Usage:
//
//    img_extract [OPTION]... [name.cel|name.cl2]...
//
// Flags:
//
//    -celini="cel.ini"
//            Path to an ini file containing image information.
//            Note: 'cl2.ini' will be used for files that have the '.cl2' extension.
//    -mpqdump="mpqdump/"
//            Path to an extracted MPQ file.
//    -mpqini="mpq.ini"
//            Path to an ini file containing relative path information.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mewrnd/blizzconv/images/imgarchive"
	"github.com/mewrnd/blizzconv/images/imgconf"
	"github.com/mewrnd/blizzconv/mpq"
)

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
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [name.cel|name.cl2]...\n", os.Args[0])
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
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
