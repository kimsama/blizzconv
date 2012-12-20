// Package imgarchive implements support for extracting CEL and CL2 archives.
package imgarchive

import "fmt"
import "os"
import "path"
import "strings"

import "github.com/mewrnd/blizzconv/images/imgconf"
import "github.com/mewrnd/blizzconv/mpq"

// Extract extracts CEL and CL2 archives.
func Extract(archiveName string) (err error) {
	imageCount, found := imgconf.GetImageCount(archiveName)
	if !found {
		return fmt.Errorf("no archived images in %q.", archiveName)
	}
	archivePath, err := mpq.GetPath(archiveName)
	if err != nil {
		return err
	}
	fr, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer fr.Close()
	fws, err := createOutputImages(archivePath, imageCount)
	if err != nil {
		return err
	}
	defer closeFiles(fws)
	ext := path.Ext(archiveName)
	switch ext {
	case ".cel":
		err = ExtractCel(fr, fws)
	case ".cl2":
		err = ExtractCl2(fr, fws)
	default:
		return fmt.Errorf("imgarchive.Extract: unknown extension: %q.", ext)
	}
	if err != nil {
		return fmt.Errorf("imgarchive.Extract: error while extracting %q: %s.", archiveName, err)
	}
	return nil
}

// createOutputImages creates the output images of the archive. Note: remember
// to close the writers while done using them.
func createOutputImages(archivePath string, imageCount int) (fws []*os.File, err error) {
	posExt := strings.LastIndex(archivePath, ".")
	if posExt == -1 {
		return nil, fmt.Errorf("no extensions located for %q.", path.Base(archivePath))
	}
	for imageNum := 0; imageNum < imageCount; imageNum++ {
		imgPath := fmt.Sprintf("%s%d%s", archivePath[:posExt], imageNum, archivePath[posExt:])
		w, err := os.Create(imgPath)
		if err != nil {
			return nil, err
		}
		fws = append(fws, w)
	}
	return fws, nil
}

// closeFiles ranges through the file slice and closes each file.
func closeFiles(fws []*os.File) {
	for _, fw := range fws {
		fw.Close()
	}
}
