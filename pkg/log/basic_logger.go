package log

import (
    "log"
    "os"
    "sync"
)

type BasicLogger struct {
    ErrorLog   *log.Logger
    InfoLog    *log.Logger
    WarningLog *log.Logger
}

// Singleton
var (
    instance *BasicLogger
    once     sync.Once
)

func GetLoggerInstance() *BasicLogger {
    if instance == nil {
        once.Do(func() {
            instance = &BasicLogger{
                ErrorLog: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
                InfoLog: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
                WarningLog: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
            }
        })
    }
    return instance
}

func (bl *BasicLogger) Error(msg string) {
    bl.ErrorLog.Println(msg)
}

func (bl *BasicLogger) ErrorFormat(format string, err error) {
    bl.ErrorLog.Printf(format, err.Error())
}

func (bl *BasicLogger) Info(msg string) {
    bl.InfoLog.Println(msg)
}

func (bl *BasicLogger) Warn(msg string) {
    bl.WarningLog.Println(msg)
}
