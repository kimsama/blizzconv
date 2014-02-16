package cel

import (
	"image"
	"image/color"
	"image/draw"
)

// DecodeFrameType0 returns an image after decoding the frame in the following
// way:
//
//    1) Range through the frame, one byte at the time.
//       - Each byte corresponds to a color index of the palette.
//       - Set one regular pixel per byte, using the color index to locate the
//         color in the palette.
//
// Type0 corresponds to a plain 32x32 images, with no transparency.
func DecodeFrameType0(frame []byte, width int, height int, pal color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	for _, b := range frame {
		setPixel(img, pal[b])
	}
	return img
}

// DecodeFrameType2 returns an image after decoding the frame in the following
// way:
//
//    1) Dump one line of 32 pixels at the time.
//       - The illustration below tells if a pixel is transparent or regular.
//       - Only regular and zero (transparent) pixels are explicitly stored in
//         the frame content. The transparent pixels are implicitly referred
//         from the illustration.
//
// Below is an illustration of the 32x32 image, where a space represents an
// implicit transparent pixel, a '0' represents an explicit transparent pixel
// and an 'x' represents an explicit regular pixel.
//
// Note: The output image will be "upside-down" compared to the illustration.
//
//    +--------------------------------+
//    |                                |
//    |                            00xx|
//    |                            xxxx|
//    |                        00xxxxxx|
//    |                        xxxxxxxx|
//    |                    00xxxxxxxxxx|
//    |                    xxxxxxxxxxxx|
//    |                00xxxxxxxxxxxxxx|
//    |                xxxxxxxxxxxxxxxx|
//    |            00xxxxxxxxxxxxxxxxxx|
//    |            xxxxxxxxxxxxxxxxxxxx|
//    |        00xxxxxxxxxxxxxxxxxxxxxx|
//    |        xxxxxxxxxxxxxxxxxxxxxxxx|
//    |    00xxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |    xxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |00xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |00xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |    xxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |    00xxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |        xxxxxxxxxxxxxxxxxxxxxxxx|
//    |        00xxxxxxxxxxxxxxxxxxxxxx|
//    |            xxxxxxxxxxxxxxxxxxxx|
//    |            00xxxxxxxxxxxxxxxxxx|
//    |                xxxxxxxxxxxxxxxx|
//    |                00xxxxxxxxxxxxxx|
//    |                    xxxxxxxxxxxx|
//    |                    00xxxxxxxxxx|
//    |                        xxxxxxxx|
//    |                        00xxxxxx|
//    |                            xxxx|
//    |                            00xx|
//    +--------------------------------+
//
// Type2 corresponds to a 32x32 images of a left facing triangle.
func DecodeFrameType2(frame []byte, width int, height int, pal color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	decodeCounts := []int{0, 4, 4, 8, 8, 12, 12, 16, 16, 20, 20, 24, 24, 28, 28, 32, 32, 32, 28, 28, 24, 24, 20, 20, 16, 16, 12, 12, 8, 8, 4, 4}
	for lineNum, decodeCount := range decodeCounts {
		var zeroCount int
		switch lineNum {
		case 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31:
			zeroCount = 2
		}
		regularCount := decodeCount - zeroCount
		decodeLineTransparencyLeft(img, setPixel, frame, regularCount, zeroCount, pal)
		frame = frame[decodeCount:]
	}
	return img
}

