package server

import "path/filepath"

func MakeRecordFilePath(imagesFolder string, id string) (path string) {
	return filepath.Join(imagesFolder, MakeFileName(id))
}

func MakeFileName(id string) (fileName string) {
	return id + "." + ImageFileExt
}
