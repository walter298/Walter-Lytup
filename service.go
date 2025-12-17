package main

import "encoding/json"

type RequestJson struct {
	HandlerID string          `json:"ServiceID"`
	JsonBody  json.RawMessage `json:"Body"`
}

type ServiceResult struct {
	Code  int
	Error string
	Yield any
}

type Service interface {
	Run(RequestJson) ServiceResult
}

type ServiceMap = map[string]Service
