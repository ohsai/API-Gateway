package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"proxy/mycrypt"
)

func request_authentication(req *http.Request) (bool, error) {
	authtoken := &Signin_Resp_to_Client{}
	authtoken_str := req.Header.Get("AuthToken")
	if authtoken_str == "" {
		return false,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken header does not exist")
	}
	err := json.Unmarshal([]byte(authtoken_str), authtoken)
	if err != nil {
		return false,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken header not in form of authentication token")
	}
	validity := mycrypt.CheckMAC((authtoken.Username + authtoken.Role), authtoken.Hash, Config_ptr.Auth_key)
	if validity {
		return validity, nil
	} else {
		return validity,
			errors.New(AUTHENTICATION_TOKEN_ERROR + ERROR_STRING_SEPARATOR +
				"AuthToken failed authentication")
	}
}
