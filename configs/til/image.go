package til

import "image"
import "image/draw"

import "github.com/mewkiz/blizzconv/configs/min"

// Image returns an image constructed from the square's pillars. The pillars are
// arranged as illustrated below, forming a square:
//
//           top
//
//            /\
//    left   /\/\   right
//           \/\/
//            \/
//
//          bottom
func (square Square) Image(pillars []min.Pillar, levelFrames []image.Image) (img image.Image) {
   // the square is two pillars in width.
   width := min.PillarWidth * 2
   // the square is one pillar and one block in height.
   height := pillars[0].Height() + min.BlockHeight
   dst := image.NewRGBA(image.Rect(0, 0, width, height))
   imgTop := pillars[square.PillarNumTop].Image(levelFrames)
   imgRight := pillars[square.PillarNumRight].Image(levelFrames)
   imgLeft := pillars[square.PillarNumLeft].Image(levelFrames)
   imgBottom := pillars[square.PillarNumBottom].Image(levelFrames)
   pointTop := image.Pt(min.PillarWidth/2, 0)
   pointRight := image.Pt(min.PillarWidth, min.BlockHeight/2)
   pointLeft := image.Pt(0, min.BlockHeight/2)
   pointBottom := image.Pt(min.PillarWidth/2, min.BlockHeight)
   bounds := imgTop.Bounds()
   draw.Draw(dst, bounds.Add(pointTop), imgTop, image.ZP, draw.Over)
   draw.Draw(dst, bounds.Add(pointRight), imgRight, image.ZP, draw.Over)
   draw.Draw(dst, bounds.Add(pointLeft), imgLeft, image.ZP, draw.Over)
   draw.Draw(dst, bounds.Add(pointBottom), imgBottom, image.ZP, draw.Over)
   return dst
}
