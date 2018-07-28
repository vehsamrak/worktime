package src

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

const LogDirectory = ".worktime/"
const LogPath = "worktime.log"

type FileProcessor struct {
	ErrorHandler *ErrorHandler
}

func (fileProcessor *FileProcessor) OpenFile() *os.File {
	logDirectory := fileProcessor.getLogDirectory()
	logPath := fileProcessor.GetFilePath()
	var _, err = os.Stat(logPath)

	if os.IsNotExist(err) {
		fmt.Println("Log file doesn't exist. Creating new one at", logPath)
	}

	os.MkdirAll(logDirectory, 0777)
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	fileProcessor.ErrorHandler.Check(err)

	return file
}

func (fileProcessor *FileProcessor) getLogDirectory() string {
	homeDirectory, _ := homedir.Dir()

	return homeDirectory + "/" + LogDirectory
}

func (fileProcessor *FileProcessor) GetFilePath() string {
	return fileProcessor.getLogDirectory() + LogPath
}

func (fileProcessor *FileProcessor) ClearLogFile() {
	err := ioutil.WriteFile(fileProcessor.GetFilePath(), []byte(""), 0644)
	fileProcessor.ErrorHandler.Check(err)
}