// DecodeFrameType3 returns an image after decoding the frame in the following
// way:
//
//    1) Dump one line of 32 pixels at the time.
//       - The illustration below tells if a pixel is transparent or regular.
//       - Only regular and zero (transparent) pixels are explicitly stored in
//         the frame content. The transparent pixels are implicitly referred
//         from the illustration.
//
// Below is an illustration of the 32x32 image, where a space represents an
// implicit transparent pixel, a '0' represents an explicit transparent pixel
// and an 'x' represents an explicit regular pixel.
//
// Note: The output image will be "upside-down" compared to the illustration.
//
//    +--------------------------------+
//    |                                |
//    |xx00                            |
//    |xxxx                            |
//    |xxxxxx00                        |
//    |xxxxxxxx                        |
//    |xxxxxxxxxx00                    |
//    |xxxxxxxxxxxx                    |
//    |xxxxxxxxxxxxxx00                |
//    |xxxxxxxxxxxxxxxx                |
//    |xxxxxxxxxxxxxxxxxx00            |
//    |xxxxxxxxxxxxxxxxxxxx            |
//    |xxxxxxxxxxxxxxxxxxxxxx00        |
//    |xxxxxxxxxxxxxxxxxxxxxxxx        |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxx00    |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxx    |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx00|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx00|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxx    |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxx00    |
//    |xxxxxxxxxxxxxxxxxxxxxxxx        |
//    |xxxxxxxxxxxxxxxxxxxxxx00        |
//    |xxxxxxxxxxxxxxxxxxxx            |
//    |xxxxxxxxxxxxxxxxxx00            |
//    |xxxxxxxxxxxxxxxx                |
//    |xxxxxxxxxxxxxx00                |
//    |xxxxxxxxxxxx                    |
//    |xxxxxxxxxx00                    |
//    |xxxxxxxx                        |
//    |xxxxxx00                        |
//    |xxxx                            |
//    |xx00                            |
//    +--------------------------------+
//
// Type3 corresponds to a 32x32 images of a right facing triangle.
func DecodeFrameType3(frame []byte, width int, height int, pal color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	decodeCounts := []int{0, 4, 4, 8, 8, 12, 12, 16, 16, 20, 20, 24, 24, 28, 28, 32, 32, 32, 28, 28, 24, 24, 20, 20, 16, 16, 12, 12, 8, 8, 4, 4}
	for lineNum, decodeCount := range decodeCounts {
		var zeroCount int
		switch lineNum {
		case 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31:
			zeroCount = 2
		}
		regularCount := decodeCount - zeroCount
		decodeLineTransparencyRight(img, setPixel, frame, regularCount, zeroCount, pal)
		frame = frame[decodeCount:]
	}
	return img
}

