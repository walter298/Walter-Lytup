package main

func main() {
	serviceMap := ServiceMap{
		"test": TestService{},
		"weather": MakeWeatherService(),
	}
	allowedResponseCodes := []int{500, 400, 406, 200}
	server := MakeServer(Address_1, serviceMap, allowedResponseCodes)

	serverManager := ServerManager{}
	serverManager.AddAsyncServer(server)
	serverManager.Run()
}