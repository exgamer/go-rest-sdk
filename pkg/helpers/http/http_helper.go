package http

import (
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"github.com/motemen/go-loghttp"
	"io"
	"io/ioutil"
	"net/http"
)

func DoPostHttpRequest(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return DoHttpRequest("POST", url, headers, body)
}

func DoGetHttpRequest(url string, headers map[string]string) ([]byte, error) {
	return DoHttpRequest("GET", url, headers, nil)
}

func DoHttpRequest(method string, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	client := http.Client{
		Transport: &loghttp.Transport{},
	}

	req, _ := http.NewRequest(method, url, body)

	for n, v := range headers {
		req.Header.Set(n, v)
	}

	response, err := client.Do(req)

	//io.Copy(os.Stdout, response.Body)

	if err != nil {
		logger.LogError(err)

		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}
