package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"strconv"
)

type mark struct {
	Time          string `json:"time"`
	Status        string `json:"status"`
	DinnerMinutes int `json:"dinnerMinutes"`
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
	fileName := "./worktime.log"

	file, error := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if error != nil {
		panic(error)
	}

	defer file.Close()

	jsonEncodedMark, _ := json.Marshal(mark)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	if _, error = file.WriteString(logString); error != nil {
		panic(error)
	}
}
