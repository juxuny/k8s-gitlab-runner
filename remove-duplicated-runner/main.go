package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	host         string
	privateToken string
	runnerName   string
)

func exit(v string) {
	fmt.Println(v)
	os.Exit(-1)
}

func doGet(api string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		exit("invalid url: " + api)
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("PRIVATE-TOKEN", privateToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request faild: " + err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request faild: " + err.Error())
	}
	return json.Unmarshal(data, out)
}

func doDelete(api string, out interface{}) error {
	req, err := http.NewRequest(http.MethodDelete, api, nil)
	if err != nil {
		exit("invalid url: " + api)
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("PRIVATE-TOKEN", privateToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request faild: " + err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request faild: " + err.Error())
	}
	if len(data) == 0 && resp.StatusCode == http.StatusOK {
		return nil
	}
	return json.Unmarshal(data, out)
}

func parse() {
	flag.StringVar(&host, "h", "http://gitlab.com", "The host for gitlab")
	flag.StringVar(&privateToken, "t", "", "The AccessToken created by administrator")
	flag.StringVar(&runnerName, "n", "", "The name of new runner")
	flag.Parse()

	if host == "" {
		exit("-h host cannot be empty")
	}
	if privateToken == "" {
		exit("-t privateName cannot be empty")
	}
	if runnerName == "" {
		exit("-n runnerName cannot be empty")
	}
}

func genApi(path string, values ...url.Values) string {
	q := ""
	if len(values) > 0 {
		q = values[0].Encode()
	}
	if q != "" {
		return fmt.Sprintf("%s%s?%s", host, path, q)
	}
	return fmt.Sprintf("%s%s", host, path)
}

func main() {
	parse()

	// doGet all runners
	var runners Runners
	if err := doGet(genApi("/api/v4/runners/all"), &runners); err != nil {
		exit("doGet all runners failed: " + err.Error())
		return
	}
	runners = runners.Filter(func(r Runner, index int) bool {
		return r.Description == runnerName
	})
	if len(runners) == 0 {
		fmt.Println("runner not found:" + runnerName)
		return
	}
	fmt.Println("remove duplicated runner")
	for _, runner := range runners {
		var out interface{}
		if err := doDelete(genApi(fmt.Sprintf("/api/v4/runners/%d", runner.Id)), &out); err != nil {
			exit("delete duplicated runner failed:" + err.Error())
			return
		}
	}
	fmt.Println("success")
}
