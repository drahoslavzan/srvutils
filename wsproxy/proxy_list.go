package wsproxy

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type (
	ProxyList struct {
		list []string
	}
)

func NewProxyList(link string) (*ProxyList, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("error parsing link: %w", err)
	}

	q := u.Query()

	skip := 0
	qs := q.Get("skip")
	if qs != "" {
		var err error
		skip, err = strconv.Atoi(qs)
		if err != nil {
			return nil, fmt.Errorf("invalid skip parameter '%s': %w", qs, err)
		}
		if skip < 0 {
			return nil, fmt.Errorf("skip parameter cannot be negative: %d", skip)
		}
	}

	take := 0
	qt := q.Get("take")
	if qt != "" {
		var err error
		take, err = strconv.Atoi(qt)
		if err != nil {
			return nil, fmt.Errorf("invalid take parameter '%s': %w", qt, err)
		}
		if take < 1 {
			return nil, fmt.Errorf("take parameter must be positive: %d", take)
		}
	}

	list, err := parseProxies(link)
	if err != nil {
		return nil, err
	}

	if skip > 0 {
		if skip >= len(list) {
			return nil, fmt.Errorf("skipping %d proxies out of %d", skip, len(list))
		}
		list = list[skip:]
	}
	if take > 0 {
		list = list[:min(len(list), take)]
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
