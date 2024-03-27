package helper


func HandleErrorWithPanic(err error) {
	if err != nil {
		SaveToLogError(err.Error())
		panic(err.Error())
	}
}