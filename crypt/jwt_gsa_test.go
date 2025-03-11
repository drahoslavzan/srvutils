package crypt

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type mockKeyFetcher struct {
	keys map[string]string
}

func createTestToken(claims JWTClaims, pk string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))
	token.Header["kid"] = "key"

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	return token.SignedString(key)
}

func TestGSATokenParser_Parse(t *testing.T) {
	validClaims := JWTClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().UTC().Add(time.Hour).Unix(),
		"aud":  "my-firebase-app",
		"iss":  "https://securetoken.google.com/my-firebase-app",
	}

	expiredClaims := JWTClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().UTC().Add(-2 * time.Hour).Unix(),
		"exp":  time.Now().UTC().Add(-1 * time.Hour).Unix(),
		"aud":  "my-firebase-app",
		"iss":  "https://securetoken.google.com/my-firebase-app",
	}

	validToken, err := createTestToken(validClaims, privateKey)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	invalidToken, err := createTestToken(validClaims, invalidKey)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	expiredToken, err := createTestToken(expiredClaims, privateKey)
	if err != nil {
		t.Fatalf("Failed to create expired test token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		want    JWTClaims
		wantErr bool
	}{
		{
			name:  "valid token",
			token: validToken,
			want:  validClaims,
		},
		{
			name:    "invalid token",
			token:   invalidToken,
			wantErr: true,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFetcher := &mockKeyFetcher{keys: map[string]string{"key": publicKey}}
			parser := NewGSATokenParser("https://example.com", mockFetcher)

			got, err := parser.Parse(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Convert numeric values to int64 for comparison
				for k, v := range got {
					if f, ok := v.(float64); ok {
						got[k] = int64(f)
					}
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func (m *mockKeyFetcher) FetchKeys(url string) (map[string]string, error) {
	return m.keys, nil
}

func (m *mockKeyFetcher) LoadFromCache(url string) (map[string]string, error) {
	return m.keys, nil
}

const privateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCqvgrSCy+CQyQ9
9qD7YYvWr3QyRr7H29S8JxAinpNGVDfYwnKwo24jjVQsf79Y3coL6OJKbzTOIdCq
1PIi6TttP+PXo7IHjMTcs3J8CRBAKzVSFvlQkOR25uDtbKC3X3EzvOIhn4ofx/lL
E3WH2+9VQi/ruW5+CA5A8VZZ+5FJRdMXyj1Ighey9dC3cd8VydbPP6rvQcJvtGLt
a2n8BkoKc7DPnZvMoq/ncdzOn8TQNST2O/x0tNBSPV4w2Ny2sp7T5z28F9U3a1FE
PAMe8R2uTeXC8IwROtRV55n/8MyW2zYv8LXE1xhW9icgllRBJhCKHwre2L2m8DXP
wryMFDplAgMBAAECggEARl9gr1TL5okM0Xsqu6FnVPtozWYKMO6ivk/tXw9zzQte
Hzl5suMRPJb+B/85GwQCyCjax+oQ3hri2d491HTDuRIZsuG1uVXkN8DAYW4M5B3K
8sQkSvgFKhqbv0/D6ABu9G+X1lrev52Y2sAw82eLO901ShdZ+pkQYuT1fc3pgkSX
EKLx4UwivIQf5adgzmn96ZsxI5sn0ISM8NMkcroy4b/fXugfIrjwkQol/tDJxs77
1mxDfWAUFLJ4KB6t97BYHsicIzQoXsgquDS4o9hleLxggSs7q3X4uayWt6VS8cYQ
r6RANnqWur9D+ohfKhfDW9MYF+A49HTKewTuAO0I+wKBgQDsFKcjbeX+AIIv9Ojh
8Jvat6oIzto+pHWGUe2IEEgLm1EwzoPCWvI/9j9AlzfBxinoF/w8j4rCvqidR1Co
IYi/2fnFZ5QMJd3drLGcJNZrNBGuYpogcvPjtp0SmhcDDdmPBzDfeBbtqvGA9nmB
LSOhma58DUV+eJIdud+1hpJHYwKBgQC5JhSSaNQB6+VR+dtqINrnAP54cSPRdIhv
f95yEvRln7RKWFh9ey9VMTmrq3R1qcxITu4dz4MssqNfuF6OUjPUMJKHq9A/repy
7gNuQUjHO1lcppKai0OWf/SqDf1cJ7JbBmARpAGyTiJ3vIWGVIS3K7UTavlfY3Zt
kueXMZsVlwKBgQDTNfMGnob6xWe0Eg/cPFCj0Fe+g5n8G1TN9DRn2/Eo+S8dVFXL
J2S+VsfmOKP4qBrL+9F+OQnzC3J9K9V8ZmwbAKAWvYFVkPc7IQrR4J840B/VfX/Y
8h9DJhjHELbv049GCC/wbldNEPf54gl4yXKsXHsfnxwCd0p6b9Y2aIwscwKBgE6n
V7vu3onGbdAKZeTK+lOCP5ho7/9uEvvTBWvOk5aMZuniaA0+hJgbZlWAa+QEcy0Q
ouV1H1Ogu/jQ+RJa53uv+r+6BKjYuC2E4V44S7OfidHrTYJrwRWxW/3WHZjFoGY8
6hj0ZGgb+1aEdvDe/NZXsgACxNd6CHh6HPpE744XAoGAbEOi3IsY0nnhXGP+pGvq
3odGN160QWT9/PukAvEi9c3+SxYi8UtZHcTIFVo+Q8UJykqVt0QEFP3R1avqyTsD
AFAgKGh7op5IIu0csz9ydtahOcXqPSAE1c2DTSv6doaNvvmWUYoVZQ3LkerMl9Fr
WJMfmEiGCSOMrQCwRvOKI74=
-----END PRIVATE KEY-----
`

const publicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqr4K0gsvgkMkPfag+2GL
1q90Mka+x9vUvCcQIp6TRlQ32MJysKNuI41ULH+/WN3KC+jiSm80ziHQqtTyIuk7
bT/j16OyB4zE3LNyfAkQQCs1Uhb5UJDkdubg7Wygt19xM7ziIZ+KH8f5SxN1h9vv
VUIv67lufggOQPFWWfuRSUXTF8o9SIIXsvXQt3HfFcnWzz+q70HCb7Ri7Wtp/AZK
CnOwz52bzKKv53Hczp/E0DUk9jv8dLTQUj1eMNjctrKe0+c9vBfVN2tRRDwDHvEd
rk3lwvCMETrUVeeZ//DMlts2L/C1xNcYVvYnIJZUQSYQih8K3ti9pvA1z8K8jBQ6
ZQIDAQAB
-----END PUBLIC KEY-----
`

const invalidKey = `
-----BEGIN PRIVATE KEY-----
MIIEuwIBADANBgkqhkiG9w0BAQEFAASCBKUwggShAgEAAoIBAQCqhcODXxlrbMCL
Pa+9nWpH0fmWAajm7vhqdYSxgLEX27AtvVGdIPbXqkPWOVwvE4FJbz+9HEzrUs04
0O2iveOR6xVsdNvwK3RCk+jkYKT+rG41uEabcA41ehJqasbksYJThhf3w0rIRNtW
Pz9MMGk7IhaNPBVSRqyjdEX76aQE1wlLR/81dsZUIOSRJbJvCXLI3ZmS3mR4LwSn
i1RvHht5zjgvaQDXg2gZ4t4pCkwjLtdPeJn1XMxOkVgxoZxByoD0gS0wgc5mjMON
SLfZUH5uVm79GPIgVbrlRoyCIJpVzox6jnxk7hGyBV4yT6kzVjeYvkqxqEXjTjdz
0/+lb/ZpAgMBAAECgf8i2+Ab+v+MoQQQyDYk4l6CUBUK6qFHPug0MHyad18R1tct
LTEsmJCIMJuthLb+Pf6FWeNYcBtJVU6epBdFguFX8wwYynWA+LZ2Oxu2PrYmrxkt
4iiM81tJk1WFwPjnx2HdBoyKbwyFOY5HvfuH7QhAuQyN0rqJFz7eKzwjQ0qZslrv
MmGTntAcMADcJQQiFGuh7t/VBtHE1wrvQ/wo26HhdMVswYj8EK2EIMI1cZ0/n6qR
0dTVcXkDKSoJP0QESh2GfEjxoJGH/Y1kXIxdbs7f7PDiVGZUvDMHMbLyXq3YrDdy
AI5GBYt7E2Yp2xq+8ydeBCyaQmJ2sWqcQTaoRz0CgYEA8HZSyWXyfiricYd8dI9Q
UHY7HrIGd/5CIqw0O6WxzpjCz/9TZhJOhXlg76bVS2A5DBz1lrTZ0vl9eUeAuGIa
XeNGFsv8pB5MvAfi/NxSXFwg8rODHmQJjkX/dnsMxeVAfCuEApMC6X59ZryR19+a
lr5rgYhSm1Omsw302vvqpmMCgYEAtYqDKD5jZTMGdG+p2Qzo+VXMBqAa2ejW+emG
xs8VRF/L4xAnvKaM9PsjpR3Ad6r310vau5JT0Zm5U1pBp+HdJZqFZDDFgzBRjp5H
Cpf9SVS8GB9/PcCHirkSnK0PI0ezQkdy08iSXS7dWB+0WpzX1r3Fic2IRtIgxCL+
/FP0s8MCgYAkB9bHzrrTJOHhWQfQ/1htdgnNw6cse7C1OVBqT52g80rdl8iLVtrl
LRbVUg5LyRNDOWOjPV4WOsQOVCR5fFmvD8sEx3QHs3KUCip88RZ2OGfHdhaDi0HT
S7HHsxBq5rsO4AZbzGN3UTjBGChSTHMBe27obeDS3WnxEnpKBc1XOQKBgQCpOP5d
axp34QJpUxU+MByTHvjaTC/7ZGHP/3EUrUAjxjBl7k88OPw3+EoXxg38/q+cTycL
pbDgNq1cF0wQVCgyv0EMTbIvQcEkckHCjD8cNhJHYkTXqTovg6jnxyHPPyzH4ZYV
+GcG2YKWfKc/t+gyUh9q/t3DNmg4rG6HgzMcxQKBgHCnsI7AAV1HoQ2PlpjePWse
ej5K0fE69J7nG5P+32tCYFzaUC1hMCxTw62YmTps5AbbYbwc1VTzX+xGOxpN/Bej
lTTIN3ZSNnVwEgbdntP8Hy9/eIxUOwIATTyKn5K/ogT86FIanH26wNfbHbXb2usv
Xy6twLiLom48SXtUeNcm
-----END PRIVATE KEY-----
`
