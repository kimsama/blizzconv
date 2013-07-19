package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mewrnd/blizzconv/configs/sol"
	"github.com/mewrnd/blizzconv/mpq"
)

func init() {
	flag.Usage = usage
	flag.StringVar(&mpq.ExtractPath, "mpqdump", "mpqdump/", "Path to an extracted MPQ file.")
	flag.StringVar(&mpq.IniPath, "mpqini", "mpq.ini", "Path to an ini file containing relative path information.")
	flag.Parse()
	err := mpq.Init()
	if err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [name.sol]...\n", os.Args[0])
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	for _, solName := range flag.Args() {
		err := solDump(solName)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func solDump(solName string) (err error) {
	solids, err := sol.Parse(solName)
	if err != nil {
		return err
	}
	for pillarNum, solid := range solids {
		fmt.Println("pillarNum:", pillarNum)
		if solid.Sol0x01 {
			fmt.Println("   0x01:", solid.Sol0x01)
		}
		if solid.Sol0x02 {
			fmt.Println("   0x02:", solid.Sol0x02)
		}
		if solid.Sol0x04 {
			fmt.Println("   0x04:", solid.Sol0x04)
		}
		if solid.Sol0x08 {
			fmt.Println("   0x08:", solid.Sol0x08)
		}
		if solid.Sol0x10 {
			fmt.Println("   0x10:", solid.Sol0x10)
		}
		if solid.Sol0x20 {
			fmt.Println("   0x20:", solid.Sol0x20)
		}
		if solid.Sol0x40 {
			fmt.Println("   0x40:", solid.Sol0x40)
		}
		if solid.Sol0x80 {
			fmt.Println("   0x80:", solid.Sol0x80)
		}
	}
	return nil
}
