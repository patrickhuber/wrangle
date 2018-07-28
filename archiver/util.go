package archiver

import (
	fp "github.com/patrickhuber/wrangle/filepath"
)

func commonDirectory(filePaths ...string) string {
	switch len(filePaths) {
	case 0:
		return ""
	case 1:
		return fp.Dir(filePaths[0])
	}

	firstSplit := fp.SplitAll(filePaths[0])
	segments := make([]string, 0)
	for i, segment := range firstSplit {
		matchAll := true
		for _, unit := range filePaths[1:] {
			unitSplit := fp.SplitAll(unit)
			if i >= len(unitSplit) || unitSplit[i] != segment {
				matchAll = false
				break
			}
		}
		if matchAll {
			segments = append(segments, segment)
		}
	}

	// in cases where one of the files happens to match the common directory
	// return the parent of the common directory
	commonDir := fp.Join(segments...)
	for _, filePath := range filePaths {
		filePath = fp.ToSlash(filePath)
		if filePath == commonDir {
			return fp.Dir(filePath)
		}
	}
	return commonDir
}
