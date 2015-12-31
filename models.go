package sssaas

type response struct {
	SharedSecrets []string `json:"sharedSecrets"`
}

type Config struct {
	Remote  []string
	Local   string
	Shares  []string
	Minimum int
	Timeout int
}
