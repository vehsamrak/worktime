package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const DefaultDinnerDuration = 30
const LogName = "worktime.log"
const TimeFormat = "2006-01-02 15:04"
const TimeFormatDate = "01-02"
const TimeFormatShort = "15:04"
const WorkHoursNumber = 8

type workDay struct {
	StartTime     string `json:"startTime"`
	StopTime      string `json:"stopTime"`
	DinnerMinutes int    `json:"dinner"`
}

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		help()
		return
	}

	command := arguments[0]
	var parameter string
	var secondParameter int

	if len(arguments) > 1 {
		parameter = arguments[1]
	}

	if len(arguments) > 2 {
		secondParameter, _ = strconv.Atoi(arguments[2])
	}

	switch command {
	case "start":
		start(workDay{StartTime: time.Now().Format(TimeFormat)})
	case "stop":
		updateLastRecord(workDay{StopTime: time.Now().Format(TimeFormat)})
	case "dinner":
		if parameter != "" {
			dinnerMinutes, _ := strconv.Atoi(parameter)
			updateLastRecord(workDay{DinnerMinutes: dinnerMinutes})
		} else {
			help()
		}
	case "time":
		var verboseLog bool
		var tailNumber int

		if parameter == "full" {
			verboseLog = true
		} else if parameter != "" {
			secondParameter, _ = strconv.Atoi(parameter)
		}

		if secondParameter > 0 {
			tailNumber = secondParameter
		}

		countTime(tailNumber, verboseLog)
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
	fmt.Println("   time full\t\tПросморт полного лога рабочего времени")
	fmt.Println("   time full [X]\tПросморт лога рабочего времени за X последних дней")
	fmt.Println("   help \t\tПросмотр текущей справки")
}

func openFile() *os.File {
	logPath := getFilePath()
	var _, err = os.Stat(logPath)

	if os.IsNotExist(err) {
		fmt.Println("Log file not exist. Creating new one at", logPath)
		var file, err = os.Create(logPath)
		checkError(err)
		defer file.Close()
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR, 0644)
	checkError(err)

	return file
}

func getFilePath() string {
	logPath, _ := homedir.Dir()

	return logPath + "/" + LogName
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
		_, err := file.WriteString(logString)
		checkError(err)
	}

	if lastWorkDay.DinnerMinutes == 0 {
		lastWorkDay.DinnerMinutes = DefaultDinnerDuration
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

	if len(workDays) > 0 {
		lastWorkDay = workDays[len(workDays)-1]
	}

	if len(workDays) > 0 {
		workDays = workDays[:len(workDays)-1]
	}

	return lastWorkDay, workDays
}

func start(workDay workDay) {
	file := openFile()
	defer file.Close()

	lastWorkDay, _ := getWorkDays(file)

	if lastWorkDay.StartTime != "" {
		lastStartDate, err := time.Parse(TimeFormat, lastWorkDay.StartTime)
		checkError(err)

		if lastStartDate.Day() == time.Now().Day() {
			fmt.Printf("Current work day was already started. Please edit %v if you like.\n", getFilePath())

			return
		}
	}

	jsonEncodedMark, _ := json.Marshal(workDay)
	logString := fmt.Sprintln(string(jsonEncodedMark))

	fmt.Println(logString)

	_, err := file.WriteString(logString)
	checkError(err)
}

func countTime(tailNumber int, verboseLog bool) {
	file := openFile()
	defer file.Close()

	lastWorkDay, workDays := getWorkDays(file)
	workDays = append(workDays, lastWorkDay)

	if verboseLog {
		fmt.Println("Дата  | Начал Конец | Обед \t| Переработка")
		fmt.Println("---------------------------------------------------")
	}

	cutWorkDaysStatistics := tailNumber > 0 && len(workDays) >= tailNumber

	if cutWorkDaysStatistics {
		workDays = workDays[len(workDays)-tailNumber-1:]
	}

	var minuteBalance float64
	for _, workDay := range workDays {
		startTime, err := time.Parse(TimeFormat, workDay.StartTime)
		checkError(err)

		if workDay.StopTime == "" {
			continue
		}

		stopTime, err := time.Parse(TimeFormat, workDay.StopTime)
		checkError(err)

		dinnerDuration := time.Duration(workDay.DinnerMinutes) * time.Minute
		expectedWorkDayDuration := time.Duration(WorkHoursNumber * time.Hour)
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
				startTime.Format(TimeFormatDate),
				startTime.Format(TimeFormatShort),
				stopTime.Format(TimeFormatShort),
				workDay.DinnerMinutes,
				workTimingString))
		}

		minuteBalance = minuteBalance + fullDayMinutes
	}

	if verboseLog {
		fmt.Println("===================================================")
	}

	if cutWorkDaysStatistics {
		fmt.Printf("Показано время за последних дней: %v.\n", tailNumber)
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
