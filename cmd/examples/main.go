package main

import (
	"fmt"
	"log"

	"github.com/bishopfox/awsservicemap"
)

func main() {

	// Create a variable of type AwsServiceMap
	// JsonFileSource options: "EMBEDDED_IN_PACKAGE", "DOWNLOAD_FROM_AWS"
	// When using "EMBEDDED_IN_PACKAGE" this package does not make any external HTTP requests, but the data might be slightly out of date
	// When using "DOWNLOAD_FROM_AWS" this package makes an external HTTP request, but you get real-time data.
	// With the new caching feature, the data is fetched only once per instance regardless of how many method calls are made.

	servicemap := &awsservicemap.AwsServiceMap{
		//JsonFileSource: "EMBEDDED_IN_PACKAGE",
		JsonFileSource: "DOWNLOAD_FROM_AWS",
	}

	//  Example of how you can also use the constructor pattern to simulate "instantiating" a new service map "object"
	// servicemap1 := awsservicemap.NewServiceMap()
	// servicemap1.JsonFileSource = "EMBEDDED_IN_PACKAGE"

	// Check what regions support grafana?
	regions, err := servicemap.GetRegionsForService("grafana")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(regions)
	// Check what services are supported in eu-south-2
	services, err := servicemap.GetServicesForRegion("eu-south-2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(services)

	// List all of the regions
	totalRegions, err := servicemap.GetAllRegions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(totalRegions)

	// List all of the services
	totalServices, err := servicemap.GetAllServices()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(totalServices)

	// Check if franddetector is supported in eu-south-2?
	res, err := servicemap.IsServiceInRegion("ec2", "eu-west-1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
