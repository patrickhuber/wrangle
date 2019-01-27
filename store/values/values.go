package values

type User struct {
	Username string
	Password string
}

type Password string

type Certificate struct {
	PublicKey            string
	PrivateKey           string
	CertificateAuthority string
}

type RSA struct {
	PublicKey  string
	PrivateKey string
}

type SSH struct {
	PublicKey  string
	PrivateKey string
}

type Structured map[string]interface{}
