package capman

import "path/filepath"

func makeRecordFilePath(imagesFolder string, id string) (path string) {
	return filepath.Join(imagesFolder, id+"."+ImageFileExt)
}
