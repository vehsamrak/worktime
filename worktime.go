package main

import (
	"fmt"
	"os"
	"time"
)

type mark struct {
	time   time.Time
	status string
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
		log("start", time.Now())
	case "stop":
		log("stop", time.Now())
	case "dinner":
		if parameter != "" {
			dinner(parameter)
		} else {
			help()
		}
	default:
		help()
	}
}
func dinner(dinnerMinutesString string) {
	log("dinner | " + dinnerMinutesString, time.Now())
}

func help() {
	fmt.Println("Использование: worktime (start|stop|dinner minutes)")
}

func log(status string, time time.Time) {
	fileName := "./worktime.log"

	file, error := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if error != nil {
		panic(error)
	}

	defer file.Close()

	logString := fmt.Sprintln(time.Format("2006-01-02 15:04:05"), "|", status)

	if _, error = file.WriteString(logString); error != nil {
		panic(error)
	}
}
