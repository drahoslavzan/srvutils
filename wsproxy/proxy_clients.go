package wsproxy

import (
	"fmt"
	"math/rand"

	"go.uber.org/zap"
)

type (
	ProxyClients struct {
		clients []*Client
	}
)

func (m *ProxyClients) Remove(client *Client) {
	sz := len(m.clients)
	if sz < 1 {
		return
	}
	if client == nil {
		return
	}

	new := make([]*Client, 0, sz-1)
	for _, c := range m.clients {
		if c == client {
			continue
		}

		new = append(new, c)
	}

	m.clients = new
}

func (m *ProxyClients) RemoveFailed(indices []int, logger *zap.Logger) int {
	if len(m.clients) < 1 {
		return 0
	}
	if len(indices) < 1 {
		return len(m.clients)
	}

	rem := make(map[int]struct{}, len(indices))
	for _, i := range indices {
		rem[i] = struct{}{}
	}

	new := make([]*Client, 0, len(m.clients)-len(indices))
	for i, c := range m.clients {
		if _, remove := rem[i]; remove {
			logger.Warn("proxy client failed", zap.String("client", c.String()))
			continue
		}

		new = append(new, c)
	}

	m.clients = new
	return len(new)
}

func (m *ProxyClients) At(idx int) (*Client, error) {
	sz := len(m.clients)
	if sz < 1 {
		return nil, fmt.Errorf("no clients available")
	}
	if idx < 0 || idx >= sz {
		return nil, fmt.Errorf("invalid index: %d", idx)
	}

	return m.clients[idx], nil
}

func (m *ProxyClients) First() (*Client, error) {
	return m.At(0)
}

func (m *ProxyClients) Random() (*Client, error) {
	return m.At(rand.Intn(len(m.clients)))
}

func (m *ProxyClients) Count() int {
	return len(m.clients)
}
