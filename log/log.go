package log

import (
    "bytes"
    "io"
    logs "log"
    "os"
)

var logger *logs.Logger

func Init() {
    writer1 := &bytes.Buffer{}
    writer2 := os.Stdout
    writer3, err := os.OpenFile("run.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        logs.Fatalf("create file log.txt failed: %v", err)
    }
    logger = logs.New(io.MultiWriter(writer1, writer2, writer3), "", logs.Lshortfile|logs.LstdFlags)
}

func Info(text string) {
    logger.Printf("[info]:%s", text)
}

func Debug(text string) {
    logger.Printf("[debug]:%s", text)
}

func Warn(text error) {
    logger.Panic("[warning]:%s", text)
}

func Error(text error) {
    logger.Fatal("[error]:%s", text)
}
