package main

import (
	"encoding/json"
	"io"
	//"log"
	"net/http"
	"proxy/mycrypt"
)

func post_filter(proxyRes *http.Response, w http.ResponseWriter) error {
	req_serv := proxyRes.Header.Get("Service")
	//log.Println(proxyRes.Status)
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
	if uri_head(uri_tail(req_serv)) == "signin" && proxyRes.StatusCode == 200 {
		//if signin
		// Create Auth token
		b, err := create_auth_token(proxyRes.Body)
		if err != nil {
			return err
		}
		//copy header
		copy_header(proxyRes.Header, w, "Content-Length")
		//Supplementary Header
		w.Header().Set("Content-Type", "application/json")
		//Deep copy body
		_, err = w.Write(b)
		if err != nil {
			return err
		}
		return nil
	} else {
		//if signup or wrong responseresponse
		return copy_response(proxyRes, w)
	}
}
func create_auth_token(Body io.ReadCloser) ([]byte, error) {
	temp := &Signin_Resp_from_Auth{}
	jsonparseerr := json.NewDecoder(Body).Decode(temp)
	if jsonparseerr != nil {
		return nil, jsonparseerr
	}
	authresp := Signin_Resp_to_Client{
		Username: temp.Username,
		Role:     temp.Role,
		Hash:     mycrypt.CreateMAC(temp.Username+temp.Role, auth_key),
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
