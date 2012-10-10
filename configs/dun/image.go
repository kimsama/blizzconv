package dun

import "image"
import "image/draw"

import "github.com/mewrnd/blizzconv/configs/min"

// Image returns an image constructed from the pillars associated with each
// coordinate of the dungeon map.
//
// ref: GetPillarRect (illustration of map coordinate system)
func (dungeon *Dungeon) Image(colCount, rowCount int, pillars []min.Pillar, levelFrames []image.Image) (img image.Image) {
	pillarHeight := pillars[0].Height()
	mapWidth := colCount*min.BlockWidth + rowCount*min.BlockWidth
	mapHeight := colCount*(min.BlockHeight/2) + rowCount*(min.BlockHeight/2) + (pillarHeight - min.BlockHeight)
	dst := image.NewRGBA(image.Rect(0, 0, mapWidth, mapHeight))
	for row := 0; row < rowCount; row++ {
		for col := 0; col < colCount; col++ {
			pillarNum := dungeon[col][row]
			if pillarNum != -1 {
				rect := GetPillarRect(col, row, mapWidth, pillarHeight)
				src := pillars[pillarNum].Image(levelFrames)
				draw.Draw(dst, rect, src, image.ZP, draw.Over)
			}
		}
	}
	return dst
}

// GetPillarRect returns an image.Rectangle based on the col and row
// coordinates. The calculations are based on the map coordinate system
// illustrated below:
//
// Map coordinate system:
//                 (0, 0)
//
//                   /\
//                r /\/\ c
//               o /\/\/\ o
//              w /\/\/\/\ l
//               /\/\/\/\/\
//    (0, 111)   \/\/\/\/\/   (111, 0)
//                \/\/\/\/
//                 \/\/\/
//                  \/\/
//                   \/
//
//               (111, 111)
func GetPillarRect(col, row, mapWidth, pillarHeight int) (rect image.Rectangle) {
	minX := mapWidth/2 - min.BlockWidth - row*min.BlockWidth + col*min.BlockWidth
	minY := row*(min.BlockHeight/2) + col*(min.BlockHeight/2)
	maxX := minX + min.PillarWidth
	maxY := minY + pillarHeight
	return image.Rect(minX, minY, maxX, maxY)
}
