package http

import (
	"github.com/exgamer/go-rest-sdk/pkg/helpers/http/httpHelperStruct"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"github.com/motemen/go-loghttp"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func DoPostHttpRequest(url string, headers map[string]string, body io.Reader) (httpHelperStruct.HttpResponse, error) {
	return DoHttpRequest("POST", url, headers, body)
}

func DoGetHttpRequest(url string, headers map[string]string) (httpHelperStruct.HttpResponse, error) {
	return DoHttpRequest("GET", url, headers, nil)
}

func DoHttpRequest(method string, url string, headers map[string]string, body io.Reader) (httpHelperStruct.HttpResponse, error) {
	client := http.Client{
		Timeout:   30 * time.Second,
		Transport: &loghttp.Transport{},
	}

	req, _ := http.NewRequest(method, url, body)

	for n, v := range headers {
		req.Header.Set(n, v)
	}

	response, err := client.Do(req)

	r := httpHelperStruct.HttpResponse{
		Url:     url,
		Method:  method,
		Headers: headers,
	}

	if err != nil {
		logger.LogError(err)

		return r, err
	}

	r.Status = response.Status
	r.StatusCode = response.StatusCode

	rBody, bErr := ioutil.ReadAll(response.Body)

	if bErr != nil {
		logger.LogError(bErr)

		return r, err
	}

	r.Body = rBody

	defer response.Body.Close()

	return r, bErr
}