// DecodeFrameType4 returns an image after decoding the frame in the following
// way:
//
//    1) Dump one line of 32 pixels at the time.
//       - The illustration below tells if a pixel is transparent or regular.
//       - Only regular and zero (transparent) pixels are explicitly stored in
//         the frame content. The transparent pixels are implicitly referred
//         from the illustration.
//
// Below is an illustration of the 32x32 image, where a space represents an
// implicit transparent pixel, a '0' represents an explicit transparent pixel
// and an 'x' represents an explicit regular pixel.
//
// Note: The output image will be "upside-down" compared to the illustration.
//
//    +--------------------------------+
//    |                            00xx|
//    |                            xxxx|
//    |                        00xxxxxx|
//    |                        xxxxxxxx|
//    |                    00xxxxxxxxxx|
//    |                    xxxxxxxxxxxx|
//    |                00xxxxxxxxxxxxxx|
//    |                xxxxxxxxxxxxxxxx|
//    |            00xxxxxxxxxxxxxxxxxx|
//    |            xxxxxxxxxxxxxxxxxxxx|
//    |        00xxxxxxxxxxxxxxxxxxxxxx|
//    |        xxxxxxxxxxxxxxxxxxxxxxxx|
//    |    00xxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |    xxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |00xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    +--------------------------------+
//
// Type4 corresponds to a 32x32 images of a left facing trapezoid.
func DecodeFrameType4(frame []byte, width int, height int, pal color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	decodeCounts := []int{4, 4, 8, 8, 12, 12, 16, 16, 20, 20, 24, 24, 28, 28, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	for lineNum, decodeCount := range decodeCounts {
		var zeroCount int
		switch lineNum {
		case 0, 2, 4, 6, 8, 10, 12, 14:
			zeroCount = 2
		}
		regularCount := decodeCount - zeroCount
		decodeLineTransparencyLeft(img, setPixel, frame, regularCount, zeroCount, pal)
		frame = frame[decodeCount:]
	}
	return img
}

// DecodeFrameType5 returns an image after decoding the frame in the following
// way:
//
//    1) Dump one line of 32 pixels at the time.
//       - The illustration below tells if a pixel is transparent or regular.
//       - Only regular and zero (transparent) pixels are explicitly stored in
//         the frame content. The transparent pixels are implicitly referred
//         from the illustration.
//
// Below is an illustration of the 32x32 image, where a space represents an
// implicit transparent pixel, a '0' represents an explicit transparent pixel
// and an 'x' represents an explicit regular pixel.
//
// Note: The output image will be "upside-down" compared to the illustration.
//
//    +--------------------------------+
//    |xx00                            |
//    |xxxx                            |
//    |xxxxxx00                        |
//    |xxxxxxxx                        |
//    |xxxxxxxxxx00                    |
//    |xxxxxxxxxxxx                    |
//    |xxxxxxxxxxxxxx00                |
//    |xxxxxxxxxxxxxxxx                |
//    |xxxxxxxxxxxxxxxxxx00            |
//    |xxxxxxxxxxxxxxxxxxxx            |
//    |xxxxxxxxxxxxxxxxxxxxxx00        |
//    |xxxxxxxxxxxxxxxxxxxxxxxx        |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxx00    |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxx    |
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx00|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    |xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx|
//    +--------------------------------+
//
// Type5 corresponds to a 32x32 images of a right facing trapezoid.
func DecodeFrameType5(frame []byte, width int, height int, pal color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	setPixel := GetPixelSetter(width, height)
	decodeCounts := []int{4, 4, 8, 8, 12, 12, 16, 16, 20, 20, 24, 24, 28, 28, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	for lineNum, decodeCount := range decodeCounts {
		var zeroCount int
		switch lineNum {
		case 0, 2, 4, 6, 8, 10, 12, 14:
			zeroCount = 2
		}
		regularCount := decodeCount - zeroCount
		decodeLineTransparencyRight(img, setPixel, frame, regularCount, zeroCount, pal)
		frame = frame[decodeCount:]
	}
	return img
}

// decodeLineTransparencyLeft decodes a line of the frame, where regularCount
// represent the number of explicit regular pixels, zeroCount the number of
// explicit transparent pixels and the rest of the line is transparent. Each
// line is assumed to have a width of 32 pixels.
func decodeLineTransparencyLeft(dst draw.Image, setPixel func(dst draw.Image, c color.Color), frame []byte, regularCount, zeroCount int, pal color.Palette) {
	// Total number of explicit pixels.
	decodeCount := zeroCount + regularCount

	// Implicit transparent pixels.
	for i := decodeCount; i < 32; i++ {
		setPixel(dst, color.Transparent)
	}
	// Explicit transparent pixels (zeroes).
	for i := 0; i < zeroCount; i++ {
		setPixel(dst, color.Transparent)
	}
	// Explicit regular pixels.
	for i := zeroCount; i < decodeCount; i++ {
		setPixel(dst, pal[frame[i]])
	}
}

// decodeLineTransparencyRight decodes a line of the frame, where regularCount
// represent the number of explicit regular pixels, zeroCount the number of
// explicit transparent pixels and the rest of the line is transparent. Each
// line is assumed to have a width of 32 pixels.
func decodeLineTransparencyRight(dst draw.Image, setPixel func(dst draw.Image, c color.Color), frame []byte, regularCount, zeroCount int, pal color.Palette) {
	// Total number of explicit pixels.
	decodeCount := zeroCount + regularCount

	// Explicit regular pixels.
	for i := 0; i < regularCount; i++ {
		setPixel(dst, pal[frame[i]])
	}

	// Explicit transparent pixels (zeroes).
	for i := 0; i < zeroCount; i++ {
		setPixel(dst, color.Transparent)
	}

	// Implicit transparent pixels.
	for i := decodeCount; i < 32; i++ {
		setPixel(dst, color.Transparent)
	}
}
