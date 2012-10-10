package main

import "flag"
import "fmt"
import "log"
import "os"
import "path"

import "github.com/mewrnd/blizzconv/images/imgarchive"
import "github.com/mewrnd/blizzconv/images/imgconf"
import "github.com/mewrnd/blizzconv/mpq"

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
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]... [name.cel|name.cl2]...\n", os.Args[0])
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
