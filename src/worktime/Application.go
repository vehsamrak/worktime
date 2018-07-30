package worktime

import "github.com/Vehsamrak/worktime/src/message"

type Application struct {
}

func (application Application) getHelpMessage() message.Message {
	return message.New("help message")
}
