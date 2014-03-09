// Package dunconf implements functions for retrieving relevant information
// required for parsing DUN files.
package dunconf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mewpkg/goini"
)

var dict ini.Dict

// IniPath is the path to an ini file which provides information about the
// starting coordinates of a given DUN file.
var IniPath string

// Init loads an ini file which provides relevant information required for
// parsing DUN files
func Init() (err error) {
	dict, err = ini.Load(IniPath)
	if err != nil {
		return err
	}
	return nil
}

// DungeonNames returns a slice of dungeon names based on the ini file.
func DungeonNames() (dungeonNames []string) {
	for dungeonName := range dict {
		if dungeonName == "" || strings.HasSuffix(dungeonName, ".dun") {
			continue
		}
		dungeonNames = append(dungeonNames, dungeonName)
	}
	sort.Strings(dungeonNames)
	return dungeonNames
}

// GetColStart returns the starting col of a given DUN file.
func GetColStart(dunName string) (colStart int, err error) {
	colStart, found := dict.GetInt(dunName, "col_start")
	if !found {
		return 0, fmt.Errorf("col_start not found for %q.", dunName)
	}
	return colStart, nil
}

// GetRowStart returns the starting row of a given DUN file.
func GetRowStart(dunName string) (rowStart int, err error) {
	rowStart, found := dict.GetInt(dunName, "row_start")
	if !found {
		return 0, fmt.Errorf("row_start not found for %q.", dunName)
	}
	return rowStart, nil
}

// GetDunNames returns the DUN file names of a given dungeon map.
func GetDunNames(dungeonName string) (dunNames []string, err error) {
	rawDunNames, found := dict.GetString(dungeonName, "duns")
	if !found {
		return nil, fmt.Errorf("duns not found for %q.", dungeonName)
	}
	return strings.Split(rawDunNames, ","), nil
}

// GetColCount returns the number of cols of a given dungeon map.
func GetColCount(dungeonName string) (colCount int, err error) {
	colCount, found := dict.GetInt(dungeonName, "col_count")
	if !found {
		return 0, fmt.Errorf("col_count not found for %q.", dungeonName)
	}
	return colCount, nil
}

// GetRowCount returns the number of rows of a given dungeon map.
func GetRowCount(dungeonName string) (rowCount int, err error) {
	rowCount, found := dict.GetInt(dungeonName, "row_count")
	if !found {
		return 0, fmt.Errorf("row_count not found for %q.", dungeonName)
	}
	return rowCount, nil
}
