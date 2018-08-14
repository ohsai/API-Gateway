package main

import (
	"net/http"
	"net/url"
)

func deepcopy_request_with_other_url(req *http.Request, newurl *url.URL) (*http.Request, error) {
	proxyReq, err := http.NewRequest(req.Method, newurl.String(), req.Body) // Copy Body
	if err != nil {
		//log.Println("Unable to create forward request! : ", err.Error())
		return proxyReq, err
	}

	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr) // Add forwardedness header

	for header, values := range req.Header { //Copy Headers
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	return proxyReq, err
}
func format_request_plain(newurl *url.URL, req *http.Request) (*http.Request, error) {
	proxyReq, err := deepcopy_request_with_other_url(req, newurl)
	//proxyReq.Header.Add("Service", req.URL.RequestURI())
	return proxyReq, err
}
