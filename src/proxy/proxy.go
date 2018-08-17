package main

import (
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 5 {
		panic("Usage : (program name) <port> <yaml path MSA structure> 	<yaml path Region structure> <yaml path Configuration> 	")
	}
	proxy_init(os.Args[1:2][0], os.Args[2:3][0], os.Args[3:4][0], os.Args[4:5][0])
}
func proxy_init(PORT string, MSAyamlpath string, RSyamlpath string, Configjsonpath string) {
	config_manager_init(MSAyamlpath, RSyamlpath, Configjsonpath)
	healthchecker_init()
	redis_init(PORT)
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
