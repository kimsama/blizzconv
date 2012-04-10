package min

import "image"
import "image/draw"

// The width and height of a pillar block in pixels.
const (
   BlockWidth  = 32
   BlockHeight = 32
)

// PillarImage returns an image constructed from the pillar's blocks.
//
// ref: BlockRect (block arrangement illustration)
func PillarImage(levelFrames []image.Image, pillar Pillar) (img image.Image) {
   // the pillar is two blocks in width.
   width := BlockWidth * 2
   // the pillar is five (for l1.min, l2.min and l3.min) or eight (for l4.min
   // and town.min) blocks in height.
   height := BlockHeight * len(pillar.Blocks) / 2
   dst := image.NewRGBA(image.Rect(0, 0, width, height))
   // draw blocks on the left side of the pillar.
   blockNumStartLeft := len(pillar.Blocks) - 2
   drawSide(dst, levelFrames, pillar, blockNumStartLeft)
   // draw blocks on the right side of the pillar.
   blockNumStartRight := len(pillar.Blocks) - 1
   drawSide(dst, levelFrames, pillar, blockNumStartRight)
   return dst
}

// BlockRect is a map from blockNum to an image.Rectangle of the block.
//
// The size of each pillar block is 32x32 pixels. The blocks are arranged as
// illustrated below, forming a pillar:
//
//    +----+----+
//    |  0 |  1 |
//    +----+----+
//    |  2 |  3 |
//    +----+----+
//    |  4 |  5 |
//    +----+----+
//    |  6 |  7 |
//    +----+----+
//    |  8 |  9 |
//    +----+----+
//    | 10 | 11 |
//    +----+----+
//    | 12 | 13 |
//    +----+----+
//    | 14 | 15 |
//    +----+----+
var BlockRect = map[int]image.Rectangle{
   // even blockNum
   0:  image.Rect(0, 32*0, 32, 32*1),
   2:  image.Rect(0, 32*1, 32, 32*2),
   4:  image.Rect(0, 32*2, 32, 32*3),
   6:  image.Rect(0, 32*3, 32, 32*4),
   8:  image.Rect(0, 32*4, 32, 32*5),
   10: image.Rect(0, 32*5, 32, 32*6),
   12: image.Rect(0, 32*6, 32, 32*7),
   14: image.Rect(0, 32*7, 32, 32*8),
   // odd blockNum
   1:  image.Rect(32, 32*0, 64, 32*1),
   3:  image.Rect(32, 32*1, 64, 32*2),
   5:  image.Rect(32, 32*2, 64, 32*3),
   7:  image.Rect(32, 32*3, 64, 32*4),
   9:  image.Rect(32, 32*4, 64, 32*5),
   11: image.Rect(32, 32*5, 64, 32*6),
   13: image.Rect(32, 32*6, 64, 32*7),
   15: image.Rect(32, 32*7, 64, 32*8),
}

// drawSide draws each block on one side of the pillar, starting from the bottom
// and going to top.
func drawSide(dst draw.Image, levelFrames []image.Image, pillar Pillar, blockNumStart int) {
   var moveUp, first bool
   first = true
   for blockNum := blockNumStart; blockNum >= 0; blockNum -= 2 {
      block := pillar.Blocks[blockNum]
      if block.IsValid {
         if first {
            switch block.Type {
            // if the first block in a section is type 1, 4 or 5 the entire
            // section of blocks should move up.
            case 1, 4, 5:
               moveUp = true
            default:
               moveUp = false
            }
            first = false
         }
         rect := BlockRect[blockNum]
         if moveUp {
            rect.Min.Y--
            rect.Max.Y--
         }
         draw.Draw(dst, rect, levelFrames[block.FrameNum], image.ZP, draw.Src)
      } else {
         // if the entire block is transparent, start a new section.
         first = true
      }
   }
}
