package main

import (
	//"fmt"
	//"github.com/jinzhu/copier"
	"log"
	"net/http"
	"os"
	"strings"
)

var MSA_ptr *MSA
var auth_key string
var HealthChecker_ptr *MSA

func main() {
	if len(os.Args) != 3 {
		panic("Usage : (program name) <port> <yaml path for MSA structure>")
	}
	proxy_init(os.Args[1:2][0], os.Args[2:3][0])
}
func proxy_init(PORT string, MSAyamlpath string) {
	MSA_ptr = yamlDecoder(MSAyamlpath)
	auth_key = "private_key"
	MSA_print(MSA_ptr)
	healthchecker_init()
	createListener(PORT)
}
func createListener(PORT string) {
	http.Handle("/", new(proxyHandler))
	http.ListenAndServe(":"+PORT, nil)
}

type proxyHandler struct {
	http.Handler
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	proxyReq, prefilter_err := prefilter(req)
	if prefilter_err != nil {
		filter_error_handler("pre_filter", w, prefilter_err)
		return
	}
	proxyRes, routing_err := routing_filter(proxyReq)
	if routing_err != nil {
		filter_error_handler("routing_filter", w, routing_err)
		return
	}
	post_err := post_filter(proxyRes, w)
	if post_err != nil {
		filter_error_handler("post_filter", w, post_err)
		return
	}
}
func filter_error_handler(filter_name string, w http.ResponseWriter, err error) {
	log.Println(filter_name+" terminated proxy with error : ", err.Error())
	error_type := strings.Split(err.Error(), ERROR_STRING_SEPARATOR)[0]
	if error_type == AUTHENTICATION_TOKEN_ERROR {
		w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	} else if error_type == RESOURCE_NONEXISTENT_ERROR {
		w.WriteHeader(http.StatusNotFound)

	} else if error_type == NO_AVAILABLE_INSTANCE_ERROR {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
	return
}

var AUTHENTICATION_TOKEN_ERROR string = "AuthTokenError"
var ERROR_STRING_SEPARATOR string = "$"
var RESOURCE_NONEXISTENT_ERROR string = "NotFoundError"
var NO_AVAILABLE_INSTANCE_ERROR string = "NoInstanceError"
