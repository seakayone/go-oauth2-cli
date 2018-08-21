package main

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequestClientGrant(t *testing.T) {
	host := "https://example.com"
	cid := "cid"
	cpw := "cpw"
	uid := "uid"
	upw := "upw"
	typ := clientGrant
	request, err := createRequest(&host, &cid, &cpw, &uid, &upw, &typ)

	assert.NotNil(t, request, "Request was nil")
	assertNoError(t, err)
	assertAuthorizationHeader(t, request)
	assertContentTypeHeader(t, request)
}

func assertNoError(t *testing.T, err int) bool {
	return assert.EqualValues(t, 0, err, "Error should be zero")
}

func assertAuthorizationHeader(t *testing.T, request *http.Request) bool {
	return assertRequestHeaderValue(t, request, "Authorization", "Basic Y2lkOmNwdw==")
}

func assertContentTypeHeader(t *testing.T, request *http.Request) bool {
	return assertRequestHeaderValue(t, request, "Content-Type", "application/x-www-form-urlencoded")
}

func assertRequestHeaderValue(t *testing.T, actual *http.Request, headerName string, headerValueExpected string) bool {
	return assert.EqualValues(t, headerValueExpected, actual.Header.Get(headerName), "Header %q not set correctly", headerName)
}
func TestCreateRequestPasswordGrant(t *testing.T) {
	host := "https://example.com"
	cid := "cid"
	cpw := "cpw"
	uid := "uid"
	upw := "upw"
	typ := passwordGrant
	request, err := createRequest(&host, &cid, &cpw, &uid, &upw, &typ)

	assert.NotNil(t, request, "Request was nil")
	assertNoError(t, err)
	assertAuthorizationHeader(t, request)
	assertContentTypeHeader(t, request)
}
