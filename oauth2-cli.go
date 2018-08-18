package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"net/url"
)

const (
	clientGrant   = "client_credentials"
	passwordGrant = "password"
)

const (
	error   = 1
	success = 0
)
const usage = `Usage of oauth2-cli:

oauth2-cli [opts]

	oauth2-cli retrieves an OAuth2 access token using client or password grant

`

func main() {
	os.Exit(Run())
}

func Run() int {
	host, cid, cpw, uid, upw, typ := parseFlags()

	req, e := createRequest(host, cid, cpw, uid, upw, typ)
	if e != success {
		return e
	}

	body, e := sendRequest(req)
	if e != success {
		return e
	}

	token, e := extractAccessToken(body)
	if e != success {
		return e
	}

	fmt.Println(token)
	return success
}

func parseFlags() (*string, *string, *string, *string, *string, *string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
	host := flag.String("host", "http://localhost:9094/token", "authorization server url")
	cid := flag.String("cid", "foo", "client id")
	cpw := flag.String("cpw", "bar", "client secret")
	uid := flag.String("uid", "fizz", "end user id")
	upw := flag.String("upw", "buzz", "end user secret")
	typ := flag.String("typ", clientGrant, "grant type, can be "+clientGrant+" or "+passwordGrant)

	flag.Parse()
	return host, cid, cpw, uid, upw, typ
}

func createRequest(host *string, cid *string, cpw *string, uid *string, upw *string, typ *string) (*http.Request, int) {
	data := url.Values{}
	if *typ == clientGrant {
		data.Add("grant_type", clientGrant)
	} else if *typ == passwordGrant {
		data.Add("grant_type", passwordGrant)
		data.Add("username", *uid)
		data.Add("password", *upw)
	} else {
		fmt.Println("Unknown grant type (typ parameter was: '" + *typ + "')")
		return nil, error
	}
	return formDataRequestWithBody(host, cid, cpw, data)
}

func formDataRequestWithBody(host *string, cid *string, cpw *string, data url.Values) (*http.Request, int) {
	req, err := http.NewRequest("POST", *host, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, error
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(*cid, *cpw)

	return req, success
}

func sendRequest(req *http.Request) ([]byte, int) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error response from token endpoint (HTTP Status %d):\n", res.StatusCode)
		fmt.Fprintf(os.Stderr, string(body))
		return nil, error
	}

	return body, success
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	Expiry      int    `json:"expires_in"`
}

func extractAccessToken(body []byte) (string, int) {
	var atr AccessTokenResponse
	err := json.Unmarshal(body, &atr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse response: \"")
		fmt.Fprintf(os.Stderr, "%s", err)
		fmt.Fprintf(os.Stderr, "\"\n")
		fmt.Fprintf(os.Stderr, "Response was:\n")
		fmt.Fprintf(os.Stderr, string(body)[:200])
		fmt.Fprintf(os.Stderr, "...")
		return "", error
	}
	return atr.AccessToken, success
}
