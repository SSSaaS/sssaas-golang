package sssaas

import (
	"encoding/json"
	"errors"
	"github.com/SSSaaS/sssa-golang"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Shares struct {
	SharedSecrets []string `json:"sharedSecrets"`
}

func GetSecret(serveruris []string, tokens []string, shares []string, timeout int) (string, error) {

	var results []string = shares
	var wg sync.WaitGroup

	var has_err = false
	var err_mesg = ""

	duration := time.Duration(time.Duration(timeout) * time.Second)

	for i := range serveruris {
		wg.Add(1)
		go func() {
			client := &http.Client{
				Timeout: duration,
			}

			req, _ := http.NewRequest("GET", serveruris[i]+"?key="+tokens[i], nil)
			req.Header.Set("User-Agent", "sssaas-golang v0 v0.0.1")
			res, _ := client.Do(req)

			if res.StatusCode != 200 {

				has_err = true
				err_mesg += strconv.Itoa(res.StatusCode) + ": " + res.Status + "; "
			} else {
				defer res.Body.Close()
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					has_err = true
					err_mesg += err.Error() + ": "
				}

				current := Shares{}

				json.Unmarshal([]byte(data), &current)

				for j := range current.SharedSecrets {
					results = append(results, current.SharedSecrets[j])
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if has_err {
		return "", errors.New(err_mesg)
	}

	sort.Strings(results)
	clen := len(results)
	i := 0

	for i < clen-1 {
		if results[i] == results[i+1] {
			results = append(results[:i], results[i+1:]...)
			clen = len(results)
		} else {
			i += 1
		}
	}

	return sssa.Combine(results), nil
}
