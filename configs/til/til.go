// Package til implements functionality for parsing TIL files.
//
// TIL files contain information about how to arrange the pillars, which are
// constructed based on the MIN format, in order to form a square. Below is a
// description of the TIL format:
//
// TIL format:
//    squares []Square
//
// Square format:
//    PillarNumTop    uint16
//    PillarNumRight  uint16
//    PillarNumLeft   uint16
//    PillarNumBottom uint16
//
// ref: Image (pillar arrangement illustration)
package til

import "encoding/binary"
import "io"
import "os"

import "github.com/mewkiz/blizzconv/mpq"

// Square is constructed of four pillars (top, right, left and bottom).
//
// ref: Image (pillar arrangement illustration)
type Square struct {
	PillarNumTop    int
	PillarNumRight  int
	PillarNumLeft   int
	PillarNumBottom int
}

// Parse parses a given TIL file and returns a slice of squares, based on the
// TIL format described above.
func Parse(tilName string) (squares []Square, err error) {
	tilPath, err := mpq.GetPath(tilName)
	if err != nil {
		return nil, err
	}
	fr, err := os.Open(tilPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	for {
		var x [4]uint16
		err = binary.Read(fr, binary.LittleEndian, &x)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		square := Square{
			PillarNumTop:    int(x[0]),
			PillarNumRight:  int(x[1]),
			PillarNumLeft:   int(x[2]),
			PillarNumBottom: int(x[3]),
		}
		squares = append(squares, square)
	}
	return squares, nil
}
