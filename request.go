package main

import "net/http"

type Request struct {
	SrID string         `json:"Srvc"`
	Seed map[string]any `json:"Seed"`
}
type ResponseHandler struct {
	Code    string
	Program func(*http.Request, string, map[string]any) (int, string, any)
}

var DHI0_SPRegister []*ResponseHandler = []*ResponseHandler{}
var DNI0_AllowedResponseCode []int = []int{500, 400, 406, 200}
var DNI0_ResponseHeaders [][]string = [][]string{[]string{"Content-Type", "application/json"}}
