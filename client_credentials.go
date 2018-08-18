package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"flag"
	"os"
	"bytes"
	"mime/multipart"
)

func main() {
	host, cid, pw := parseFlags()
	req := createClientCredentialsRequest(host, cid, pw)
	body := sendRequest(req)
	token := extractAccessToken(body)
	fmt.Println(token)
}

func parseFlags() (*string, *string, *string) {
	host := flag.String("host", "http://localhost:9094/token", "client secret")
	cid := flag.String("cid", "foo", "client id")
	pw := flag.String("cpw", "bar", "client secret")
	flag.Parse()
	return host, cid, pw
}

func createClientCredentialsRequest(host *string, cid *string, pw *string) (*http.Request) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("grant_type", "client_credentials")
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", *host, bodyBuf)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", contentType)
	req.SetBasicAuth(*cid, *pw)
	return req
}

func sendRequest(req *http.Request) []byte {
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
		fmt.Fprintf(os.Stderr, "Error fetching access token:\n")
		fmt.Fprintf(os.Stderr, string(body))
		os.Exit(1)
	}

	return body
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	Expiry      int    `json:"expires_in"`
}

func extractAccessToken(body []byte) string {
	var atr AccessTokenResponse
	err := json.Unmarshal(body, &atr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse response: \"")
		fmt.Fprintf(os.Stderr, "%s", err)
		fmt.Fprintf(os.Stderr, "\"\n")
		fmt.Fprintf(os.Stderr, "Response was:\n")
		fmt.Fprintf(os.Stderr, string(body)[:200])
		fmt.Fprintf(os.Stderr, "...")
		os.Exit(1)
	}
	return atr.AccessToken
}
