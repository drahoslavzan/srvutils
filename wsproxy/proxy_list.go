package wsproxy

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type (
	ProxyList struct {
		list []string
	}
)

func NewProxyList(url string) (*ProxyList, error) {
	list, err := parseProxies(url)
	if err != nil {
		return nil, err
	}

	ret := &ProxyList{
		list: list,
	}

	return ret, nil
}

func (m *ProxyList) List() []string {
	return m.list
}

func (m *ProxyList) Clients(cfg *Config) (*ProxyClients, error) {
	clients := make([]*Client, len(m.list))
	for i, url := range m.list {
		cln, err := NewHttpProxyClient(url, cfg)
		if err != nil {
			return nil, fmt.Errorf("new http proxy client: %w", err)
		}

		clients[i] = cln
	}

	ret := &ProxyClients{
		clients: clients,
	}

	return ret, nil
}

func parseProxies(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer res.Body.Close()

	var list []string

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, ":")

		if len(s) != expectedProxyStringSplitCount {
			return nil, fmt.Errorf("invalid proxy list, got line: %s", line)
		}

		list = append(list, fmt.Sprintf("http://%s:%s@%s:%s", s[2], s[3], s[0], s[1]))
	}

	return list, nil
}

const expectedProxyStringSplitCount = 4
