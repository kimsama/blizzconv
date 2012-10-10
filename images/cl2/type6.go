package cl2

import "image"
import "image/color"

import "github.com/mewrnd/blizzconv/images/cel"

// DecodeFrameType6 returns an image after decoding the frame in the following
// way:
//
//    1) Read one byte (chunkSize).
//    2) If chunkSize is positive, set that many transparent pixels.
//    3) If chunkSize is negative, invert it's sign.
//       3a) If chunkSize is below or equal to 65, read that many bytes.
//          - Each byte read this way corresponds to a color index of the
//            palette.
//          - Set one regular pixel per byte, using the color index to locate
//            the color in the palette.
//       3b) If chunkSize is above 65, subtract 65 from it and read one byte.
//          - The byte read this way corresponds to a color index of the
//            palette.
//          - Set chunkSize regular pixels, using the color index to locate the
//            color in the palette.
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
// Type6 is the only type for CL2 images.
func DecodeFrameType6(frame []byte, width int, height int, pal color.Palette) (img image.Image) {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := cel.GetPixelSetter(width, height)
	pos := 0
	for pos < len(frame) {
		chunkSize := int(int8(frame[pos]))
		pos++
		if chunkSize >= 0 {
			// transparent pixels
			for i := 0; i < chunkSize; i++ {
				setPixel(rgba, color.RGBA{})
			}
		} else {
			chunkSize = -chunkSize
			if chunkSize <= 65 {
				// regular pixels
				for i := 0; i < chunkSize; i++ {
					setPixel(rgba, pal[frame[pos]])
					pos++
				}
			} else {
				chunkSize -= 65
				// RLE encoded pixels
				for i := 0; i < chunkSize; i++ {
					setPixel(rgba, pal[frame[pos]])
				}
				pos++
			}
		}
	}
	return rgba
}
