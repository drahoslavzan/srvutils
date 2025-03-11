package crypt

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
)

type httpKeyFetcher struct {
	cacheDir string
}

func NewHTTPKeyFetcher(cacheDir string) *httpKeyFetcher {
	if cacheDir != "" {
		os.MkdirAll(cacheDir, 0755)
	}

	return &httpKeyFetcher{
		cacheDir: cacheDir,
	}
}

func (m *httpKeyFetcher) FetchKeys(url string) (map[string]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var keys map[string]string
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&keys); err != nil {
		return nil, err
	}

	if m.cacheDir != "" {
		m.saveToCache(url, keys)
	}

	return keys, nil
}

func (m *httpKeyFetcher) LoadFromCache(url string) (map[string]string, error) {
	if m.cacheDir == "" {
		return nil, errors.New("caching disabled")
	}

	filename := m.getCacheFilename(url)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var keys map[string]string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func (m *httpKeyFetcher) getCacheFilename(url string) string {
	// Create a hash of the URL to use as filename
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return filepath.Join(m.cacheDir, hash+".json")
}

func (m *httpKeyFetcher) saveToCache(url string, keys map[string]string) error {
	if m.cacheDir == "" {
		return nil
	}

	filename := m.getCacheFilename(url)
	data, err := json.Marshal(keys)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
