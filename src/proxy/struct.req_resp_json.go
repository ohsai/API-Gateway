package main

type AuthSigninReq struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}
type AuthSignupReq struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
	Role     string `json:"role", db:"role"`
}

type Signin_Resp_from_Auth struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
type Signin_Resp_to_Client struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Hash     []byte `json:"hash"`
}
