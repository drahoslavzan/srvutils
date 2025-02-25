package crypt

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	gsaTokenParser struct {
		url       string
		pubKeys   map[string]string
		mut       sync.RWMutex
		lastFetch time.Time
	}
)

func newGSATokenParser(url string) *gsaTokenParser {
	return &gsaTokenParser{
		url: url,
	}
}

func (m *gsaTokenParser) KeyByID(id string) []byte {
	m.mut.RLock()
	defer m.mut.RUnlock()

	return []byte(m.pubKeys[id])
}

func (m *gsaTokenParser) FetchKeys() error {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m.lastFetch.Add(fetchDelay).After(time.Now()) {
		// recently fetched
		return nil
	}

	res, err := http.Get(m.url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var keys map[string]string
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&keys); err != nil {
		return err
	}

	m.pubKeys = keys
	m.lastFetch = time.Now()

	return nil
}

func (m *gsaTokenParser) Parse(token string) (JWTClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid key id")
		}

		pk := m.KeyByID(kid)
		if pk == nil {
			if err := m.FetchKeys(); err != nil {
				return nil, err
			}
			pk = m.KeyByID(kid)
			if pk == nil {
				return nil, errors.New("public key not found")
			}
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pk))
		if err != nil {
			return nil, errors.New("invalid public key")
		}

		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return JWTClaims(claims), nil
	}

	return nil, errors.New("invalid token")
}

const fetchDelay = 5 * time.Minute
