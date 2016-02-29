package apiversion

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var httpClient *http.Client
var testServer *httptest.Server

func testVendorCall(acceptHeader, responseExpected string, t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL, nil)
	req.Header.Add("Accept", acceptHeader)

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Error("Unexpected err:", err)
		t.FailNow()
	}

	if resp.StatusCode != 200 {
		t.Errorf("%s Expected 200 - Got %d \n", acceptHeader, resp.StatusCode)
		t.FailNow()
	}

	value, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		t.Error("Unexpected error when reading response:", err)
		t.Fail()
	}

	if string(value) != responseExpected {
		t.Errorf("%s Response unexpected: %s \n", acceptHeader, responseExpected)
	}
}
func TestRoutingSucceed(t *testing.T) {

	muxV1 := http.NewServeMux()
	muxV1.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "version 1")

	}))

	muxV2 := http.NewServeMux()
	muxV2.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "version 2")
	}))

	apiV1 := NewVersion("v1", muxV1)
	apiV2 := NewVersion("v2", muxV2)

	vendor, err := NewVendorMiddleware("mybusiness.com", apiV1, apiV2)
	if err != nil {
		t.Error("Unexpected initialization")
		t.FailNow()
	}

	// Server
	testServer = httptest.NewServer(vendor)
	defer testServer.Close()

	// Client
	httpClient = &http.Client{}

	testVendorCall("application/vnd.mybusiness.com", "version 2", t)
	testVendorCall("application/vnd.mybusiness.com-v2", "version 2", t)
	testVendorCall("application/vnd.mybusiness.com-v1", "version 1", t)
}
