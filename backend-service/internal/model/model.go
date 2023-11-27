package model

type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type GeneralResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
