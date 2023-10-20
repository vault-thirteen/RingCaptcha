package rcs

import (
	"os"
	"strings"
)

const (
	ConfigurationFilePathDefault = "RCS.json"
)

type CommandLineArguments struct {
	ConfigurationFilePath string
}

func ReadCLA() (cla *CommandLineArguments, err error) {
	cla = &CommandLineArguments{}

	if len(os.Args) != 2 {
		cla.ConfigurationFilePath = ConfigurationFilePathDefault
		return cla, nil
	}

	cla.ConfigurationFilePath = strings.TrimSpace(os.Args[1])

	return cla, nil
}

// IsDefaultFile tells whether the default file path is used for the
// configuration file.
func (cla *CommandLineArguments) IsDefaultFile() (isDefaultFile bool) {
	return cla.ConfigurationFilePath == ConfigurationFilePathDefault
}
