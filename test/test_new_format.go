package main

import (
	"fmt"
	"log"

	"github.com/bishopfox/awsservicemap"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	servicemap := &awsservicemap.AwsServiceMap{
		JsonFileSource: "DOWNLOAD_FROM_AWS",
	}

	fmt.Println("=== Testing awsservicemap with new AWS API format ===\n")

	// Test IsServiceInRegion
	fmt.Println("1. Testing IsServiceInRegion for EC2:")
	testRegions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}
	for _, region := range testRegions {
		res, err := servicemap.IsServiceInRegion("ec2", region)
		if err != nil {
			fmt.Printf("   ERROR checking %s: %v\n", region, err)
		} else {
			fmt.Printf("   EC2 in %s: %v\n", region, res)
		}
	}

	// Test GetRegionsForService
	fmt.Println("\n2. Testing GetRegionsForService for EC2:")
	regions, err := servicemap.GetRegionsForService("ec2")
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
	} else {
		fmt.Printf("   Found EC2 in %d regions\n", len(regions))
		fmt.Printf("   First 5 regions: %v\n", regions[:min(5, len(regions))])
	}

	// Test GetAllServices
	fmt.Println("\n3. Testing GetAllServices:")
	services, err := servicemap.GetAllServices()
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
	} else {
		fmt.Printf("   Found %d total services\n", len(services))
		fmt.Printf("   First 10 services: %v\n", services[:min(10, len(services))])
	}

	// Test GetServicesForRegion
	fmt.Println("\n4. Testing GetServicesForRegion for us-east-1:")
	usEast1Services, err := servicemap.GetServicesForRegion("us-east-1")
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
	} else {
		fmt.Printf("   Found %d services in us-east-1\n", len(usEast1Services))
		fmt.Printf("   First 10 services: %v\n", usEast1Services[:min(10, len(usEast1Services))])
	}

	// Test GetAllRegions
	fmt.Println("\n5. Testing GetAllRegions:")
	allRegions, err := servicemap.GetAllRegions()
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
	} else {
		fmt.Printf("   Found %d total regions\n", len(allRegions))
		fmt.Printf("   All regions: %v\n", allRegions)
	}

	// Test a few more services
	fmt.Println("\n6. Testing other services:")
	otherServices := []string{"rds", "lambda", "s3", "iam", "eks"}
	for _, svc := range otherServices {
		res, err := servicemap.IsServiceInRegion(svc, "us-east-1")
		if err != nil {
			fmt.Printf("   ERROR checking %s: %v\n", svc, err)
		} else {
			fmt.Printf("   %s in us-east-1: %v\n", svc, res)
		}
	}

	fmt.Println("\n=== All tests completed successfully! ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
