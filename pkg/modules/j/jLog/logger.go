package jLog

import (
	"log"
	"os"
	"strconv"
	"time"
)

func PrintInfoJLog(serviceName string, method string, uri string, status int, requestId string, message string) {
	PrintJLog("INFO", serviceName, method, uri, status, requestId, message)
}

func PrintErrorJLog(serviceName string, method string, uri string, status int, requestId string, message string) {
	PrintJLog("ERROR", serviceName, method, uri, status, requestId, message)
}

func PrintJLog(level string, serviceName string, method string, uri string, status int, requestId string, message string) {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	dateTime := time.Now().Format("2006-01-02 15:04:05.345")
	serviceData := "[" + serviceName + "," + requestId + "]"
	rData := "[" + method + "," + uri + "," + strconv.Itoa(status) + "]"

	log.Println(dateTime + " " + level + " " + serviceData + rData + " " + message)
}
