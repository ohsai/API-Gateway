package main

import (
	"errors"
	//"log"
	"math/rand"
	"net/http"
	"net/url"
	"proxy/mycrypt"
	"strings"
)

func prefilter(req *http.Request) (*http.Request, error) {
	//format URL
	load_balancer_info := []string{strings.Split(req.RemoteAddr, ":")[0], req.URL.Path}
	newurl, err := format_url(req.URL, load_balancer_info)
	if err != nil {
		return req, err
	}
	//format Request and test Authenticity
	var proxyReq *http.Request
	var valid bool
	if url_to_service(req.URL) == "auth" {
		proxyReq, err = format_request_auth(newurl, req)
	} else {
		//auth layer
		if valid, err = request_authentication(req); valid {
			//request formatting layer
			proxyReq, err = format_request_regular(newurl, req)
		}
	}
	if err != nil {
		return req, err
	}

	return proxyReq, err
}
func format_url(url_in *url.URL, load_balancer_info []string) (*url.URL, error) {
	uri_input := url_in.RequestURI()
	inst_list, serv_err := service2instlist(uri_input)
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
func service2instlist(uri_input string) ([]string, error) {
	var inst_list []string
	requested_service := uri_head(uri_input)
	available_service_flag := false
	for _, cur_service := range HealthChecker_ptr.Service {
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
	} else if len(inst_list) == 0 {
		return nil,
			errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
				"No instance available for particular uri")
	}
	//log.Println(inst_list)
	return inst_list, nil
}
func choose_instance(instance_list []string, load_balancer_info []string) (*url.URL, error) {
	instance_chosen, err := load_balance(instance_list, "random", load_balancer_info)
	if err != nil {
		return nil, err
	}
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
func load_balance(instance_list []string, policy string, load_balancer_info []string) (string, error) {
	var instance_chosen string = ""
	var chosen_index int
	for i := 0; i < len(instance_list); i++ {
		//Implement round robin / weighted_round_robin / random
		if policy == "round_robin" {

		} else if policy == "weighted_round_robin" {
			//put weight in msa.yaml

		} else if policy == "ip_hash" {
			chosen_index = mycrypt.String_modhash(load_balancer_info[0], len(instance_list))

			//need url
		} else if policy == "url_hash" {
			chosen_index = mycrypt.String_modhash(
				load_balancer_info[0]+load_balancer_info[1], len(instance_list))
			//need url
		} else { //random
			chosen_index = rand.Intn(len(instance_list))
		}
		instance_chosen = instance_list[chosen_index]
		//Implement one_more_check, retry
		if err := ping(instance_chosen); err == nil {
			return instance_chosen, nil
		}
	}
	return "", errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
		"instance failed before HealthChecker check")
}
