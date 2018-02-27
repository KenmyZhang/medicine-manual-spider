package app

import (
	"crypto/tls"
	"errors"
	iconv "github.com/djimenez/iconv-go"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func httpGet(url string, conv bool) (string, error) {
	var req *http.Request
	var httpError error
	time.Sleep(150 * time.Millisecond)
	if req, httpError = http.NewRequest("GET", url, nil); httpError != nil {
		return "", httpError
	}

	resp, err := HttpClient().Do(req)
	if err != nil {
		return "", errors.New("http get error:" + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("http read error:" + err.Error())
	}

	if conv == true {
		out := make([]byte, len(body))
		out = out[:]

		iconv.Convert(body, out, "gb2312", "utf-8")
		src := string(out)
		return src, nil
	} else {
		return string(body), nil
	}
}

const (
	connectTimeout = 6 * time.Second
	requestTimeout = 30 * time.Second
)

func HttpClient() *http.Client {
	return secureHttpClient
}

var (
	secureHttpClient = createHttpClient(false)
)

func createHttpClient(enableInsecureConnections bool) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   connectTimeout,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   connectTimeout,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: enableInsecureConnections,
			},
		},
		Timeout: requestTimeout,
	}

	return client
}
