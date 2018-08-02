package main

import (
	"log"
	"net/http"
)

func routing_filter(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	proxyRes, err := client.Do(req)
	if err != nil {
		log.Println("Unable to create forward response! : ", err.Error())
	}
	//Pass headers crucial for proxy
	proxyRes.Header.Add("Service", req.Header.Get("Service"))
	return proxyRes, err

}
