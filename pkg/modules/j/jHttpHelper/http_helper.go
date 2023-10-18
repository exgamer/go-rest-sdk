package jHttpHelper

import (
	"github.com/exgamer/go-rest-sdk/pkg/helpers/http"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"github.com/exgamer/go-rest-sdk/pkg/modules/j/jLog"
	"github.com/exgamer/go-rest-sdk/pkg/modules/j/jStructures"
	"io"
	"strings"
	"time"
)

func DoPostHttpRequest(requestData *jStructures.RequestData, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return GetResponseBody(requestData, "POST", url, headers, body)
}

func DoGetHttpRequest(requestData *jStructures.RequestData, url string, headers map[string]string) ([]byte, error) {
	return GetResponseBody(requestData, "GET", url, headers, nil)
}

func GetResponseBody(requestData *jStructures.RequestData, method string, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	start := time.Now()
	response, err := http.DoHttpRequest(method, url, headers, body)
	execTime := time.Since(start)

	if err != nil {
		logger.LogError(err)

		return nil, err
	}

	messageBuilder := strings.Builder{}
	messageBuilder.WriteString("Url: " + response.Method + " " + response.Status + " " + response.Url)
	messageBuilder.WriteString(" Exec time:" + execTime.String())

	if err != nil {
		messageBuilder.WriteString(" Error:" + err.Error())
	}

	if response.StatusCode != 200 {
		messageBuilder.WriteString(" Response:" + string(response.Body))
	}

	jLog.PrintInfoJLog(requestData.ServiceName, requestData.RequestMethod, requestData.RequestHost+requestData.RequestUrl, 0, requestData.RequestId, messageBuilder.String())

	return response.Body, nil
}
