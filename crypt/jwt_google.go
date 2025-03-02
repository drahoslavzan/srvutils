package crypt

type jwtGoogleParser struct{}

func JWTGoogleParser() jwtGoogleParser {
	go firebaseTP.FetchKeys()
	go gcpTP.FetchKeys()

	return jwtGoogleParser{}
}

// Verify and parse Firebase auth token. The subject is the user ID.
func (m jwtGoogleParser) ParseFirebaseToken(token string) (JWTClaims, error) {
	return firebaseTP.Parse(token)
}

// Verify and parse GCP auth token. The subject is the service account ID.
func (m jwtGoogleParser) ParseGCPToken(token string) (JWTClaims, error) {
	return gcpTP.Parse(token)
}

var (
	firebaseTP = newGSATokenParser("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	gcpTP      = newGSATokenParser("https://www.googleapis.com/oauth2/v1/certs")
)
