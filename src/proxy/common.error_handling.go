package main

import (
	//"log"
	"net/http"
	"strings"
)

func filter_error_handler(filter_name string, w http.ResponseWriter, err error) {
	//log.Println(filter_name+" terminated proxy with error : ", err.Error())
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
