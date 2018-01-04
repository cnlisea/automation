package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Auth Type
const (
	NoAuth = iota
	BearerTokenAuth
)

const (
	HttpBodyNNil = iota
	HttpBodyTypeJson
	HttpBodyTypeXml
	HttpBodyTypePostForm
)

func HttpRequest(url string, body []byte, bodyType int, method string, authType int, token ...string) (int, []byte, error) {
	c := new(http.Client)

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if nil != err {
		return 0, nil, err
	}

	// set request head Content-Type
	switch bodyType {
	case HttpBodyTypeJson:
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	case HttpBodyTypeXml:
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	case HttpBodyTypePostForm:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// set auth
	switch authType {
	case BearerTokenAuth:
		if len(token) < 1 || len(token[0]) == 0 {
			return 0, nil, errors.New("bad parameter token(Bearer) ")
		}
		req.Header.Set("Authorization", "Bearer "+token[0])
	}

	// send request
	res, err := c.Do(req)
	if nil != err {
		return 0, nil, err
	}
	defer res.Body.Close()

	// read all body
	result, err := ioutil.ReadAll(res.Body)

	return res.StatusCode, result, err
}

func ValidMethod(method string) bool {
	var status bool

	switch strings.ToUpper(method) {
	case http.MethodGet:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodOptions:
		status = true
	default:
		status = false
	}

	return status
}
