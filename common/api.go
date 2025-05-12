package common

type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

