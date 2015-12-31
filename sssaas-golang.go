package sssaas

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/SSSaaS/sssa-golang"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"os"
	"sort"
	"sync"
	"fmt"
	"time"
)

func FromYAML(content []byte, key string) (string, error) {
	var output map[string]map[string]interface{}
	err := yaml.Unmarshal(content, &output)
	if err != nil {
		return "", err
	}

	var obj Config
	if _, ok := output[key]["remote"]; ok {
		for i := range output[key]["remote"].([]interface{}) {
			obj.Remote = append(obj.Remote, output[key]["remote"].([]interface{})[i].(string))
		}
	}
	if _, ok := output[key]["shares"]; ok {
		for i := range output[key]["shares"].([]interface{}) {
			obj.Remote = append(obj.Remote, output[key]["shares"].([]interface{})[i].(string))
		}
	}
	if _, ok := output[key]["local"]; ok {
		obj.Local = output[key]["local"].(string)
	}
	if _, ok := output[key]["minimum"]; ok {
		obj.Minimum = output[key]["minimum"].(int)
	}
	if _, ok := output[key]["timeout"]; ok {
		obj.Minimum = output[key]["timeout"].(int)
	}

	return FromConfig(obj)
}

func FromConfig(obj Config) (string, error) {
	if obj.Local != "" {
		fh, err := os.Open(obj.Local)
		if err != nil {
			return "", err
		}
		defer fh.Close()

		r := bufio.NewReader(fh)
		line, err := r.ReadString('\n')

		for err == nil {
			last := len(line)-1
			share := line[:last]

			if sssa.IsValidShare(share) {
				obj.Shares = append(obj.Shares, share)
			}
			line, err = r.ReadString('\n')
		}

		if line != "" && sssa.IsValidShare(line) {
			obj.Shares = append(obj.Shares, line)
		}
	}

	return GetSecret(obj.Remote, obj.Shares, obj.Minimum, obj.Timeout)
}

func GetSecret(endpoints []string, shares []string, minimum int, timeout int) (string, error) {
	var results []string = shares
	var wg sync.WaitGroup
	var global_err []error
	var done bool = false

	if timeout == 0 {
		timeout = 300
	}

	duration := time.Duration(time.Duration(timeout) * time.Second)
	for i := range endpoints {
		wg.Add(1)
		go func() {
			client := &http.Client{
				Timeout: duration,
			}

			req, err := http.NewRequest("GET", endpoints[i], nil)
			if err != nil {
				global_err = append(global_err, err)

				if !done {
					wg.Done()
				}
				return
			}

			req.Header.Set("User-Agent", "sssaas-golang v0 v0.0.1")
			res, err := client.Do(req)
			if err != nil {
				global_err = append(global_err, err)

				if !done {
					wg.Done()
				}
				return
			}

			if res.StatusCode != 200 {
				global_err = append(global_err, errors.New(strconv.Itoa(res.StatusCode) + ":" + res.Status))

				if !done {
					wg.Done()
				}
			} else {
				defer res.Body.Close()
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					global_err = append(global_err, err)
				}

				current := response{}

				err = json.Unmarshal([]byte(data), &current)
				if err != nil {
					global_err = append(global_err, err)
				}

				for j := range current.SharedSecrets {
					results = append(results, current.SharedSecrets[j])
				}
			}

			results = removeDuplicates(results)

			if len(results) >= minimum && !done {
				done = true
				for _ = range endpoints {
					wg.Done()
				}
			}
			if !done {
				wg.Done()
			}
		}()
	}

	wg.Wait()

	results = removeDuplicates(results)

	if len(results) < minimum {
		fmt.Println("results")
		if len(global_err) >= 1 {
			return "", global_err[0]
		} else {
			return "", errors.New("Could not meet minimum shares!")
		}
	}

	return sssa.Combine(results), nil
}

func removeDuplicates(data []string) []string {
	results := data
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

	return results

}
