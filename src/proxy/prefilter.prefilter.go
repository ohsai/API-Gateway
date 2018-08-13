package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Region_Structure struct {
	Cur_region_code string
	Regions         []Region
}
type Region struct {
	Region_code    string
	Available_Zone []string
}

var RS_ptr *Region_Structure

func yamlDecoderRS(RSyamlpath string) *Region_Structure {
	rs_out := new(Region_Structure)
	yaml_file, yaml_open_err := ioutil.ReadFile(RSyamlpath)
	if yaml_open_err != nil {
		log.Println("error while opening yaml", yaml_open_err.Error())
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &rs_out)
	if yaml_decode_err != nil {
		log.Println("error while unmarshal on yaml", yaml_decode_err.Error())
	}
	RS_print(rs_out)
	return rs_out
}
func RS_print(RS_in *Region_Structure) {
	log.Println("Region Structure :")
	for _, cur_region := range RS_in.Regions {
		fmt.Println("  ", cur_region.Region_code)
		for _, cur_AZ := range cur_region.Available_Zone {
			fmt.Println("    ", cur_AZ)
		}
	}
	fmt.Println("cur region code : ", RS_in.Cur_region_code)
}
func region2AZlist(region_code_in string) ([]string, error) {
	var AZ_list []string
	//requested_service := uri_head(uri_input)
	available_zone_flag := false
	for _, cur_region := range RS_ptr.Regions {
		if region_code_in == cur_region.Region_code {
			AZ_list = cur_region.Available_Zone
			available_zone_flag = true
			break
		}
	}
	if available_zone_flag == false {
		return nil,
			errors.New(RESOURCE_NONEXISTENT_ERROR + ERROR_STRING_SEPARATOR +
				"No region proxy exists for particular locale")
	} else if len(AZ_list) == 0 {
		return nil,
			errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
				"No Available Zone instance for particular locale")
	}
	//log.Println(inst_list)
	return AZ_list, nil
}
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
		//Else if region is wrong
		//Format url to pass request to right regional proxy
		//
		//Check region structure, Find closest region
		AZ_list, region_err := region2AZlist(req_region_code)
		if region_err != nil {
			return req, region_err
		}
		//Choose AZ
		url_out, inst_err := choose_instance(AZ_list, load_balancer_info)
		if inst_err != nil {
			return req, inst_err
		}
		url_out.Path = req.URL.Path
		//	log.Println(url_out)
		//modify url to that region proxy
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
	instance_chosen, err := load_balance(instance_list, "round_robin", load_balancer_info)
	if err != nil {
		return nil, err
	}
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
