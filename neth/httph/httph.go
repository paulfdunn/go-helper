// Package httph provides helper functions for the net/http package.
package httph

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/paulfdunn/go-helper/osh/runtimeh"
)

// Header key/value pairs are set when calling CollectURL.
type Header struct {
	Key   string
	Value string
}

// BodyUnmarshal - Unmarshals a request body (JSON) into an object. On any error the header
// is written with the appropriate http.Status; callers should not write header status.
func BodyUnmarshal(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return runtimeh.SourceInfoError("reading body", err)
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return runtimeh.SourceInfoError("unmarshal body", err)
	}

	return nil
}

// URLCollectionData - Functions that collect data from multiple URLs will return instance(s) of this
// structure, in order to allow association of URL, Byte (data), and errors.
type URLCollectionData struct {
	URL      string
	Bytes    []byte
	Response *http.Response
	Err      error
}

// CollectURL - Pass in a URL, request timeout, HTTP method to use, and get back
// the body of the request. HTTP method MUST be one of: [MethodGet, MethodHead]
func CollectURL(urlIn string, timeout time.Duration, method string, headers []Header) ([]byte, *http.Response, error) {
	var req *http.Request
	u, err := url.Parse(urlIn)
	if err != nil {
		return []byte{}, nil, runtimeh.SourceInfoError("error parsing urlIn", err)
	}

	var reqErr error
	switch method {
	case http.MethodGet:
		req, reqErr = http.NewRequest(http.MethodGet, u.String(), nil)
	case http.MethodHead:
		req, reqErr = http.NewRequest(http.MethodHead, u.String(), nil)
	default:
		return nil, nil, runtimeh.SourceInfoError("", fmt.Errorf("invalid method: %s", method))
	}

	if reqErr != nil {
		return nil, nil, runtimeh.SourceInfoError("Error creating http.Request", reqErr)
	}
	req.Header.Set("Connection", "close")
	for _, hdr := range headers {
		req.Header.Set(hdr.Key, hdr.Value)
	}
	req.Close = true

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			// This timeout is require in order to prevent "too many open file" errors.
			Timeout:   timeout,
			KeepAlive: timeout,
		}).Dial}
	client := http.Client{Timeout: timeout, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, resp, runtimeh.SourceInfoError("CollectURL client error", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("resp.Body.Close() error:%+v\n", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		return body, resp, err
	}

	return body, resp, err
}

// CollectURLs - Pass in a slice of URLs, request timeout, HTTP method to use, and
// get back a slice of URLCollectionData with results.
// The URLs are processed in parallel using threads number of parallel requests.
func CollectURLs(urls []string, timeout time.Duration, method string, threads int, headers []Header) []URLCollectionData {
	// Channel to feed work to the go routines
	tasks := make(chan string, threads)
	// Channel to return data from the workers.
	workerOut := make(chan URLCollectionData, len(urls))
	// Data to return to caller
	var returnData []URLCollectionData

	// Spawn threads number of workers
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(sendResult chan URLCollectionData) {
			for url := range tasks {
				b, resp, e := CollectURL(url, timeout, method, headers)
				sendResult <- URLCollectionData{url, b, resp, e}
			}
			wg.Done()
		}(workerOut)
	}

	for _, url := range urls {
		tasks <- url
	}
	close(tasks)

	wg.Wait()
	// Workers are done, all data should have already been returned.
	close(workerOut)
	for r := range workerOut {
		returnData = append(returnData, r)
	}

	return returnData
}

// RequestUsername will return the username of the request when using basic or digest
// authentication; if it can be determined.
func RequestUsername(r *http.Request) string {
	// r.Header["Authorization"] is a slice of strings. I.E.
	// Basic authentication.
	// "Authorization":[]string{"Basic YWRtaW46YWRtaW4="},
	// Digest authentication
	// r.Header["Authorization"] is a slice of strings. I.E.
	// "Authorization":[]string{"Digest username=\"admin\", realm=\"Western Digital Corporation\", nonce=\"AHYBbBIPrPRMzsDo\",...}
	for _, v := range r.Header["Authorization"] {
		splits := strings.Split(v, ",")
		for _, split := range splits {
			if strings.Contains(split, "username") {
				u := strings.Split(split, "=")
				if len(u) == 2 {
					return strings.Replace(u[1], `"`, ``, -1)
				}

				return ""
			} else if strings.Contains(split, "Basic ") {
				u := strings.Split(split, " ")
				if len(u) == 2 {
					user, err := base64.StdEncoding.DecodeString(u[1])
					if err != nil {
						return ""
					}
					userSplit := strings.Split(string(user), ":")
					return string(userSplit[0])
				}

				return ""
			}
		}
	}

	return ""
}
