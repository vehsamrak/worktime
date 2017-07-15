package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"strconv"
)

import "github.com/mitchellh/go-homedir"

const LOG_NAME = "/worktime.log"

type mark struct {
	Time          string `json:"time"`
	Status        string `json:"status"`
	DinnerMinutes int `json:"dinner"`
}

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		help()
		return
	}

	command := arguments[0]
	var parameter string

	if len(arguments) > 1 {
		parameter = arguments[1]
	}

	switch command {
	case "start":
		log(mark{Time: time.Now().Format("2006-01-02 15:04:05"), Status: "start"})
	case "stop":
		log(mark{Time: time.Now().Format("2006-01-02 15:04:05"), Status: "stop"})
	case "dinner":
		if parameter != "" {
			dinnerMinutes, _ := strconv.Atoi(parameter)
			log(mark{Time: time.Now().Format("2006-01-02 15:04:05"), Status: "dinner", DinnerMinutes: dinnerMinutes})
		} else {
			help()
		}
	default:
		help()
	}
}

func help() {
	fmt.Println("Использование: worktime (start|stop|dinner minutes)")
}

func log(mark mark) {
	logPath, _ := homedir.Dir()
	logPath = logPath + LOG_NAME
	var _, error = os.Stat(logPath)

	if os.IsNotExist(error) {
		fmt.Println("Log file not exist. Creating new one at", logPath)
		var file, error = os.Create(logPath)
		checkError(error)
		defer file.Close()
	}

	file, error := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	checkError(error)

	defer file.Close()

	jsonEncodedMark, _ := json.Marshal(mark)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	_, error = file.WriteString(logString)
	checkError(error)
}

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}
