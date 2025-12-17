package main

type TestService struct {}

func (TestService) Run(RequestJson) ServiceResult {
	return ServiceResult{200, "", "hello world"}
}