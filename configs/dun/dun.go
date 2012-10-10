// Package dun implements functionality for parsing DUN files.
//
// DUN files contain information about how to arrange the squares, which are
// constructed based on the TIL format, in order to form a dungeon. Below is a
// description of the DUN format:
//
// DUN format:
//    dunWidth        uint16
//    dunHeight       uint16
//    squareNumsPlus1 [dunWidth][dunHeight]uint16
package dun

import "encoding/binary"
import "fmt"
import "os"
import "path"

import "github.com/mewrnd/blizzconv/configs/dunconf"
import "github.com/mewrnd/blizzconv/configs/til"
import "github.com/mewrnd/blizzconv/mpq"

// The maximum number of cols and rows in a dungeon map.
const (
	ColMax = 112
	RowMax = 112
)

// Dungeon maps from a col and a row to a pillarNum. Each pillarNum value is
// initialized to -1, which corresponds to a transparent pillar.
type Dungeon [ColMax][RowMax]int

// New returns a new Dungeon, where all pillarNum values have been initialized
// to -1.
func New() (dungeon *Dungeon) {
	dungeon = new(Dungeon)
	for row := 0; row < RowMax; row++ {
		for col := 0; col < ColMax; col++ {
			dungeon[col][row] = -1
		}
	}
	return dungeon
}

// Parse parses a given DUN file and stores each pillarNum at a coordinate in
// the dungeon, based on the DUN format described above.
//
// Below is a description of how the squares are positioned on the dungeon map:
//    1) Start at the coordinates colStart, rowStart.
//    2) Place a square.
//       - Each square is two cols in width and two rows in height.
//    3) Increment col with two.
//    4) goto 2) dunWidth number of times.
//    5) Increment row with two.
//    6) goto 2) dunHeight number of times.
//
// ref: GetPillarRect (illustration of map coordinate system)
func (dungeon *Dungeon) Parse(dunName string) (err error) {
	dunPath, err := mpq.GetPath(dunName)
	if err != nil {
		return err
	}
	fr, err := os.Open(dunPath)
	if err != nil {
		return err
	}
	defer fr.Close()
	var tmp [2]uint16
	err = binary.Read(fr, binary.LittleEndian, &tmp)
	if err != nil {
		return err
	}
	dunWidth := int(tmp[0])
	dunHeight := int(tmp[1])
	colStart, err := dunconf.GetColStart(dunName)
	if err != nil {
		return err
	}
	rowStart, err := dunconf.GetRowStart(dunName)
	if err != nil {
		return err
	}
	nameWithoutExt, err := GetLevelName(dunName)
	if err != nil {
		return err
	}
	squares, err := til.Parse(nameWithoutExt + ".til")
	if err != nil {
		return err
	}
	row := rowStart
	for i := 0; i < dunHeight; i++ {
		col := colStart
		for j := 0; j < dunWidth; j++ {
			var x uint16
			err = binary.Read(fr, binary.LittleEndian, &x)
			if err != nil {
				return err
			}
			squareNumPlus1 := int(x)
			square := til.Square{
				PillarNumTop:    -1,
				PillarNumRight:  -1,
				PillarNumLeft:   -1,
				PillarNumBottom: -1,
			}
			if squareNumPlus1 != 0 {
				square = squares[squareNumPlus1-1]
			}
			dungeon[col][row] = square.PillarNumTop
			dungeon[col+1][row] = square.PillarNumRight
			dungeon[col][row+1] = square.PillarNumLeft
			dungeon[col+1][row+1] = square.PillarNumBottom
			col += 2
		}
		row += 2
	}
	return nil
}

// GetLevelName returns the level name (without extension) of a given DUN file.
func GetLevelName(dunName string) (nameWithoutExt string, err error) {
	relDunPath, err := mpq.GetRelPath(dunName)
	if err != nil {
		return "", err
	}
	dunDir, _ := path.Split(relDunPath)
	switch dunDir {
	case "levels/l1data/":
		nameWithoutExt = "l1"
	case "levels/l2data/":
		nameWithoutExt = "l2"
	case "levels/l3data/":
		nameWithoutExt = "l3"
	case "levels/l4data/":
		nameWithoutExt = "l4"
	case "levels/towndata/":
		nameWithoutExt = "town"
	default:
		return "", fmt.Errorf("invalid dunDir (%s).", dunDir)
	}
	return nameWithoutExt, nil
}
