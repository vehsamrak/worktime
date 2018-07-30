package message

type ConsoleMessage struct {
	message string
}

func (consoleMessage ConsoleMessage) getMessage() string {
	return consoleMessage.message
}

func New(messageText string) Message {
	return &ConsoleMessage{message: messageText}
}
