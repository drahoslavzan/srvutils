package crypt

import (
	"errors"
	"fmt"

	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	KeyFetcher interface {
		FetchKeys(url string) (map[string]string, error)
		LoadFromCache(url string) (map[string]string, error)
	}

	GSATokenParser struct {
		url       string
		fetcher   KeyFetcher
		pubKeys   map[string]string
		mut       sync.RWMutex
		lastFetch time.Time
	}
)

const fetchDelay = 5 * time.Minute

func NewGSATokenParser(url string, fetcher KeyFetcher) *GSATokenParser {
	fetcher.LoadFromCache(url)

	return &GSATokenParser{
		url:     url,
		fetcher: fetcher,
	}
}

func (m *GSATokenParser) Parse(token string) (JWTClaims, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	}

	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid key id")
		}

		pk := m.keyByID(kid)
		if len(pk) < 1 {
			if err := m.fetchKeys(); err != nil {
				return nil, err
			}
			pk = m.keyByID(kid)
			if pk == nil {
				return nil, errors.New("public key not found")
			}
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pk))
		if err != nil {
			return nil, fmt.Errorf("invalid public key: %s", pk)
		}

		return pubKey, nil
	}, opts...)
	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return JWTClaims(claims), nil
	}

	return nil, errors.New("invalid token")
}

func (m *GSATokenParser) keyByID(id string) []byte {
	m.mut.RLock()
	defer m.mut.RUnlock()

	return []byte(m.pubKeys[id])
}

func (m *GSATokenParser) fetchKeys() error {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m.lastFetch.Add(fetchDelay).After(time.Now()) {
		// recently fetched
		return nil
	}

	keys, err := m.fetcher.FetchKeys(m.url)
	if err != nil {
		return err
	}

	m.pubKeys = keys
	m.lastFetch = time.Now()

	return nil
}
