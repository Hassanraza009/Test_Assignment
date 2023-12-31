package main

import (
	"fmt"

	"test/rest"
	"test/service"
)

func main() {

	fmt.Println("#==================================#")
	fmt.Println("#===========Starting Server =======#")
	fmt.Println("#==================================#")

	/*
	* Initiate Service Layer Container
	 */

	serviceContainer := service.NewServiceContainer()

	/*
	* Initiate Rest Server
	 */
	rest.StartServer(serviceContainer)

	fmt.Println("========== Rest Server Started ============")
	fmt.Println("========== Server Running ============")

	select {}

}
