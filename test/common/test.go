package common

import "path/filepath"

const OutputFolderName = "out"

func GetTestDataFolder() (folderPath string) {
	return filepath.Join("test", "data")
}
