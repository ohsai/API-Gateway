package main

import (
	//"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic("Usage : (program name) <port> <yaml path for MSA structure>")
	}
	proxy_init(os.Args[1:2][0], os.Args[2:3][0])
}
func proxy_init(PORT string, Ryamlpath string) {
	//MSA_ptr = yamlDecoder(MSAyamlpath)
	//auth_key = "private_key"
	Region_ptr = yamlDecoderR(Ryamlpath)
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
