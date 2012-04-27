// dun_dump is a tool for constructing dungeons, based on the information
// retrieved from a given DUN file, and storing these dungeons as png images.
//
// Usage:
//
//    dun_dump [OPTIONS]... [name.dun]...
//
// The OPTIONS are:
//
//    -a=false
//            Dump all dungeons.
//    -celini="cel.ini"
//            Path to an ini file containing image information.
//            Note: 'cl2.ini' will be used for files that have the '.cl2' extension.
//    -mpqdump="mpqdump/"
//            Path to an extracted MPQ file.
//    -mpqini="mpq.ini"
//            Path to an ini file containing relative path information.
package documentation
