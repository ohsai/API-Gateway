package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func prefilter(req *http.Request) (*http.Request, error) {
	// Region check
	req_region_code := req.Header.Get("Accept-Language")
	load_balancer_info := []string{strings.Split(req.RemoteAddr, ":")[0], req.URL.Path}
	if req_region_code == "" || req_region_code == RS_ptr.Cur_region_code { //on right region
		newurl, err := format_url(req.URL, load_balancer_info)
		if err != nil {
			return req, err
		}
		var proxyReq *http.Request
		var valid bool
		if url_to_service(req.URL) == "auth" {
			proxyReq, err = format_request_plain(newurl, req)
		} else {
			if valid, err = request_authentication(req); valid {
				proxyReq, err = format_request_after_auth(newurl, req)
			}
		}
		if err != nil {
			return req, err
		}
		return proxyReq, err

	} else { //Wrong region
		AZ_list, region_err := region2AZlist(req_region_code) //Find closest region
		if region_err != nil {
			return req, region_err
		}
		url_out, inst_err := choose_instance(AZ_list, load_balancer_info) //Choose one of region elbs
		if inst_err != nil {
			return req, inst_err
		}
		url_out.Path = req.URL.Path
		proxyReq, err := format_request_plain(url_out, req)
		if err != nil {
			return req, err
		}
		return proxyReq, err
	}
}
func format_url(url_in *url.URL, load_balancer_info []string) (*url.URL, error) {
	uri_input := url_in.RequestURI()
	inst_list, serv_err := service2instlist(uri_head(uri_input)) //Choose service
	if serv_err != nil {
		return url_in, serv_err
	}

	url_out, inst_err := choose_instance(inst_list, load_balancer_info) //Choose instance
	if inst_err != nil {
		return url_out, inst_err
	}
	url_out.Path = uri_tail(uri_input)
	return url_out, nil
}

func service2instlist(requested_service string) ([]string, error) {
	var inst_list []string
	available_service_flag := false
	for _, cur_service := range HealthChecker_ptr.Service { //find from healthchecker
		if requested_service == cur_service.Name {
			inst_list = cur_service.Available_instance
			available_service_flag = true
			break
		}
	}
	if available_service_flag == false {
		return nil,
			errors.New(RESOURCE_NONEXISTENT_ERROR + ERROR_STRING_SEPARATOR +
				"No service exists for particular uri")
	} else if len(inst_list) == 0 {
		return nil,
			errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
				"No instance available for particular uri")
	}
	return inst_list, nil
}
func choose_instance(instance_list []string, load_balancer_info []string) (*url.URL, error) {
	instance_chosen, err := load_balance(instance_list, Config_ptr.Load_balancer_policy, load_balancer_info)
	if err != nil {
		return nil, err
	}
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
