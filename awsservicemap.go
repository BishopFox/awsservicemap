package awsservicemap

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed data/aws-service-regions.json
var awsJson embed.FS

type serviceEntry struct {
	ID string `json:"id"`
}

type regionalServiceData struct {
	ServiceEntries []serviceEntry `json:"prices"`
}

// Checks if element is part of array.
func contains(element string, array []string) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func parseJson() regionalServiceData {

	jsonFile, err := awsJson.ReadFile("data/aws-service-regions.json")
	if err != nil {
		fmt.Println(err)
	}

	var serviceData regionalServiceData
	json.Unmarshal([]byte(jsonFile), &serviceData)
	return serviceData
}

func GetAllRegions() []string {
	totalRegions := []string{}
	serviceData := parseJson()
	for _, id := range serviceData.ServiceEntries {
		region := strings.Split(id.ID, ":")[1]
		if !contains(region, totalRegions) {
			totalRegions = append(totalRegions, region)
		}
	}
	return totalRegions

}

func GetRegionsForService(reqService string) []string {
	regionsForServiceMap := map[string][]string{}

	serviceData := parseJson()
	for _, id := range serviceData.ServiceEntries {
		service := strings.Split(id.ID, ":")[0]
		if _, ok := regionsForServiceMap[service]; !ok {
			regionsForServiceMap[service] = nil
		}
		region := strings.Split(id.ID, ":")[1]
		if _, ok := regionsForServiceMap[service]; ok {
			regionsForServiceMap[service] = append(regionsForServiceMap[service], region)
		}
	}
	return regionsForServiceMap[reqService]

}

func GetAllServices() []string {
	totalServices := []string{}

	serviceData := parseJson()
	for _, id := range serviceData.ServiceEntries {
		service := strings.Split(id.ID, ":")[0]
		if !contains(service, totalServices) {
			totalServices = append(totalServices, service)
		}
	}
	return totalServices

}

func GetServicesForRegion(reqRegion string) []string {
	servicesForRegionMap := map[string][]string{}
	serviceData := parseJson()

	for _, id := range serviceData.ServiceEntries {

		region := strings.Split(id.ID, ":")[1]
		if _, ok := servicesForRegionMap[region]; !ok {
			servicesForRegionMap[region] = nil
		}

		service := strings.Split(id.ID, ":")[0]
		if _, ok := servicesForRegionMap[region]; ok {
			servicesForRegionMap[region] = append(servicesForRegionMap[region], service)
		}
	}
	return servicesForRegionMap[reqRegion]
}
