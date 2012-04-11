package min

import "image"
import "image/draw"

// The width and height of a pillar block in pixels.
const (
   BlockWidth  = 32
   BlockHeight = 32
)

// PillarWidth is the width of a pillar in pixels.
const PillarWidth = BlockWidth * 2

// Width returns the width of the pillar in pixels.
func (pillar Pillar) Width() int {
   // the pillar is two blocks in width.
   return PillarWidth
}

// Height returns the height of the pillar in pixels.
func (pillar Pillar) Height() int {
   // the pillar is five (for l1.min, l2.min and l3.min) or eight (for l4.min
   // and town.min) blocks in height.
   return BlockHeight * len(pillar.Blocks) / 2
}

// Image returns an image constructed from the pillar's blocks.
//
// ref: BlockRect (block arrangement illustration)
func (pillar Pillar) Image(levelFrames []image.Image) (img image.Image) {
   dst := image.NewRGBA(image.Rect(0, 0, pillar.Width(), pillar.Height()))
   // draw blocks on the left side of the pillar.
   blockNumStartLeft := len(pillar.Blocks) - 2
   pillar.drawSide(dst, levelFrames, blockNumStartLeft)
   // draw blocks on the right side of the pillar.
   blockNumStartRight := len(pillar.Blocks) - 1
   pillar.drawSide(dst, levelFrames, blockNumStartRight)
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
   0:  image.Rect(0, BlockHeight*0, BlockWidth, BlockHeight*1),
   2:  image.Rect(0, BlockHeight*1, BlockWidth, BlockHeight*2),
   4:  image.Rect(0, BlockHeight*2, BlockWidth, BlockHeight*3),
   6:  image.Rect(0, BlockHeight*3, BlockWidth, BlockHeight*4),
   8:  image.Rect(0, BlockHeight*4, BlockWidth, BlockHeight*5),
   10: image.Rect(0, BlockHeight*5, BlockWidth, BlockHeight*6),
   12: image.Rect(0, BlockHeight*6, BlockWidth, BlockHeight*7),
   14: image.Rect(0, BlockHeight*7, BlockWidth, BlockHeight*8),
   // odd blockNum
   1:  image.Rect(BlockWidth, BlockHeight*0, BlockWidth*2, BlockHeight*1),
   3:  image.Rect(BlockWidth, BlockHeight*1, BlockWidth*2, BlockHeight*2),
   5:  image.Rect(BlockWidth, BlockHeight*2, BlockWidth*2, BlockHeight*3),
   7:  image.Rect(BlockWidth, BlockHeight*3, BlockWidth*2, BlockHeight*4),
   9:  image.Rect(BlockWidth, BlockHeight*4, BlockWidth*2, BlockHeight*5),
   11: image.Rect(BlockWidth, BlockHeight*5, BlockWidth*2, BlockHeight*6),
   13: image.Rect(BlockWidth, BlockHeight*6, BlockWidth*2, BlockHeight*7),
   15: image.Rect(BlockWidth, BlockHeight*7, BlockWidth*2, BlockHeight*8),
}

// drawSide draws each block on one side of the pillar, starting from the bottom
// and going to top.
func (pillar Pillar) drawSide(dst draw.Image, levelFrames []image.Image, blockNumStart int) {
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
