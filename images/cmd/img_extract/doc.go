// img_extract is a tool for extracting CEL and CL2 archives.
//
// Usage:
//
//    img_extract [OPTION]... [name.cel|name.cl2]...
//
// Flags:
//
//    -celini="cel.ini"
//            Path to an ini file containing image information.
//            Note: 'cl2.ini' will be used for files that have the '.cl2' extension.
//    -mpqdump="mpqdump/"
//            Path to an extracted MPQ file.
//    -mpqini="mpq.ini"
//            Path to an ini file containing relative path information.
package main
