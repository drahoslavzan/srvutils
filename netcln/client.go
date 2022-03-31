package netcln

import (
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
		Get(link string) (status int, body string, error error)
		PostURLEncoded(link string, values map[string]string) (status int, body string, error error)
	}

	client struct {
		handler ClientHandler
	}
)

func NewClient(handler ClientHandler) Client {
	return &client{handler}
}

func (m *client) Get(link string) (int, string, error) {
	res, err := m.handler.Do(getRequest("GET", link, nil))
	if err != nil {
		return 0, "", err
	}

	text, err := getRespText(res)
	return res.StatusCode, text, err
}

func (m *client) PostURLEncoded(link string, values map[string]string) (int, string, error) {
	res, err := m.handler.Do(getRequest("POST", link, values))
	if err != nil {
		return 0, "", err
	}

	text, err := getRespText(res)
	return res.StatusCode, text, err
}

func getRequest(method string, link string, values map[string]string) *http.Request {
	var body io.Reader
	if values != nil {
		vals := url.Values{}
		for k, v := range values {
			vals.Add(k, v)
		}
		body = strings.NewReader(vals.Encode())
	}

	req, err := http.NewRequest(method, link, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req
}

func getRespText(res *http.Response) (string, error) {
	defer res.Body.Close()
	text, err := ioutil.ReadAll(res.Body)
	return string(text), err
}
