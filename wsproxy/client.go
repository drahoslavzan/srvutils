package wsproxy

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

type (
	Except func(int, error) bool

	Client struct {
		url string
		cln *http.Client
		cfg *Config
	}
)

func NewHttpProxyClient(proxyURL string, cfg *Config) (*Client, error) {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("url parse of %s failed: %w", proxyURL, err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	cln := &Client{
		url: proxyURL,
		cln: client,
		cfg: cfg,
	}

	return cln, nil
}

func NewSocks5ProxyClient(proxyURL string) (*Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyURL, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("cannot create proxy request to %s: %w", proxyURL, err)
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
	}

	cln := &Client{
		url: proxyURL,
		cln: client,
	}

	return cln, nil
}

func (m *Client) String() string {
	parsed, err := url.Parse(m.url)
	if err != nil {
		return m.url
	}

	return parsed.Host
}

func (m *Client) Get(url string) (io.ReadCloser, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot create get request to %s: %w", url, err)
	}

	req.Header.Set("User-Agent", m.cfg.UserAgent)

	res, err := m.cln.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("get request to %s failed: %w", url, err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, res.StatusCode, fmt.Errorf("get request to %s failed with status %d", url, res.StatusCode)
	}

	return res.Body, res.StatusCode, nil
}

func (m *Client) Post(url string, data []byte, header http.Header) (io.ReadCloser, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("cannot create post request to %s: %w", url, err)
	}

	req.Header = header

	res, err := m.cln.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("post request to %s failed: %w", url, err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, res.StatusCode, fmt.Errorf("post request to %s failed with status %d", url, res.StatusCode)
	}

	return res.Body, res.StatusCode, nil
}

func (m *Client) RetryGet(url string, except Except, logger *zap.Logger) (io.ReadCloser, int, error) {
	return m.retry(url, m.Get, except, logger)
}

func (m *Client) RetryPost(url string, data []byte, header http.Header, except Except, logger *zap.Logger) (io.ReadCloser, int, error) {
	fn := func(url string) (io.ReadCloser, int, error) {
		return m.Post(url, data, header)
	}

	return m.retry(url, fn, except, logger)
}

func (m *Client) retry(url string, fn func(url string) (io.ReadCloser, int, error), except Except, logger *zap.Logger) (io.ReadCloser, int, error) {
	var lastErr error
	for i := 0; i < m.cfg.MaxRetries; i++ {
		body, status, err := fn(url)
		if err == nil {
			return body, status, nil
		}
		if except != nil && except(status, err) {
			return body, status, err
		}
		if status >= 400 && status < 500 {
			return body, status, err
		}

		logger.Warn("request failed", zap.String("url", url), zap.Error(err))

		backoff := time.Duration((1<<i)*500) * time.Millisecond  // exponential base
		jitter := time.Duration(rand.Int63n(int64(backoff / 2))) // up to 50% jitter
		totalBackoff := backoff + jitter

		lastErr = err
		time.Sleep(totalBackoff)
	}

	return nil, 0, fmt.Errorf("request failed after %d attempts, last error: %w", m.cfg.MaxRetries, lastErr)
}
