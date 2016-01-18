package sssaas

import (
	"testing"
)

func TestFromYAML(t *testing.T) {
	output, err := FromYAML([]byte("database_password:\n    remote: ['http://localhost:8765/api/v0/request/test?key=sssaas-library-test-allowed']\n    local: ./test/secrets.sssa\n    minimum: 4\n"), "database_password")
	if err != nil {
		t.Fatal(err)
	}
	if output != "test-pass" {
		t.Fatal("Unexpected output: " + output)
	}
}

func TestFromConfig(t *testing.T) {
	var c Config
	c.Remote = []string{"http://localhost:8765/api/v0/request/test?key=sssaas-library-test-allowed"}
	c.Shares = []string{"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=", "wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0="}
	c.Minimum = 4

	output, err := FromConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	if output != "test-pass" {
		t.Fatal("Unexpected output: " + output)
	}
}

func TestCreateCombine(t *testing.T) {
	apis := []string{"http://localhost:8765/api/v0/request/test?key=sssaas-library-test-allowed"}
	shares := []string{"j8-Y4_7CJvL8aHxc8WMMhP_K2TEsOkxIHb7hBcwIBOo=T5-EOvAlzGMogdPawv3oK88rrygYFza3KSki2q8WEgs=", "wGXxa_7FPFSVqdo26VKdgFxqVVWXNfwSDQyFmCh2e5w=8bTrIEs0e5FeiaXcIBaGwtGFxeyNtCG4R883tS3MsZ0="}
	for aid := range apis {
		api := []string{apis[aid]}

		result, err := GetSecret(api, shares, 4, 0)
		if err != nil || result != "test-pass" {
			t.Fatal("Fatal: creating and combining returned invalid data")
		}
	}
}
