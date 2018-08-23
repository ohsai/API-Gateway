package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"proxy/mycrypt"
)

func post_filter(proxyRes *http.Response, w http.ResponseWriter) error {
	req_serv := proxyRes.Header.Get("Service")
	if req_serv == "" {
		return errors.New("Service header not appended by proxy")
	}
	var err error
	if uri_head(req_serv) == "auth" {
		err = format_response_auth(proxyRes, w)
	} else {
		err = format_response(proxyRes, w)
	}
	if err != nil {
		return err
	}
	return nil
}
func format_response_auth(proxyRes *http.Response, w http.ResponseWriter) error {
	req_serv := proxyRes.Header.Get("Service")
	//if signin, create Auth token
	if uri_head(uri_tail(req_serv)) == "signin" && proxyRes.StatusCode == 200 {
		b, err := create_auth_token(proxyRes.Body)
		if err != nil {
			return err
		}
		copy_header(proxyRes.Header, w, "Content-Length")
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(b) //write auth token to response body
		if err != nil {
			return err
		}
		return nil
	} else {
		return copy_response(proxyRes, w)
	}
}

//Create json type auth token
func create_auth_token(Body io.ReadCloser) ([]byte, error) {
	temp := &Signin_Resp_from_Auth{}
	jsonparseerr := json.NewDecoder(Body).Decode(temp)
	if jsonparseerr != nil {
		return nil, jsonparseerr
	}
	authresp := Signin_Resp_to_Client{
		Username: temp.Username,
		Role:     temp.Role,
		Hash:     mycrypt.CreateMAC(temp.Username+temp.Role, Config_ptr.Auth_key),
	}
	b, err := json.Marshal(authresp)
	if err != nil {
		return nil, err
	}
	return b, err
}
func format_response(proxyRes *http.Response, w http.ResponseWriter) error {
	return copy_response(proxyRes, w)
}
func copy_response(proxyRes *http.Response, w http.ResponseWriter) error {
	w.WriteHeader(proxyRes.StatusCode)
	copy_header(proxyRes.Header, w)
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
