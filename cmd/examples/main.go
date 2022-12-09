package main

import (
	"fmt"
	"log"

	"github.com/bishopfox/awsservicemap/pkg/awsservicemap"
)

func main() {

	// Instantiate a new servicemap object
	// JsonFileSource options: "EMBEDDED_IN_PACKAGE", "DOWNLOAD_FROM_AWS"
	servicemap := &awsservicemap.AwsServiceMap{
		JsonFileSource: "EMBEDDED_IN_PACKAGE",
	}

	//  Example of how to Instantiate a new servicemap object using constructor
	// JsonFileSource options: "EMBEDDED_IN_PACKAGE", "DOWNLOAD_FROM_AWS"
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
	res, err := servicemap.IsServiceInRegion("frauddetector", "eu-south-2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
