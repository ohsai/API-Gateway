package main

import (
	"io"
	//"log"
	"net/http"
)

func post_filter(proxyRes *http.Response, w http.ResponseWriter) error {
	err := format_response(proxyRes, w)
	if err != nil {
		return err
	}
	return nil
}
func format_response(proxyRes *http.Response, w http.ResponseWriter) error {
	return copy_response(proxyRes, w)
}
func copy_response(proxyRes *http.Response, w http.ResponseWriter) error {
	w.WriteHeader(proxyRes.StatusCode)
	// copy header
	copy_header(proxyRes.Header, w)
	//copy body
	if _, err := io.Copy(w, proxyRes.Body); err != nil {
		return err
	}
	proxyRes.Body.Close()
	return nil
}
func copy_header(response_headers http.Header, w http.ResponseWriter, decline_header ...string) {
	for cur_header, values := range response_headers {
		if !contains(decline_header, cur_header) {
			for _, value := range values {

				w.Header().Add(cur_header, value)
			}
		}
	}

}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
