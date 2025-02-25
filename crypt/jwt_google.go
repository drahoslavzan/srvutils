package crypt

type JWTGoogleParser struct{}

func NewJWTGoogleParser() *JWTGoogleParser {
	go firebaseTP.FetchKeys()
	go gcpTP.FetchKeys()

	return &JWTGoogleParser{}
}

// Verify and parse Firebase auth token. The subject is the service account ID.
func (m *JWTGoogleParser) ParseFirebaseToken(token string) (JWTClaims, error) {
	return firebaseTP.Parse(token)
}

// Verify and parse GCP auth token. The subject is the service account ID.
func (m *JWTGoogleParser) ParseGCPToken(token string) (JWTClaims, error) {
	return gcpTP.Parse(token)
}

var firebaseTP = newGSATokenParser("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
var gcpTP = newGSATokenParser("https://www.googleapis.com/oauth2/v1/certs")
