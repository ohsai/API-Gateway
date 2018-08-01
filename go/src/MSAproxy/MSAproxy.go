package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Microservice struct {
	Name     string
	Instance []string
}

type MSA struct {
	Name    string
	Service []Microservice
}

var MSA_ptr *MSA

func main() {
	if len(os.Args) != 3 {
		panic("Usage : (program name) <port> <yaml path for MSA structure>")
	}
	MSA_ptr = yamlDecoder(os.Args[2:3][0])
	fmt.Println(MSA_ptr.Name, MSA_ptr.Service[0].Name, MSA_ptr.Service[0].Instance[0])
	createProxy(os.Args[1:2][0])
}
func createProxy(PORT string) {
	http.Handle("/", new(proxyHandler))
	http.ListenAndServe(":"+PORT, nil)
}

type proxyHandler struct {
	http.Handler
}

func yamlDecoder(MSAyamlpath string) *MSA {
	msa_out := new(MSA)
	yaml_file, yaml_open_err := ioutil.ReadFile(MSAyamlpath)
	if yaml_open_err != nil {
		fmt.Println("error while opening yaml", yaml_open_err.Error())
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &msa_out)
	if yaml_decode_err != nil {
		fmt.Println("error while unmarshal on yaml", yaml_decode_err.Error())
	}
	fmt.Printf("---MSA : \n%+v\n", msa_out)
	return msa_out
}
func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	proxyReq, prefilter_err := prefilter(MSA_ptr, req)
	if prefilter_err != nil {
		fmt.Println("prefilter terminated proxy with error : ", prefilter_err.Error())
		w.WriteHeader(500)
		w.Write([]byte(http.StatusText(500)))
		return
	}
	proxyRes, routing_err := routing_filter(proxyReq)
	if routing_err != nil {
		fmt.Println("routing_filter terminated proxy with error : ", routing_err.Error())
		w.WriteHeader(500)
		w.Write([]byte(http.StatusText(500)))
		return
	}
	post_filter(proxyRes, w)
}
func prefilter(structure *MSA, req *http.Request) (*http.Request, error) {
	newurl, urlerr := reform_url(MSA_ptr, req.URL)
	if urlerr != nil {
		fmt.Println("Unable to route by url : ", urlerr.Error())
		return req, urlerr
	}
	proxyReq, err := http.NewRequest(req.Method, newurl.String(), req.Body)
	if err != nil {
		fmt.Println("Unable to create forward request! : ", err.Error())
		return proxyReq, err
	}

	//proxyReq.Header.Set("Host", req.Host)
	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	return proxyReq, err
}
func choose_service(structure *MSA, uri_input string) ([]string, error) {
	parts := strings.Split(uri_input, string(os.PathSeparator))
	fmt.Println("Requested Service : ", parts[1])
	var inst_list []string
	requested_service := parts[1]
	available_service_flag := false
	for _, cur_service := range structure.Service {
		if requested_service == cur_service.Name {
			inst_list = cur_service.Instance
			available_service_flag = true
			break
		}
	}
	if available_service_flag == false {
		return nil, errors.New("No service exists for particular uri")
	}
	return inst_list, nil
}
func choose_instance(instance_list []string) (*url.URL, error) {
	instance_chosen := instance_list[rand.Intn(len(instance_list))]
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}

func reform_url(structure *MSA, url_in *url.URL) (*url.URL, error) {
	uri_input := url_in.RequestURI()
	inst_list, serv_err := choose_service(structure, uri_input)
	if serv_err != nil {
		fmt.Println("Error while choosing service url :", serv_err.Error())
		return url_in, serv_err

	}
	url_out, inst_err := choose_instance(inst_list)
	if inst_err != nil {
		fmt.Println("Error while choosing instance url :", inst_err.Error())
		return url_out, inst_err
	}
	parts := strings.Split(uri_input, string(os.PathSeparator))
	rest_of_uri := strings.Join(parts[2:], string(os.PathSeparator))
	fmt.Println("Rest part of URI : ", rest_of_uri)
	fmt.Println("URI Before : ", url_out.String())
	url_out.Path = rest_of_uri
	fmt.Println("URI After : ", url_out.String())
	return url_out, nil
}
func routing_filter(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	proxyRes, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to create forward response! : ", err.Error())
	}
	return proxyRes, err

}
func post_filter(proxyRes *http.Response, w http.ResponseWriter) {
	for header, values := range proxyRes.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	io.Copy(w, proxyRes.Body)
	proxyRes.Body.Close()
}
