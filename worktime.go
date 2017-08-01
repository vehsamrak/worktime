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
const WORK_HOURS_NUMBER = 8

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
		var verboseLog bool = false

		if parameter == "full" {
			verboseLog = true
		}

		countTime(verboseLog)
	default:
		help()
	}
}

func help() {
	fmt.Println("Использование: worktime (start|stop|time [full]|dinner (minutes))")
	fmt.Println("   start \t\tОтметка о начале рабочего дня")
	fmt.Println("   stop \t\tОтметка об окончании рабочего дня")
	fmt.Println("   dinner (minutes) \tЗапись количества минут проведенных на отдыхе или обеде")
	fmt.Println("   time \t\tПросмотр временного баланса переработок или недоработок")
	fmt.Println("   time full \t\tПросморт полного лога рабочего времени")
	fmt.Println("   help \t\tПросмотр текущей справки")
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

    lastWorkDay, _ := getWorkDays(file)

    lastStartDate, error := time.Parse(TIME_FORMAT, lastWorkDay.StartTime)
    checkError(error)

    if lastStartDate.Day() == time.Now().Day() {
        fmt.Printf("Current work day was already started. Please edit %v if you like.\n", getFilePath())

        return
    }

	jsonEncodedMark, _ := json.Marshal(workDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	_, error = file.WriteString(logString)
	checkError(error)
}

func countTime(verboseLog bool) {
	file := openFile()
	defer file.Close()

	lastWorkDay, workDays := getWorkDays(file)
	workDays = append(workDays, lastWorkDay)

	if verboseLog {
		fmt.Println("Дата  | Начал Конец | Обед \t| Переработка")
		fmt.Println("---------------------------------------------------")
	}

	var minuteBalance float64
	for _, workDay := range workDays {
		startTime, error := time.Parse(TIME_FORMAT, workDay.StartTime)
		checkError(error)

		if workDay.StopTime == "" {
			continue
		}

		stopTime, error := time.Parse(TIME_FORMAT, workDay.StopTime)
		checkError(error)

		dinnerDuration := time.Duration(workDay.DinnerMinutes) * time.Minute
		expectedWorkDayDuration := time.Duration(WORK_HOURS_NUMBER * time.Hour)
		overTimeWork := stopTime.Sub(startTime) - expectedWorkDayDuration - dinnerDuration
		overTimeWorkHours := overTimeWork.Hours()

		var fullDayMinutes float64
		var dayHours float64
		if overTimeWork >= 0 {
			fullDayMinutes = math.Floor(overTimeWorkHours * 60)
			dayHours = math.Floor(fullDayMinutes / 60)
		} else {
			fullDayMinutes = math.Ceil(overTimeWorkHours * 60)
			dayHours = math.Ceil(fullDayMinutes / 60)
		}

		if verboseLog {
			dayMinutes := fullDayMinutes - (dayHours * 60)

			var workTimingString string
			if dayHours == 0 {
				workTimingString = fmt.Sprintf("%v мин", dayMinutes)
			} else {
				workTimingString = fmt.Sprintf("%v час %v мин", dayHours, math.Abs(dayMinutes))
			}

			fmt.Println(fmt.Sprintf("%v | %v %v | %v мин \t| %v",
				startTime.Format(TIME_FORMAT_DATE),
				startTime.Format(TIME_FORMAT_SHORT),
				stopTime.Format(TIME_FORMAT_SHORT),
				workDay.DinnerMinutes,
				workTimingString))
		}

		minuteBalance = minuteBalance + fullDayMinutes
	}

	if verboseLog {
		fmt.Println("===================================================")
	}

	var hourBalance float64
	var balanceStatus string
	if minuteBalance >= 0 {
		hourBalance = math.Floor(minuteBalance / 60)
		balanceStatus = "Переработка"
	} else {
		hourBalance = math.Ceil(minuteBalance / 60)
		balanceStatus = "Недоработка"
	}

	minuteBalance = minuteBalance - (hourBalance * 60)

	var totalWorkTimingString string
	if hourBalance == 0 {
		totalWorkTimingString = fmt.Sprintf("%v мин", math.Abs(minuteBalance))
	} else {
		totalWorkTimingString = fmt.Sprintf("%v час %v мин", math.Abs(hourBalance), math.Abs(minuteBalance))
	}

	fmt.Println(fmt.Sprintf("%v: %v", balanceStatus, totalWorkTimingString))
}
