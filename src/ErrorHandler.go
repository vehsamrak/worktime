package src

type ErrorHandler struct {
}

func (errorHandler *ErrorHandler) Check(error error) {
	if error != nil {
		panic(error)
	}
}
