package awsservicemap

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//go:embed data/aws-service-regions.json
var awsJson embed.FS

type AwsServiceMap struct {
	DownloadJson bool
}

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

func NewServiceMap() *AwsServiceMap {
	return &AwsServiceMap{
		DownloadJson: false,
	}
}

func parseJson(downloadJson bool) regionalServiceData {
	var serviceData regionalServiceData

	if downloadJson {
		res, err := http.Get("https://api.regional-table.region-services.aws.a2z.com/index.json")
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal([]byte(body), &serviceData)

	} else {
		jsonFile, err := awsJson.ReadFile("data/aws-service-regions.json")
		json.Unmarshal([]byte(jsonFile), &serviceData)
		if err != nil {
			fmt.Println(err)
		}
	}

	return serviceData
}

func (*AwsServiceMap) GetAllRegions() []string {
	totalRegions := []string{}
	servicemap := NewServiceMap()
	serviceData := parseJson(servicemap.DownloadJson)
	for _, id := range serviceData.ServiceEntries {
		region := strings.Split(id.ID, ":")[1]
		if !contains(region, totalRegions) {
			totalRegions = append(totalRegions, region)
		}
	}
	return totalRegions

}

func (*AwsServiceMap) GetRegionsForService(reqService string) []string {
	regionsForServiceMap := map[string][]string{}

	servicemap := NewServiceMap()
	serviceData := parseJson(servicemap.DownloadJson)
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

func (*AwsServiceMap) GetAllServices() []string {
	totalServices := []string{}

	servicemap := NewServiceMap()
	serviceData := parseJson(servicemap.DownloadJson)
	for _, id := range serviceData.ServiceEntries {
		service := strings.Split(id.ID, ":")[0]
		if !contains(service, totalServices) {
			totalServices = append(totalServices, service)
		}
	}
	return totalServices

}

func (*AwsServiceMap) GetServicesForRegion(reqRegion string) []string {
	servicesForRegionMap := map[string][]string{}
	servicemap := NewServiceMap()
	serviceData := parseJson(servicemap.DownloadJson)

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

func serviceEntryContains(element serviceEntry, array []serviceEntry) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func (*AwsServiceMap) IsServiceInRegion(reqService string, reqRegion string) bool {
	//servicesForRegionMap := map[string][]string{}
	servicemap := NewServiceMap()
	serviceData := parseJson(servicemap.DownloadJson)

	reqPair := serviceEntry{ID: fmt.Sprintf("%s:%s", reqService, reqRegion)}
	if serviceEntryContains(reqPair, serviceData.ServiceEntries) {
		return true
	} else {
		return false
	}
}
