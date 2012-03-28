/*
 *    image/cel
 */

package cel

import "image"
import "image/color"

// DecodeType1 returns an image after decoding the frame in the following way:
//    1) Read one byte (chunkSize).
//    2) If chunkSize is negative, set that many transparent pixels.
//    3) If chunkSize is positive, read that many bytes.
//       - Each byte read this way corresponds to a color index of the palette.
//       - Set one regular pixel per byte, using the color index to locate the
//         color in the palette.
//    4) goto 1 until EOF is reached.
//
// Pixels are stored "upside-down" with respect to normal image raster scan
// order, starting in the lower left corner, going from left to right, and then
// row by row from the bottom to the top of the image.
//
// Coordinate system:
//
//     [ y ]
//
//         +---+---+
//       1 |   |   |
//         +---+---+
//       0 |   |   |
//         +---+---+
//           0   1      [ x ]
//
func DecodeType1(frame []byte, width int, height int, pal color.Palette) (img image.Image) {
   rgba := image.NewRGBA(image.Rect(0, 0, width, height))
   pos := 0
   var x, y int
   y = height - 1
   set := func(c color.Color) {
      if x == width - 1 {
         x = 0
         y--
      } else {
         x++
      }
      rgba.Set(x, y, c)
   }
   for pos < len(frame) {
      chunkSize := int(int8(frame[pos]))
      pos++
      if chunkSize < 0 {
         // transparent pixels
         for i := 0; i > chunkSize; i-- {
            set(color.RGBA{})
         }
      } else {
         // regular pixels
         for i := 0; i < chunkSize; i++ {
            set(pal[frame[pos]])
            pos++
         }
      }
   }
   return rgba
}
