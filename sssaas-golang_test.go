package sssaas

import (
	"testing"
)

func TestCreateCombine(t *testing.T) {
    api := []string{"https://sssaas.cipherboy.com:8444/api/v0/request/test", "https://sssaas.cipherboy.com:8444/api/v0/request/test", "https://sssaas.cipherboy.com:8444/api/v0/request/test"}
    token := []string{"sssaas-library-test-allowed", "sssaas-library-test-allowed", "sssaas-library-test-allowed"}
    shares := []string{"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=", "wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0="}
    result, err := GetSecret(api, token, shares, 300)
	if err != nil || result != "test-pass" {
		t.Fatal("Fatal: creating and combining returned invalid data")
	}
}
