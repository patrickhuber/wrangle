package values

type User struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type Password string

type Certificate struct {
	PublicKey            string `yaml:"public_key" json:"public_key"`
	PrivateKey           string `yaml:"private_key" json:"private_key"`
	CertificateAuthority string `yaml:"ca" json:"ca"`
}

type RSA struct {
	PublicKey  string `yaml:"public_key" json:"public_key"`
	PrivateKey string `yaml:"private_key" json:"private_key"`
}

type SSH struct {
	PublicKey  string `yaml:"public_key" json:"public_key"`
	PrivateKey string `yaml:"private_key" json:"private_key"`
}

type Structured map[string]interface{}
