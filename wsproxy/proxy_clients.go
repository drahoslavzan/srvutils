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
	if idx < 0 || idx > len(m.clients) {
		return nil, fmt.Errorf("invalid index: %d", idx)
	}

	return m.clients[idx], nil
}

func (m *ProxyClients) Random() (*Client, error) {
	sz := len(m.clients)
	if sz < 1 {
		return nil, fmt.Errorf("no clients available")
	}

	idx := rand.Intn(len(m.clients))
	return m.clients[idx], nil
}

func (m *ProxyClients) Count() int {
	return len(m.clients)
}
