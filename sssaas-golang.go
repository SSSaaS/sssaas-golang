package sssaas

import (
	"encoding/json"
	"errors"
	sssa "github.com/SSSaaS/sssa-golang"
	"io/ioutil"
	"net/http"
	"sort"
	"bufio"
	"strconv"
    "os"
	"sync"
	"time"
	"gopkg.in/yaml.v2"
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
		for err != nil {
			if sssa.IsValidShare(line) {
                obj.Shares = append(obj.Shares, line)
            }
		}

        if line != "" && sssa.IsValidShare(line) {
            obj.Shares = append(obj.Shares, line)
        }
	}

    return GetSecret(obj.Remote, obj.Shares, obj.Timeout)
}

func GetSecret(serveruris []string, shares []string, timeout int) (string, error) {
	var results []string = shares
	var wg sync.WaitGroup

	var has_err = false
	var err_mesg = ""

    if timeout <= 0 {
        timeout = 300
    }

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

	results = RemoveDuplicates(results)

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
