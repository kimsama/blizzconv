package cel

///import dbg "fmt"
import "image"
import "image/color"

// GetFrameDecoder returns the appropriate function for decoding the frame.
func GetFrameDecoder(celName string, frame []byte, frameNum int) func([]byte, int, int, color.Palette) image.Image {
	frameSize := len(frame)
	switch celName {
	case "l1.cel", "l2.cel", "l3.cel", "l4.cel", "town.cel":
		// Some regular (Type1) CEL images just happen to have the exact frame
		// size of 0x220, 0x320 and 0x400. Therefor the isType* functions are
		// required.
		switch frameSize {
		case 0x400:
			if isType0(celName, frameNum) {
				///dbg.Printf("type 0: plain 32x32        [X]   - %4d\n", frameNum)
				return DecodeFrameType0
			}
		case 0x220:
			if isType2or4(frame) {
				///dbg.Printf("type 2: triangle left      <|    - %4d\n", frameNum)
				return DecodeFrameType2
			} else if isType3or5(frame) {
				///dbg.Printf("type 3: triangle right      |>   - %4d\n", frameNum)
				return DecodeFrameType3
			}
		case 0x320:
			if isType2or4(frame) {
				///dbg.Printf("type 4: trapezoid left     /|    - %4d\n", frameNum)
				return DecodeFrameType4
			} else if isType3or5(frame) {
				///dbg.Printf("type 5: trapezoid right     |\\   - %4d\n", frameNum)
				return DecodeFrameType5
			}
		default:
			///dbg.Printf("type 1: regular cel        [ ]   - %4d\n", frameNum)
		}
	}
	return DecodeFrameType1
}

// isType0 returns true if the image is a plain 32x32.
func isType0(celName string, frameNum int) bool {
	// The following frames are Type1, thus return false.
	switch celName {
	case "l1.cel":
		switch frameNum {
		case 148, 159, 181, 186, 188:
			return false
		}
	case "l2.cel":
		switch frameNum {
		case 47, 1397, 1399, 1411:
			return false
		}
	case "l4.cel":
		switch frameNum {
		case 336, 639:
			return false
		}
	case "town.cel":
		switch frameNum {
		case 2328, 2367, 2593:
			return false
		}
	}
	return true
}

// isType2or4 returns true if the image is a triangle or a trapezoid pointing to
// the left.
//
// ref: DecodeFrameType2 and DecodeFrameType4
func isType2or4(frame []byte) bool {
	zeroPositions := []int{0, 1, 8, 9, 24, 25, 48, 49, 80, 81, 120, 121, 168, 169, 224, 225}
	for _, zeroPos := range zeroPositions {
		if frame[zeroPos] != 0x00 {
			return false
		}
	}
	return true
}

// isType3or5 returns true if the image is a triangle or a trapezoid pointing to
// the right.
//
// ref: DecodeFrameType3 and DecodeFrameType5
func isType3or5(frame []byte) bool {
	zeroPositions := []int{2, 3, 14, 15, 34, 35, 62, 63, 98, 99, 142, 143, 194, 195, 254, 255}
	for _, zeroPos := range zeroPositions {
		if frame[zeroPos] != 0x00 {
			return false
		}
	}
	return true
}

// DecodeFrameType1 returns an image after decoding the frame in the following
// way:
//
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
// Type1 is the most common type for CEL images.
func DecodeFrameType1(frame []byte, width int, height int, pal color.Palette) (img image.Image) {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	pos := 0
	for pos < len(frame) {
		chunkSize := int(int8(frame[pos]))
		pos++
		if chunkSize < 0 {
			// transparent pixels
			for i := 0; i > chunkSize; i-- {
				setPixel(rgba, color.RGBA{})
			}
		} else {
			// regular pixels
			for i := 0; i < chunkSize; i++ {
				setPixel(rgba, pal[frame[pos]])
				pos++
			}
		}
	}
	return rgba
}

// GetPixelSetter returns a function that can be invoced to incrementally set
// pixels, starting in the lower left corner, going from left to right, and then
// row by row from the bottom to the top of the image.
func GetPixelSetter(width, height int) func(*image.RGBA, color.Color) {
	var x, y int
	y = height - 1
	setPixel := func(rgba *image.RGBA, c color.Color) {
		rgba.Set(x, y, c)
		if x == width-1 {
			x = 0
			y--
		} else {
			x++
		}
	}
	return setPixel
}
