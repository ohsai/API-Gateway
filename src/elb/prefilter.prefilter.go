package main

import (
	//"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Region struct {
	Region_code    string
	Available_Zone []string
}

var Region_ptr *Region

func yamlDecoderR(Ryamlpath string) (*Region, error) {
	r_out := new(Region)
	yaml_file, yaml_open_err := ioutil.ReadFile(Ryamlpath)
	if yaml_open_err != nil {
		log.Println("error while opening yaml", yaml_open_err.Error())
		return r_out, yaml_open_err
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &r_out)
	if yaml_decode_err != nil {
		log.Println("error while unmarshal on yaml", yaml_decode_err.Error())
		return r_out, yaml_decode_err
	}
	R_print(r_out)
	return r_out, nil
}
func R_print(R_in *Region) {
	log.Println("ELB$Region Structure :")
	fmt.Println("  ", R_in.Region_code)
	for _, cur_AZ := range R_in.Available_Zone {
		fmt.Println("    ", cur_AZ)
	}
	fmt.Println("ELB$cur region code : ", R_in.Region_code)
}
func prefilter(req *http.Request) (*http.Request, error) {
	load_balancer_info := []string{strings.Split(req.RemoteAddr, ":")[0], req.URL.Path}
	url_out, err_url := format_url(load_balancer_info)
	if err_url != nil {
		return req, err_url
	}
	//modify url to that region proxy
	proxyReq, err := format_request_plain(url_out, req)
	if err != nil {
		return req, err
	}
	return proxyReq, err
}
func format_url(load_balancer_info []string) (*url.URL, error) {
	//Check region structure, Find closest region
	AZ_list := Region_ptr.Available_Zone
	//Choose AZ
	url_out, inst_err := choose_instance(AZ_list, load_balancer_info)
	if inst_err != nil {
		return url_out, inst_err
	}
	url_out.Path = load_balancer_info[1]
	return url_out, nil
}

func choose_instance(instance_list []string, load_balancer_info []string) (*url.URL, error) {
	instance_chosen, err := load_balance(instance_list, "round_robin", load_balancer_info)
	if err != nil {
		return nil, err
	}
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}
