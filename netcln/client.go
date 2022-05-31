package netcln

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	userAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`
)

type (
	ClientHandler interface {
		Do(*http.Request) (*http.Response, error)
	}

	Client interface {
		Get(link string, headers map[string]string) (status int, body string, error error)
		Delete(link string, headers map[string]string) (status int, body string, error error)
		PostURLEncoded(link string, values map[string]string) (status int, body string, error error)
		Send(method, link, data string, headers map[string]string) (int, string, error)
		SendJSON(method, link string, data any, headers map[string]string) (int, string, error)
	}

	client struct {
		handler ClientHandler
	}
)

func NewClient(handler ClientHandler) Client {
	return &client{handler}
}

func (m *client) Get(link string, headers map[string]string) (int, string, error) {
	req := getRequest("GET", link, nil, headers)
	return m.send(req)
}

func (m *client) Delete(link string, headers map[string]string) (int, string, error) {
	req := getRequest("DELETE", link, nil, headers)
	return m.send(req)
}

func (m *client) PostURLEncoded(link string, values map[string]string) (int, string, error) {
	vals := url.Values{}
	for k, v := range values {
		vals.Add(k, v)
	}

	body := strings.NewReader(vals.Encode())
	req := getRequest("POST", link, body, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return m.send(req)
}

func (m *client) Send(method, link, data string, headers map[string]string) (int, string, error) {
	body := bytes.NewBufferString(data)
	req := getRequest(method, link, body, headers)
	return m.send(req)
}

func (m *client) SendJSON(method, link string, data any, headers map[string]string) (int, string, error) {
	var body io.Reader
	if data != nil {
		js, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(js)
	}

	req := getRequest(method, link, body, headers)
	req.Header.Set("Content-Type", "application/json")

	return m.send(req)
}

func (m *client) send(req *http.Request) (int, string, error) {
	res, err := m.handler.Do(req)
	if err != nil {
		return 0, "", err
	}

	text, err := getRespText(res)
	return res.StatusCode, text, err
}

func getRequest(method string, link string, body io.Reader, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, link, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", userAgent)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return req
}

func getRespText(res *http.Response) (string, error) {
	defer res.Body.Close()
	text, err := ioutil.ReadAll(res.Body)
	return string(text), err
}
