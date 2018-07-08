package model

type WorkDay struct {
	StartTime     string `json:"startTime"`
	StopTime      string `json:"stopTime"`
	DinnerMinutes int    `json:"dinner"`
	Correction    int    `json:"correction"`
	Comment       string `json:"comment"`
}
