package crypt

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type FBTParser struct{}

func NewFBTParser() *FBTParser {
	fetchGSAPubKeys()
	return &FBTParser{}
}

// Verify and parse Firebase auth token. The subject is the user ID.
func (m *FBTParser) Parse(token string) (JWTClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid key id")
		}

		pk := getGSAPubKey(kid)
		if pk == nil {
			if err := fetchGSAPubKeys(); err != nil {
				return nil, err
			}
			pk = getGSAPubKey(kid)
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

func getGSAPubKey(id string) []byte {
	gsaPubKeysMut.RLock()
	defer gsaPubKeysMut.Unlock()

	return []byte(gsaPubKeys[id])
}

func fetchGSAPubKeys() error {
	gsaPubKeysMut.Lock()
	defer gsaPubKeysMut.Unlock()

	if gsaPubKeysLastFetch.Add(15 * time.Minute).After(time.Now()) {
		// recently fetched
		return nil
	}

	res, err := http.Get(gasPubKeysURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var keys map[string]string
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&keys); err != nil {
		return err
	}

	gsaPubKeys = keys
	gsaPubKeysLastFetch = time.Now()

	return nil
}

const gasPubKeysURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

var (
	gsaPubKeys          map[string]string
	gsaPubKeysMut       sync.RWMutex
	gsaPubKeysLastFetch time.Time
)
