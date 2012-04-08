// Package mpq provides access to an extracted MPQ archive.
package mpq

import ini "github.com/glacjay/goini"

import "fmt"

var dict ini.Dict

// IniPath is the path to an ini file which provides relative path information
// for files in an extracted MPQ archive.
var IniPath string

// Init loads an ini file which provides relative path information for files in
// an extracted MPQ archive.
func Init() (err error) {
   dict, err = ini.Load(IniPath)
   if err != nil {
      return err
   }
   return nil
}

// ExtractPath is the path to an extracted MPQ file.
var ExtractPath string

// GetPath returns the full path of name.
func GetPath(name string) (path string, err error) {
   relPath, err := GetRelPath(name)
   if err != nil {
      return "", err
   }
   return ExtractPath + relPath, nil
}

// GetRelPath returns the relative path of name.
func GetRelPath(name string) (relPath string, err error) {
   relPath, found := dict.GetString(name, "path")
   if !found {
      return "", fmt.Errorf("path not found for '%s'.", name)
   }
   return relPath, nil
}
