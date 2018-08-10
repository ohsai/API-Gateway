package main

import (
	"net/url"
	"os"
	"strings"
)

func uri_head(uri_in string) string {
	temp := strings.Split(uri_in, string(os.PathSeparator))[1]
	return temp
}
func uri_tail(uri_in string) string {
	parts := strings.Split(uri_in, string(os.PathSeparator))
	temp := parts[0] + string(os.PathSeparator) + strings.Join(parts[2:], string(os.PathSeparator))
	return temp
}
func url_to_service(url_in *url.URL) string {
	return uri_head(url_in.RequestURI())
}
