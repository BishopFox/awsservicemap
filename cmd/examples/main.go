package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap/pkg/awsservicemap"
)

func main() {
	// Instantiate a new servicemap object
	servicemap := awsservicemap.NewServiceMap()

	// Check what regions support grafana?
	regions := servicemap.GetRegionsForService("grafana")
	fmt.Println(regions)
	// Check what services are supported in eu-south-2
	services := servicemap.GetServicesForRegion("eu-south-2")
	fmt.Println(services)

	// List all of the regions
	totalRegions := servicemap.GetAllRegions()
	fmt.Println(totalRegions)

	// List all of the services
	totalServices := servicemap.GetAllServices()
	fmt.Println(totalServices)

	// Check if franddetector is supported in eu-south-2?
	res := servicemap.IsServiceInRegion("frauddetector", "eu-south-2")
	fmt.Println(res)
}
