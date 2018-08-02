package main

import (
	"encoding/json"
	"errors"
	//"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"proxy/mycrypt"
	"strings"
)

func prefilter(req *http.Request) (*http.Request, error) {
	//format URL
	newurl, err := format_url(req.URL)
	if err != nil {
		//log.Println("Unable to route by url : ", err.Error())
		return req, err
	}
	//format Request and test Authenticity
	var proxyReq *http.Request
	var valid bool
	if url_to_service(req.URL) == "auth" {
		proxyReq, err = format_request_auth(newurl, req)
	} else {
		if valid, err = request_authentication(req); valid {
			proxyReq, err = format_request_regular(newurl, req)
		}
	}
	if err != nil {
		//log.Println(err.Error())
		return req, err
	}

	return proxyReq, err
}
func url_to_service(url_in *url.URL) string {
	return uri_head(url_in.RequestURI())
}

func request_authentication(req *http.Request) (bool, error) {
	//Authentication
	authtoken := &Signin_Resp_to_Client{}
	authtoken_str := req.Header.Get("AuthToken")
	if authtoken_str == "" {
		//log.Println("AuthToken header does not exist")
		return false,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken header does not exist")
	}
	err := json.Unmarshal([]byte(authtoken_str), authtoken)
	if err != nil {
		//log.Println("AuthToken header not in form of authentication token")
		return false,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken header not in form of authentication token")
	}
	//log.Println("check hash : ", mycrypt.CreateMAC(authtoken.Username+authtoken.Role, auth_key))
	//log.Println("token hash : ", authtoken.Hash)
	validity := mycrypt.CheckMAC((authtoken.Username + authtoken.Role), authtoken.Hash, auth_key)
	if validity {
		return validity, nil
	} else {
		return validity,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken failed authentication")
	}
}

func format_request_regular(newurl *url.URL, req *http.Request) (*http.Request, error) {
	proxyReq, err := deepcopy_request_with_other_url(req, newurl)
	//Reform Request
	proxyReq.Header.Add("Service", req.URL.RequestURI())
	return proxyReq, err
}
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
func format_request_auth(newurl *url.URL, req *http.Request) (*http.Request, error) {
	proxyReq, err := deepcopy_request_with_other_url(req, newurl)
	proxyReq.Header.Add("Service", req.URL.RequestURI())
	return proxyReq, err
}

func format_url(url_in *url.URL) (*url.URL, error) {
	uri_input := url_in.RequestURI()
	inst_list, serv_err := service2instlist(uri_input)
	if serv_err != nil {
		//log.Println("Error while choosing service url :", serv_err.Error())
		return url_in, serv_err
	}

	url_out, inst_err := choose_instance(inst_list)
	if inst_err != nil {
		//log.Println("Error while choosing instance url :", inst_err.Error())
		return url_out, inst_err
	}
	url_out.Path = uri_tail(uri_input)
	return url_out, nil
}
func service2instlist(uri_input string) ([]string, error) {
	var inst_list []string
	requested_service := uri_head(uri_input)
	available_service_flag := false
	for _, cur_service := range MSA_ptr.Service {
		if requested_service == cur_service.Name {
			inst_list = cur_service.Instance
			available_service_flag = true
			break
		}
	}
	if available_service_flag == false {
		return nil,
			errors.New(RESOURCE_NONEXISTENT_ERROR + ERROR_STRING_SEPARATOR +
				"No service exists for particular uri")
	}
	return inst_list, nil
}
func choose_instance(instance_list []string) (*url.URL, error) {
	instance_chosen := instance_list[rand.Intn(len(instance_list))]
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
func uri_head(uri_in string) string {
	temp := strings.Split(uri_in, string(os.PathSeparator))[1]
	return temp
}
func uri_tail(uri_in string) string {
	parts := strings.Split(uri_in, string(os.PathSeparator))
	temp := parts[0] + string(os.PathSeparator) + strings.Join(parts[2:], string(os.PathSeparator))
	return temp
}
