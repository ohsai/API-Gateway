package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func prefilter(req *http.Request) (*http.Request, error) {
	// Region check

	//Check region code (Header)
	req_region_code := req.Header.Get("Accept-Language")
	load_balancer_info := []string{strings.Split(req.RemoteAddr, ":")[0], req.URL.Path}
	//Check current region
	if req_region_code == "" || req_region_code == RS_ptr.Cur_region_code {
		// If region is right or no region code
		//format URL
		newurl, err := format_url(req.URL, load_balancer_info)
		if err != nil {
			return req, err
		}
		//format Request and test Authenticity
		var proxyReq *http.Request
		var valid bool
		if url_to_service(req.URL) == "auth" {
			proxyReq, err = format_request_plain(newurl, req)
		} else {
			//auth layer
			if valid, err = request_authentication(req); valid {
				//request formatting layer
				proxyReq, err = format_request_after_auth(newurl, req)
			}
		}
		if err != nil {
			return req, err

		}
		return proxyReq, err

	} else {
		/*
			Region Proxying
			Else if region is wrong
			Format url to pass request to right regional proxy
		*/
		//Check region structure, Find closest region
		AZ_list, region_err := region2AZlist(req_region_code)
		if region_err != nil {
			return req, region_err
		}
		//Choose one of region elbs
		url_out, inst_err := choose_instance(AZ_list, load_balancer_info)
		if inst_err != nil {
			return req, inst_err
		}
		url_out.Path = req.URL.Path
		//modify url to that region elb
		proxyReq, err := format_request_plain(url_out, req)
		if err != nil {
			return req, err
		}
		return proxyReq, err
	}
}
func format_url(url_in *url.URL, load_balancer_info []string) (*url.URL, error) {
	uri_input := url_in.RequestURI()
	inst_list, serv_err := service2instlist(uri_head(uri_input))
	if serv_err != nil {
		//log.Println("Error while choosing service url :", serv_err.Error())
		return url_in, serv_err
	}

	url_out, inst_err := choose_instance(inst_list, load_balancer_info)
	if inst_err != nil {
		//log.Println("Error while choosing instance url :", inst_err.Error())
		return url_out, inst_err
	}
	url_out.Path = uri_tail(uri_input)
	return url_out, nil
}

//func service2instlist(uri_input string) ([]string, error) {
func service2instlist(requested_service string) ([]string, error) {
	var inst_list []string
	//requested_service := uri_head(uri_input)
	available_service_flag := false
	for _, cur_service := range HealthChecker_ptr.Service {
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
	//log.Println(inst_list)
	return inst_list, nil
}
func choose_instance(instance_list []string, load_balancer_info []string) (*url.URL, error) {
	//instance_chosen, err := load_balance(instance_list, "round_robin", load_balancer_info)
	instance_chosen, err := load_balance(instance_list, Config_ptr.Load_balancer_policy, load_balancer_info)
	if err != nil {
		return nil, err
	}
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
