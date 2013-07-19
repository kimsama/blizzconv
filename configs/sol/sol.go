// Package sol implements functionality for parsing SOL files.
//
// SOL files contain information about various pillar properties, such as
// transparency and collision. Below is a description of the SOL format:
//
// SOL format:
//    // sol is a bitfield containing ###, ###, ###, ###, ###, ###, ### and ###:
//    //    ### := sol & 0x01
//    //    ### := sol & 0x02
//    //    ### := sol & 0x04 // block range (missiles and summoning of monsters).
//    //    ### := sol & 0x08 // allow transparency
//    //    ### := sol & 0x10
//    //    ### := sol & 0x20
//    //    ### := sol & 0x40
//    //    ### := sol & 0x80
//    solids []uint8
//
// The solid properties of a pillar can be obtained using the pillarNum as an
// offset into the solids array.
/// ### todo ###
///   - replace ### with something else.
/// ############
package sol

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/mewrnd/blizzconv/mpq"
)

// Solid defines the solid properties of a pillar.
type Solid struct {
	Sol0x01 bool
	Sol0x02 bool
	Sol0x04 bool
	Sol0x08 bool
	Sol0x10 bool
	Sol0x20 bool
	Sol0x40 bool
	Sol0x80 bool
}

// Parse parses a given SOL file and returns a slice of solids, based on the
// SOL format described above.
func Parse(solName string) (solids []Solid, err error) {
	solPath, err := mpq.GetPath(solName)
	if err != nil {
		return nil, err
	}
	fr, err := os.Open(solPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	var x uint8
	for {
		err = binary.Read(fr, binary.LittleEndian, &x)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		var solid Solid
		if x&0x01 != 0 {
			solid.Sol0x01 = true
		}
		if x&0x02 != 0 {
			solid.Sol0x02 = true
		}
		if x&0x04 != 0 {
			solid.Sol0x04 = true
		}
		if x&0x08 != 0 {
			solid.Sol0x08 = true
		}
		if x&0x10 != 0 {
			solid.Sol0x10 = true
		}
		if x&0x20 != 0 {
			solid.Sol0x20 = true
		}
		if x&0x40 != 0 {
			solid.Sol0x40 = true
		}
		if x&0x80 != 0 {
			solid.Sol0x80 = true
		}
		solids = append(solids, solid)
	}

	return solids, nil
}
