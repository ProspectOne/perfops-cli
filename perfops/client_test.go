package perfops

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type roundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

func newTestClient(tr roundTripper) (*Client, error) {
	c := &http.Client{Transport: tr}
	return NewClient(WithHTTPClient(c))
}

type recordingTransport struct {
	req *http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return nil, errors.New("dummy impl")
}

type respondingTransport struct {
	resp *http.Response
}

func (t *respondingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.resp, nil
}

func dummyReq(method string) *http.Request {
	return &http.Request{Method: method, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}
}

func dummyResp(status int, method, body string) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(status) + " Meep meep",
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    dummyReq(method),
		Header: http.Header{
			"Connection":     {"close"},
			"Content-Length": {strconv.Itoa(len(body))},
		},
		Close:         true,
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}
