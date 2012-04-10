// Package min implements functionality for parsing MIN files.
//
// MIN files contain information about how to arrange the frames, in a CEL image
// level file, in order to form a pillar. Below is a description of the MIN
// format.
//
// MIN format:
//    pillars []Pillar
//
// Pillar format:
//    // blocks contains 10 blocks for l1.min, l2.min and l3.min and 16 blocks
//    // for l4.min and town.min.
//    //
//    // ref: BlockRect (block arrangement illustration)
//    blocks [blockCount]uint16
//
// Block format:
//    // block is a bitfield containing both frameNumPlus1 and Type:
//    //    frameNumPlus1 := block & 0x0FFF
//    //    Type          := block & 0x7000
//    block uint16
package min

import "github.com/mewkiz/blizzconv/mpq"

import "encoding/binary"
import "io"
import "os"

// Pillar contains 10 to 16 blocks, each corresponding to a frame in a CEL image
// level files.
//
// ref: BlockRect (block arrangement illustration)
type Pillar struct {
   Blocks []Block
}

// Block contains information about which CEL decode algorithm (Type) that
// should be used to decode a specific FrameNum in a CEL image level file.
type Block struct {
   IsValid  bool
   FrameNum int
   Type     int
}

// Parse parses a given MIN file and returns a slice of pillars, based on the
// MIN format described above.
func Parse(name string) (pillars []Pillar, err error) {
   path, err := mpq.GetPath(name)
   if err != nil {
      return nil, err
   }
   fr, err := os.Open(path)
   if err != nil {
      return nil, err
   }
   defer fr.Close()
   var blockCount int
   switch name {
   case "l1.min", "l2.min", "l3.min":
      blockCount = 10
   case "l4.min", "town.min":
      blockCount = 16
   }
   tmp := make([]uint16, blockCount)
   for {
      err = binary.Read(fr, binary.LittleEndian, &tmp)
      if err != nil {
         if err == io.EOF {
            break
         }
         return nil, err
      }
      pillar := Pillar{}
      pillar.Blocks = make([]Block, blockCount)
      for i := 0; i < blockCount; i++ {
         frameNumPlus1 := int(tmp[i] & 0x0FFF)
         if frameNumPlus1 != 0 {
            pillar.Blocks[i].IsValid = true
            pillar.Blocks[i].FrameNum = frameNumPlus1 - 1
         }
         pillar.Blocks[i].Type = int(tmp[i] & 0x7000) >> 12
      }
      pillars = append(pillars, pillar)
   }
   return pillars, nil
}
