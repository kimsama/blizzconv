// img_dump is a tool for converting CEL and CL2 images into png images.
//
// Usage:
//
//    img_dump [OPTION]... [name.cel|name.cl2]...
//
// Flags:
//
//    -a
//            Dump all image files.
//    -imgini="cel.ini"
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
	"image/color"
	"log"
	"os"
	"path"
	"strings"
    //win-fix
	//"github.com/0xC3/progress/barcli"
	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewrnd/blizzconv/images/cel"
	"github.com/mewrnd/blizzconv/images/cl2"
	"github.com/mewrnd/blizzconv/images/imgarchive"
	"github.com/mewrnd/blizzconv/images/imgconf"
	"github.com/mewrnd/blizzconv/images/trn"
	"github.com/mewrnd/blizzconv/mpq"
)

// flagAll specifies if all CEL images should be dumped or not.
var flagAll bool

func init() {
	flag.Usage = usage
	flag.BoolVar(&flagAll, "a", false, "Dump all image files.")
	flag.StringVar(&imgconf.IniPath, "imgini", "cel.ini", "Path to an ini file containing image information.")
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

// bar represents the progress bar.
//win-fix
//var bar *barcli.Bar

func main() {
	if flag.NArg() > 0 {
		if path.Ext(flag.Arg(0)) == ".cl2" && imgconf.IniPath == "cel.ini" {
			imgconf.IniPath = "cl2.ini"
		}
	}
	err := imgconf.Init()
	if err != nil {
		log.Fatalln(err)
	}
	if flagAll {
	    /* win-fix
		bar, err = barcli.New(imgconf.Len())
		if err != nil {
			log.Fatalln(err)
		}
		*/
		// dump all images in the ini file.
		err := imgconf.AllFunc(dump)
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
		err := dump(imgName)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// dump extracts archived images if there are any, decodes image configs (pals)
// and dumps the image's frames, once for each image config.
func dump(imgName string) (err error) {
	if flagAll {
		//win-fix
		//bar.Inc()
	}
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
	for _, relPalPath := range relPalPaths {
		conf, err := cel.GetConf(imgName, relPalPath)
		if err != nil {
			return err
		}
		var palDir string
		if len(relPalPaths) > 1 {
			palDir = path.Base(relPalPath) + "/"
		}

		// dump the image's frames using conf (pal) with no color transitions.
		err = dumpFrames(conf, palDir, "", imgName)
		if err != nil {
			return err
		}

		relTrnPaths := imgconf.GetRelTrnPaths(imgName)
		var srcPal color.Palette
		if len(relTrnPaths) > 0 {
			srcPal = make(color.Palette, len(conf.Pal))
			copy(srcPal, conf.Pal)
		}
		for _, relTrnPath := range relTrnPaths {
			conf.Pal, err = trn.ConvertPal(srcPal, relTrnPath)
			if err != nil {
				return err
			}

			var trnDir string
			if len(relTrnPaths) > 1 {
				trnDir = path.Base(relTrnPath) + "/"
			}

			// dump the image's frames using conf (pal) with color transitions.
			err = dumpFrames(conf, palDir, trnDir, imgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// dumpFrames decodes an image's frames using a given image config (pal),
// creates a dump directory and stores each frame as a new png image.
func dumpFrames(conf *cel.Config, palDir, trnDir, imgName string) (err error) {
	// decode frames using the given image config (pal)
	imgs, err := cl2.DecodeAll(imgName, conf)
	if err != nil {
		return err
	}
	// create dumpDir
	nameWithoutExt := imgName[:len(imgName)-len(path.Ext(imgName))]
	var frameDir, pngName string
	if len(imgs) > 1 {
		frameDir = nameWithoutExt + "/"
	} else {
		pngName = nameWithoutExt + ".png"
	}
	var dumpDir string
	if len(imgs) > 0 {
		dumpDir, err = createDumpDir(frameDir, palDir, trnDir, imgName)
		if err != nil {
			return err
		}
	}
	for frameNum, img := range imgs {
		if len(imgs) > 1 {
			pngName = fmt.Sprintf("%s_%04d.png", nameWithoutExt, frameNum)
		}
		err := imgutil.WriteFile(dumpDir+pngName, img)
		if err != nil {
			return err
		}
	}
	return nil
}

// dumpPrefix is the name of the dump directory.
const dumpPrefix = "_dump_/"

// createDumpDir creates a dump directory for the image.
//
//    === [ dumpDir examples ] =================================================
//
//    --- [ one pal, one trn, one frame ] --------------------------------------
//
//       _dump_/imgDir/name.png
//
//    --- [ one pal, one trn, many frames ] ------------------------------------
//
//       _dump_/imgDir/name/name_0001.png
//       _dump_/imgDir/name/name_0002.png
//
//    --- [ one pal, many trns, one frame ] ------------------------------------
//
//       _dump_/imgDir/trn_0001/name.png
//       _dump_/imgDir/trn_0002/name.png
//
//    --- [ one pal, many trns, many frames ] ----------------------------------
//
//       _dump_/imgDir/trn_0001/name_0001.png
//       _dump_/imgDir/trn_0001/name_0002.png
//       _dump_/imgDir/trn_0002/name_0001.png
//       _dump_/imgDir/trn_0002/name_0002.png
//
//    --- [ many pals, one trn, one frame ] ------------------------------------
//
//       _dump_/imgDir/pal_0001/name.png
//       _dump_/imgDir/pal_0002/name.png
//
//    --- [ many pals, one trn, many frames ] ----------------------------------
//
//       _dump_/imgDir/name/pal_0001/name_0001.png
//       _dump_/imgDir/name/pal_0001/name_0002.png
//       _dump_/imgDir/name/pal_0002/name_0001.png
//       _dump_/imgDir/name/pal_0002/name_0002.png
//
//    --- [ many pals, many trns, one frame ] ----------------------------------
//
//       _dump_/imgDir/pal_0001/trn_0001/name.png
//       _dump_/imgDir/pal_0001/trn_0002/name.png
//       _dump_/imgDir/pal_0002/trn_0001/name.png
//       _dump_/imgDir/pal_0002/trn_0002/name.png
//
//    --- [ many pals, many trns, many frames ] --------------------------------
//
//       _dump_/imgDir/name/pal_0001/trn_0001/name_0001.png
//       _dump_/imgDir/name/pal_0001/trn_0001/name_0002.png
//       _dump_/imgDir/name/pal_0001/trn_0002/name_0001.png
//       _dump_/imgDir/name/pal_0001/trn_0002/name_0002.png
//       _dump_/imgDir/name/pal_0002/trn_0001/name_0001.png
//       _dump_/imgDir/name/pal_0002/trn_0001/name_0002.png
//       _dump_/imgDir/name/pal_0002/trn_0002/name_0001.png
//       _dump_/imgDir/name/pal_0002/trn_0002/name_0002.png
func createDumpDir(frameDir, palDir, trnDir, imgName string) (dumpDir string, err error) {
	imgPath, err := mpq.GetRelPath(imgName)
	if err != nil {
		return "", err
	}
	imgDir, _ := path.Split(imgPath)
	dumpDir = path.Clean(dumpPrefix+imgDir+frameDir+palDir+trnDir) + "/"
	// prevent directory traversal
	if !strings.HasPrefix(dumpDir, dumpPrefix) {
		return "", fmt.Errorf("path (%s) contains no dump prefix (%s).", dumpDir, dumpPrefix)
	}
	err = os.MkdirAll(dumpDir, 0755)
	if err != nil {
		return "", err
	}
	return dumpDir, nil
}
