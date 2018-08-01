package main

import (
	"./mycrypt"
	"encoding/json"
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
var auth_key string

func main() {
	if len(os.Args) != 3 {
		panic("Usage : (program name) <port> <yaml path for MSA structure>")
	}
	proxy_init(os.Args[1:2][0], os.Args[2:3][0])
}
func proxy_init(PORT string, MSAyamlpath string) {
	MSA_ptr = yamlDecoder(MSAyamlpath)
	auth_key = "private_key"
	createListener(PORT)
	fmt.Println(MSA_ptr.Name, MSA_ptr.Service[0].Name, MSA_ptr.Service[0].Instance[0])
}
func createListener(PORT string) {
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
func prefilter(structure *MSA, req *http.Request) (*http.Request, string, error) {
	req_serv, newurl, urlerr := reform_url(MSA_ptr, req.URL)
	if urlerr != nil {
		fmt.Println("Unable to route by url : ", urlerr.Error())
		return req, req_serv, urlerr
	}
	var proxyReq *http.Request
	var err error
	if req_serv == "auth" {
		proxyReq, err = reform_request_auth(newurl, req)
	} else {
		proxyReq, err = reform_request(newurl, req)
	}
	if err != nil {
		fmt.Println("Unable to create request from url : ", err.Error())
		return proxyReq, req_serv, err
	}
	return proxyReq, req_serv, err
}
func reform_request(newurl *url.URL, req *http.Request) (*http.Request, error) {
	//Authentication
	authtoken := &AuthRespwToken{}
	authtoken_str := req.Header.Get("AuthToken")
	err := json.Unmarshal([]byte(authtoken_str), authtoken)
	validity := mycrypt.CheckMAC((authtoken.Username + authtoken.Role), authtoken.Hash, auth_key)
	if validity == false {
		return req, errors.New("Invalid authentication token")
	}

	proxyReq, err := http.NewRequest(req.Method, newurl.String(), req.Body)
	if err != nil {
		fmt.Println("Unable to create forward request! : ", err.Error())
		return proxyReq, err
	}

	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	return proxyReq, err
}
func reform_request_auth(newurl *url.URL, req *http.Request) (*http.Request, error) {
	proxyReq, err := http.NewRequest(req.Method, newurl.String(), req.Body)
	if err != nil {
		fmt.Println("Unable to create forward request! : ", err.Error())
		return proxyReq, err
	}

	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	return proxyReq, err
}

func service2instlist(structure *MSA, uri_input string) (string, []string, error) {
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
		return "", nil, errors.New("No service exists for particular uri")
	}
	return requested_service, inst_list, nil
}
func choose_instance(instance_list []string) (*url.URL, error) {
	instance_chosen := instance_list[rand.Intn(len(instance_list))]
	url_out, err := url.ParseRequestURI(instance_chosen)
	return url_out, err
}

func reform_url(structure *MSA, url_in *url.URL) (string, *url.URL, error) {
	uri_input := url_in.RequestURI()
	req_serv, inst_list, serv_err := service2instlist(structure, uri_input)
	if serv_err != nil {
		fmt.Println("Error while choosing service url :", serv_err.Error())
		return req_serv, url_in, serv_err

	}

	url_out, inst_err := choose_instance(inst_list)
	if inst_err != nil {
		fmt.Println("Error while choosing instance url :", inst_err.Error())
		return req_serv, url_out, inst_err
	}
	parts := strings.Split(uri_input, string(os.PathSeparator))
	rest_of_uri := strings.Join(parts[2:], string(os.PathSeparator))
	url_out.Path = rest_of_uri
	return req_serv, url_out, nil
}
func routing_filter(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	proxyRes, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to create forward response! : ", err.Error())
	}
	return proxyRes, err

}

type AuthResp struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
type AuthRespwToken struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Hash     string `json:"hash"`
}

func post_filter(req_serv string, proxyRes *http.Response, w http.ResponseWriter) error {
	fmt.Println("requested service : ", req_serv)

	if req_serv == "auth" {
		err := reform_response_auth(proxyRes, w)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		err := reform_response(proxyRes, w)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}
func reform_response_auth(proxyRes *http.Response, w http.ResponseWriter) error {
	temp := &AuthResp{}
	jsonparseerr := json.NewDecoder(proxyRes.Body).Decode(temp)
	if jsonparseerr != nil {
		fmt.Println(jsonparseerr.Error())
		return jsonparseerr
	}
	authresp := AuthRespwToken{
		Username: temp.Username,
		Role:     temp.Role,
		Hash:     string(mycrypt.CreateMAC(temp.Username+temp.Role, auth_key)),
	}
	b, err := json.Marshal(authresp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(authresp)
	fmt.Println(string(b))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func reform_response(proxyRes *http.Response, w http.ResponseWriter) error {
	for header, values := range proxyRes.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	if _, err := io.Copy(w, proxyRes.Body); err != nil {
		fmt.Println(err.Error())
		return err
	}
	proxyRes.Body.Close()
	return nil
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	proxyReq, req_serv, prefilter_err := prefilter(MSA_ptr, req)
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
	post_filter(req_serv, proxyRes, w)
}
