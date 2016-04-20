package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	data map[string]interface{}
)

type Response struct {
	Path    string
	Method  string
	Code    int
	Body    []byte
	Err     error
	Cookies map[string]string
}

func SlurpFromURL(t string) []byte {
	r, err := http.Get(t)
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	return b
}

func SlurpFromFile(filePath string) []byte {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	return data
}

func Post(u *url.URL, data []byte, key string) (*Response, error) {
	req, err := http.NewRequest("POST", u.String(),
		bytes.NewReader(data))
	if err != nil {
		return &Response{Path: u.Path, Method: "POST", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)
	fmt.Println(req.URL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "POST", Err: err}, err
	}

	r, err := ProcessResponse(res)
	r.Path = u.Path
	r.Method = "POST"
	r.Err = err
	return r, err
}

func Get(u *url.URL, key string) (*Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return &Response{Path: u.Path, Method: "GET", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)
	fmt.Println(req.URL)

	client := http.Client{Timeout: (2999 * time.Millisecond)}

	res, err := client.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "GET", Err: err}, err
	}

	r, err := ProcessResponse(res)
	r.Path = u.Path
	r.Method = "GET"
	r.Err = err
	return r, err
}

func Put(u *url.URL, content []byte, key string) (*Response, error) {
	req, err := http.NewRequest("PUT", u.String(), bytes.NewReader(content))
	if err != nil {
		return &Response{Path: u.Path, Method: "PUT", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "PUT", Err: err}, err
	}

	r, err := ProcessResponse(res)
	r.Path = u.Path
	r.Method = "PUT"
	r.Err = err
	return r, err
}

func Del(u *url.URL, key string) (*Response, error) {
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return &Response{Path: u.Path, Method: "DELETE", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "DELETE", Err: err}, err
	}

	r, err := ProcessResponse(res)
	r.Path = u.Path
	r.Method = "DELETE"
	r.Err = err
	return r, err
}

func ProcessResponse(res *http.Response) (r *Response, err error) {
	ret := &Response{
		Code: res.StatusCode,
	}

	ret.Body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
