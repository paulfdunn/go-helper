package httph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testUser     = "testUser"
	testPassword = "password"
)

var (
	reqUser string
)

func TestBodyUnmarshal(t *testing.T) {
	id := "this is a test struct"
	type testStruct struct {
		ID string
	}
	tsIn := testStruct{ID: id}
	tsOut := testStruct{}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		BodyUnmarshal(w, r, &tsOut)
	}))
	defer testServer.Close()

	b, err := json.Marshal(tsIn)
	if err != nil {
		t.Errorf("marshal error: %v", err)
		return
	}
	resp, err := http.Post(testServer.URL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("POST error: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code: %d", resp.StatusCode)
		return
	}
	if tsOut.ID != id {
		t.Errorf("incorrect tsOut: %+v", tsOut)
		return
	}

	resp, err = http.Post(testServer.URL, "application/json", bytes.NewBuffer(b[1:]))
	if err != nil {
		t.Errorf("POST error: %v", err)
		return
	}
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status code: %d", resp.StatusCode)
		return
	}
}

func TestCollectURL(t *testing.T) {
	returnString := `{"value":"test CollectURL"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(returnString))
		}
	}))
	defer server.Close()

	// Only HEAD and GET are supported.
	_, _, errDelete := CollectURL(server.URL, 1*time.Second, http.MethodDelete)
	if errDelete == nil {
		t.Errorf("CollectURL expected to return error on invalid method, but no error returned.")
		return
	}

	value, response, err := CollectURL(server.URL, 1*time.Second, http.MethodGet)
	if err != nil {
		t.Errorf("CollectURL returned non-nil error: %v", err)
		return
	}
	if string(value) != returnString {
		t.Errorf("Expected %s, got %s", returnString, value)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("incorrect status, expected %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestCollectURLs(t *testing.T) {
	returnString := `{"value":"test CollectURLs"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(returnString))
		}
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL}
	ucds := CollectURLs(urls, 1*time.Second, http.MethodGet, 2)
	if len(ucds) != len(urls) {
		t.Errorf("Incorrect number of URLCollectionData items returned, expected %d, got %d", len(ucds), len(urls))
		return
	}
	for _, ucd := range ucds {
		if ucd.Err != nil {
			t.Errorf("CollectURLs returned non-nil error: %v", ucd.Err)
			return
		}
		if string(ucd.Bytes) != returnString {
			t.Errorf("Expected %s, got %s", returnString, ucd.Bytes)
		}
		if ucd.Response.StatusCode != http.StatusOK {
			t.Errorf("incorrect status, expected %d, got %d", http.StatusOK, ucd.Response.StatusCode)
		}
	}
}

func TestRequestUsername(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandlerFuncUser))
	defer ts.Close()

	handler := http.HandlerFunc(testHandlerFuncUser)
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://TestRequestUsername", nil)
	req.SetBasicAuth(testUser, testPassword)
	handler.ServeHTTP(resp, req)

	if reqUser != "testUser" {
		t.Errorf("User was not correct, reqUser:%+v", reqUser)
	} else {
		//fmt.Printf("User was:%+v", reqUser)
	}
}

func testHandlerFuncUser(w http.ResponseWriter, r *http.Request) {
	reqUser = RequestUsername(r)
}
