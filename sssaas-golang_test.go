package sssaas

import (
	"testing"
	"fmt"
)

func TestFromYAML(t *testing.T) {
	output, err := FromYAML([]byte("database_password:\n    remote: ['https://api-1.sssaas.com/api/v0?key=a', 'https://api-2.sssaas.com/api/v0?key=b', 'https://api-3.sssaas.com/api/v0?key=c']\n    local: ./config/secrets.sssa\n    minimum: 4\n"), "database_password")
	if err != nil {
		t.Fatal(err)
	}
	if output != "test-pass" {
		t.Fatal("Unexpected output: " + output)
	}
}

func TestFromConfig(t *testing.T) {

}

func TestCreateCombine(t *testing.T) {
	apis := []string{"https://sssaas.cipherboy.com/api/v0/request/test", "http://sssaas.cipherboy.com:8765/api/v0/request/test"}
	tokens := []string{"sssaas-library-test-allowed", "sssaas-library-test-allowed"}
	shares := []string{"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=", "wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0="}
	for aid := range apis {
		api := []string{apis[aid]}
		token := []string{tokens[aid]}
		result, err := GetSecret(api, token, shares, 300)
		if err != nil || result != "test-pass" {
			t.Fatal("Fatal: creating and combining returned invalid data")
		}
	}
}
