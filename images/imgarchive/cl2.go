package imgarchive

import "encoding/binary"
import "io"
import "os"

// ExtractCl2 extracts CL2 images, based on the CL2 archive format described
// below.
//
// CL2 archive format:
//    // headerOffsets contains the offsets to the cl2Headers. (little endian)
//    headerOffsets  [imageCount]uint32
//    // cl2Headers contains the CL2 Headers, but with offsets relative to the current headerOffset
//    cl2Headers     [imageCount][]byte
//    // data contains the CL2 image content, excluding header data.
//    //    start: headerOffsets[imageNum] + frameOffsets[0]
//    //    end:   headerOffsets[imageNum] + frameOffsets[frameCount]
//    // Note: Both frameOffsets and frameCount are located in cl2Headers[imageNum].
//    data           []byte
func ExtractCl2(r *os.File, ws []*os.File) (err error) {
   imageCount := len(ws)
   headerOffsets := make([]uint32, imageCount)
   err = binary.Read(r, binary.LittleEndian, headerOffsets)
   if err != nil {
      return err
   }
   for imageNum := 0; imageNum < imageCount; imageNum++ {
      headerOffset := headerOffsets[imageNum]
      _, err = r.Seek(int64(headerOffset), os.SEEK_SET)
      if err != nil {
         return err
      }
      var frameCount uint32
      err = binary.Read(r, binary.LittleEndian, &frameCount)
      if err != nil {
         return err
      }
      frameOffsets := make([]uint32, frameCount+1)
      err = binary.Read(r, binary.LittleEndian, frameOffsets)
      if err != nil {
         return err
      }
      w := ws[imageNum]
      err = binary.Write(w, binary.LittleEndian, frameCount)
      if err != nil {
         return err
      }
      for frameNum := uint32(0); frameNum <= frameCount; frameNum++ {
         diff := frameOffsets[0] - (1+frameCount+1)*4
         err = binary.Write(w, binary.LittleEndian, frameOffsets[frameNum]-diff)
         if err != nil {
            return err
         }
      }
      imageStart := int64(frameOffsets[0] + headerOffset)
      _, err = r.Seek(imageStart, os.SEEK_SET)
      if err != nil {
         return err
      }
      imageSize := int64(frameOffsets[frameCount] - frameOffsets[0])
      _, err = io.CopyN(w, r, imageSize)
      if err != nil {
         return err
      }
   }
   return nil
}
