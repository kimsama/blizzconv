package imgarchive

import "encoding/binary"
import "io"
import "os"

// ExtractCel extracts CEL images, based on the CEL archive format described
// below.
//
// CEL archive format:
//    // imageOffsets contains the offsets to the CEL images. (little endian)
//    imageOffsets   [imageCount]uint32
//    // data contains the CEL image content.
//    //    start: imageOffsets[imageNum]
//    //    end:   imageOffsets[imageNum + 1]
//    // Note: the last image has only an implicit end offset, which is the end of the file.
//    data           []byte
//
func ExtractCel(r *os.File, ws []*os.File) (err error) {
   imageCount := len(ws)
   imageOffsets := make([]uint32, imageCount)
   err = binary.Read(r, binary.LittleEndian, imageOffsets)
   if err != nil {
      return err
   }
   for imageNum := 0; imageNum < imageCount; imageNum++ {
      imageStart := int64(imageOffsets[imageNum])
      if imageNum == imageCount-1 {
         // Last image, so copy all that's left.
         _, err = io.Copy(ws[imageNum], r)
         if err != nil {
            return err
         }
      } else {
         imageSize := int64(imageOffsets[imageNum+1]) - imageStart
         _, err = io.CopyN(ws[imageNum], r, imageSize)
         if err != nil {
            return err
         }
      }
   }
   return nil
}
