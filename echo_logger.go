package log4go

type echoLog struct {
}

func (e *echoLog) Write(p []byte) (n int, err error) {
	writeGormLog(string(p))
	return 0, nil
}

func NewEchoLogger() *echoLog {
	return &echoLog{}
}
