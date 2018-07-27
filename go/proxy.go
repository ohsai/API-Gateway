package main

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"net/url"
)

func main() {
	if len(os.Args) != 2 {
		panic("Error : need one and only one argument specifying port number")
	}
	createProxy(os.Args[1:2][0])
}
func createProxy(PORT string){
	http.Handle("/",new(proxyHandler))
	http.ListenAndServe(":"+PORT, nil)
}
type proxyHandler struct{
	http.Handler

}
func reform_url (url_in *url.URL) *url.URL {
	url_out := url_in 
	url_out.Host = "localhost:5000"
	if !url_out.IsAbs() {
		url_out.Scheme = "http"
	}
	return url_out
}
func (h* proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request){
	proxyReq, errp := prefilter (req)
    if errp != nil {
	    fmt.Println("prefilter terminated proxy with error : ", errp.Error())
	    return 
    }
    proxyRes, errr := routing_filter(proxyReq)
    if errr != nil {
	    fmt.Println("routing_filter terminated proxy with error : ", errr.Error())
	    return 
    }
    post_filter (proxyRes, w)
}
func prefilter(req * http.Request) (*http.Request, error){
	newurl := reform_url(req.URL)
    proxyReq, err := http.NewRequest(req.Method, newurl.String(), req.Body)
    if err != nil {
	fmt.Println("Unable to create forward request!",err.Error())
	return proxyReq, err 
    }

    proxyReq.Header.Set("Host", req.Host)
    proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

    for header, values := range req.Header {
        for _, value := range values {
            proxyReq.Header.Add(header, value)
        }
    }
    return proxyReq,err 
}
func routing_filter(req * http.Request) (*http.Response, error){
    client := &http.Client{}
    proxyRes, err := client.Do(req)
    if err != nil {
	fmt.Println("Unable to create forward response!",err.Error())
    }
    return proxyRes,err
    
}
func post_filter(proxyRes *http.Response, w http.ResponseWriter){
	for header, values := range proxyRes.Header{
		for _, value := range values{
			w.Header().Add(header, value)
		}
	}
    io.Copy(w, proxyRes.Body)
    proxyRes.Body.Close()
}


