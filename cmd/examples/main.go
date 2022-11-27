package main

import (
	"fmt"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	regions := awsservicemap.GetRegionsForService("grafana")
	fmt.Println(regions)
	services := awsservicemap.GetServicesForRegion("eu-south-2")
	fmt.Println(services)
	totalRegions := awsservicemap.GetAllRegions()
	fmt.Println(totalRegions)
	totalServices := awsservicemap.GetAllServices()
	fmt.Println(totalServices)
}
