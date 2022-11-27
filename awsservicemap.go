package awsservicemap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type serviceEntry struct {
	ID string `json:"id"`
}

type regionalServiceData struct {
	ServiceEntries []serviceEntry `json:"prices"`
}

//var serviceMap map[string]map[string]string

func parseJson() regionalServiceData {

	jsonFile, err := os.Open("data/aws-service-regions.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var serviceData regionalServiceData
	json.Unmarshal([]byte(byteValue), &serviceData)
	return serviceData
}

func CreateServiceRegionMap() regionalServiceData {
	serviceData := parseJson()
	for _, id := range serviceData.ServiceEntries {
		service := strings.Split(id.ID, ":")[0]
		region := strings.Split(id.ID, ":")[1]
		fmt.Println(service, region)
	}
	return serviceData
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
