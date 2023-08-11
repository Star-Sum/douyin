package Log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

// LogFileGet 获取日志文件地址
func LogFileGet() (*os.File, error) {
	file, err := os.OpenFile("./logfile.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Panic("Log File Error!")
		return file, err
	}
	return file, err
}

// LogWrite 写入日志信息
func LogWrite(Info string, errInput error) {
	file, err := LogFileGet()
	if err != nil {
		log.Panic("Log File Error!")
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(Info)
	if errInput != nil {
		errString := errInput.Error()
		log.Println(errString)
		log.Println(string(debug.Stack()))
	}
}

// ErrorLogWithPanic 处理错误信息,有panic
func ErrorLogWithPanic(WrongInfo string, err error) {
	fmt.Println(WrongInfo + "\n")
	LogWrite(WrongInfo, err)
	panic(err)
}

// ErrorLogWithoutPanic 处理错误信息,无panic
func ErrorLogWithoutPanic(WrongInfo string, err error) {
	fmt.Println(WrongInfo + "\n")
	LogWrite(WrongInfo, err)
}

// NormalLog 当结果正常时输出日志信息
func NormalLog(RightInfo string, err error) {
	LogWrite(RightInfo, err)
}
