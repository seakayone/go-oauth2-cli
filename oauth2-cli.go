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

const (
	clientGrant   = "client_credentials"
	passwordGrant = "password"
)

func main() {
	host, cid, cpw, uid, upw, typ := parseFlags()
	req := createRequest(host, cid, cpw, uid, upw, typ)
	body := sendRequest(req)
	token := extractAccessToken(body)
	fmt.Println(token)
}

func parseFlags() (*string, *string, *string, *string, *string, *string) {
	host := flag.String("host", "http://localhost:9094/token", "client secret")
	cid := flag.String("cid", "foo", "client id")
	cpw := flag.String("cpw", "bar", "client secret")
	uid := flag.String("uid", "fizz", "end user id")
	upw := flag.String("upw", "buzz", "end user secret")
	typ := flag.String("typ", clientGrant, "grant type, can be "+clientGrant+" or "+passwordGrant)

	flag.Parse()
	return host, cid, cpw, uid, upw, typ
}

func createRequest(host *string, cid *string, cpw *string, uid *string, upw *string, typ *string) (*http.Request) {
	var fieldWriter BodyFieldWriter
	if *typ == clientGrant {
		fieldWriter = clientGrantBodyWriter()
	} else if *typ == passwordGrant {
		fieldWriter = passwordGrantBodyWriter(uid, upw)
	} else {
		fmt.Println("Unknown grant type (typ parameter was: '" + *typ + "')")
		os.Exit(1)
	}
	return multiPartFormDataRequestWithBody(host, cid, cpw, fieldWriter)
}

type BodyFieldWriter func(bodyWriter *multipart.Writer)

func clientGrantBodyWriter() BodyFieldWriter {
	return func(bodyWriter *multipart.Writer) {
		bodyWriter.WriteField("grant_type", clientGrant)
	}
}

func passwordGrantBodyWriter(uid *string, upw *string) BodyFieldWriter {
	return func(bodyWriter *multipart.Writer) {
		bodyWriter.WriteField("grant_type", passwordGrant)
		bodyWriter.WriteField("username", *uid)
		bodyWriter.WriteField("password", *upw)
	}
}

func multiPartFormDataRequestWithBody(host *string, cid *string, cpw *string, bodyFieldWriter BodyFieldWriter) (*http.Request) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyFieldWriter(bodyWriter)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	req, err := http.NewRequest("POST", *host, bodyBuf)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", contentType)
	req.SetBasicAuth(*cid, *cpw)
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
