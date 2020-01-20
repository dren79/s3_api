package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Auth(t *testing.T) {
	responseStatus, responseBody := testHelperAuth("admin", "adminpass", t)
	assert.Equal(t, http.StatusOK, responseStatus, "should return status 200")
	assert.Equal(t, "Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJyb2xlIjoiRGV2T3BzIn0.JKGVHubFfBBk2w6_T7UDcV7hNh5CRyC-p87JenqnOtE", string(responseBody), "token not created")
}


func testHelperAuth(Name string, Password string, t *testing.T) (int, string) {
	r := ConfigureRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	reqURL := ts.URL + "/api/v1/auth"

	var data = make(url.Values)
	data.Set("Name", Name)
	data.Set("Password", Password)

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		t.Error("Error making the request")
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode, string(responseBody)
}
