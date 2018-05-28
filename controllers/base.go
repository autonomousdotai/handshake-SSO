package controllers

type JsonResponse struct {
    Status int `json:"status"`
    Message string `json:"message,omitempty"`
    Data interface{} `json:"data,omitempty"`
}
