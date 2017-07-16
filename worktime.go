package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io"
	"log"
	"bufio"
	"io/ioutil"
)

const LOG_NAME = "/worktime.log"
const TIME_FORMAT = "2006-01-02 15:04:05"

type workDay struct {
	StartTime     		string `json:"startTime"`
	StopTime       		string `json:"stopTime"`
	DinnerMinutes 		int    `json:"dinner"`
	TimeBalanceHours   	int    `json:"balanceHours"`
	TimeBalanceMinutes  int    `json:"balanceMinutes"`
}

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		help()
		return
	}

	command := arguments[0]
	//var parameter string
	//
	//if len(arguments) > 1 {
	//	parameter = arguments[1]
	//}

	switch command {
	case "start":
		start(workDay{StartTime: time.Now().Format(TIME_FORMAT)})
	case "stop":
		stop()
	//case "dinner":
	//	if parameter != "" {
	//		dinnerMinutes, _ := strconv.Atoi(parameter)
	//		log(workDay{Time: time.Now().Format("2006-01-02 15:04:05"), Status: "dinner", DinnerMinutes: dinnerMinutes})
	//	} else {
	//		help()
	//	}
	default:
		help()
	}
}

func help() {
	fmt.Println("Использование: worktime (start|stop|dinner minutes)")
}

func openFile() *os.File {
	logPath := getFilePath()
	var _, error = os.Stat(logPath)

	if os.IsNotExist(error) {
		fmt.Println("Log file not exist. Creating new one at", logPath)
		var file, error = os.Create(logPath)
		checkError(error)
		defer file.Close()
	}

	file, error := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR, 0644)
	checkError(error)

	return file
}

func getFilePath() string {
	logPath, _ := homedir.Dir()

	return logPath + LOG_NAME
}

func checkError(error error) {
	if error != nil {
		panic(error)
	}
}

func clearLogFile() {
	err := ioutil.WriteFile(getFilePath(), []byte(""), 0644)
	checkError(err)
}

//func countTime(workDays []workDay) (hours int, minutes int) {
//	for _, workDay := range workDays {
//		then, error := time.Parse(TIME_FORMAT, workDay.StartTime)
//		checkError(error)
//
//		duration := time.Since(then)
//		fmt.Println(duration.Hours())
//
//		hours += workDay.StartTime
//	}
//
//	return hours, minutes
//}

func start(workDay workDay) {
	file := openFile()
	defer file.Close()

	jsonEncodedMark, _ := json.Marshal(workDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	_, error := file.WriteString(logString)
	checkError(error)
}

func stop() {
	file := openFile()
	defer file.Close()

	bf := bufio.NewReader(file)

	var workDays []workDay

	for {
		line, _, err := bf.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		var c workDay
		json.Unmarshal(line, &c)

		workDays = append(workDays, c)
	}

	lastDay := workDays[len(workDays)-1]
	workDays = workDays[:len(workDays)-1]

	clearLogFile()

	for _, workDay := range workDays {
		jsonEncodedMark, _ := json.Marshal(workDay)
		logString := fmt.Sprintln(string(jsonEncodedMark))
		_, error := file.WriteString(logString)
		checkError(error)
	}

	lastDay.StopTime = time.Now().Format(TIME_FORMAT)

	jsonEncodedMark, _ := json.Marshal(lastDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))
	fmt.Println(logString)
	file.WriteString(logString)
}
