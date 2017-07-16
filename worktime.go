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
	"strconv"
	"math"
)

const LOG_NAME = "worktime.log"
const DEFAULT_DINNER_DURATION = 30
const TIME_FORMAT = "2006-01-02 15:04"
const TIME_FORMAT_DATE = "01-02"
const TIME_FORMAT_SHORT = "15:04"

type workDay struct {
	StartTime     		string `json:"startTime"`
	StopTime       		string `json:"stopTime"`
	DinnerMinutes 		int    `json:"dinner"`
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
		start(workDay{StartTime: time.Now().Format(TIME_FORMAT)})
	case "stop":
		updateLastRecord(workDay{StopTime: time.Now().Format(TIME_FORMAT)})
	case "dinner":
		if parameter != "" {
			dinnerMinutes, _ := strconv.Atoi(parameter)
			updateLastRecord(workDay{DinnerMinutes: dinnerMinutes})
		} else {
			help()
		}
	case "time":
		countTime()
	default:
		help()
	}
}

func help() {
	fmt.Println("Использование: worktime (start|stop|time|dinner minutes)")
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

	return logPath + "/" + LOG_NAME
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

func updateLastRecord(workDayPatch workDay) {
	file := openFile()
	defer file.Close()

	lastWorkDay, workDays := getWorkDays(file)

	clearLogFile()

	for _, workDay := range workDays {
		jsonEncodedMark, _ := json.Marshal(workDay)
		logString := fmt.Sprintln(string(jsonEncodedMark))
		_, error := file.WriteString(logString)
		checkError(error)
	}

	if lastWorkDay.DinnerMinutes == 0 {
		lastWorkDay.DinnerMinutes = DEFAULT_DINNER_DURATION
	}

	patchWordDay(&lastWorkDay, workDayPatch)

	jsonEncodedMark, _ := json.Marshal(lastWorkDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))
	fmt.Println(logString)
	file.WriteString(logString)
}

func patchWordDay(workDay *workDay, patch workDay) {
	if patch.DinnerMinutes > 0 {
		workDay.DinnerMinutes = patch.DinnerMinutes
	}

	if patch.StopTime != "" {
		workDay.StopTime = patch.StopTime
	}
}

func getWorkDays(file *os.File) (lastWorkDay workDay, workDays []workDay) {
	bf := bufio.NewReader(file)

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

	lastWorkDay = workDays[len(workDays)-1]
	workDays = workDays[:len(workDays)-1]

	return lastWorkDay, workDays
}

func start(workDay workDay) {
	file := openFile()
	defer file.Close()

	jsonEncodedMark, _ := json.Marshal(workDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	_, error := file.WriteString(logString)
	checkError(error)
}

func countTime() {
	file := openFile()
	defer file.Close()

	lastWorkDay, workDays := getWorkDays(file)
	workDays = append(workDays, lastWorkDay)

	fmt.Println("Дата  | Начал Конец | Обед \t| Переработка")
	fmt.Println("---------------------------------------------------")

	var hours float64
	var minutes float64
	for _, workDay := range workDays {
		startTime, error := time.Parse(TIME_FORMAT, workDay.StartTime)
		checkError(error)

		if workDay.StopTime == "" {
			continue
		}

		stopTime, error := time.Parse(TIME_FORMAT, workDay.StopTime)
		checkError(error)

		startTimeWithDinner := startTime.Add(time.Duration(workDay.DinnerMinutes) * time.Minute)
		dayDuration := stopTime.Sub(startTimeWithDinner)

		dayHours := math.Floor(dayDuration.Hours())
		dayMinutes := math.Floor((dayDuration.Hours() - dayHours) * 60)

		fmt.Println(fmt.Sprintf("%v | %v %v | %v мин \t| %v час %v мин",
			startTime.Format(TIME_FORMAT_DATE),
			startTime.Format(TIME_FORMAT_SHORT),
			stopTime.Format(TIME_FORMAT_SHORT),
			workDay.DinnerMinutes,
			dayHours,
			dayMinutes,
		))

		hours = hours + dayHours
		minutes = minutes + dayMinutes
	}

	fmt.Println("===================================================")

	hourBalance := hours + math.Floor(minutes / 60)
	minuteBalance := minutes - (math.Floor(minutes / 60) * 60)

	fmt.Println(fmt.Sprintf("Переработка: %v ч. %v мин.", hourBalance, minuteBalance))
}
