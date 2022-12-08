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

func (m *AwsServiceMap) parseJson() (regionalServiceData, error) {
	var serviceData regionalServiceData
	var err error

	if m.DownloadJson {
		res, err := http.Get("https://api.regional-table.region-services.aws.a2z.com/index.json")
		if err != nil {
			return serviceData, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return serviceData, err
		}
		json.Unmarshal([]byte(body), &serviceData)
		if err != nil {
			return serviceData, err
		}

	} else {
		jsonFile, err := awsJson.ReadFile("data/aws-service-regions.json")
		json.Unmarshal([]byte(jsonFile), &serviceData)
		if err != nil {
			return serviceData, err
		}
	}

	return serviceData, err
}

func (m *AwsServiceMap) GetAllRegions() ([]string, error) {
	totalRegions := []string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return totalRegions, err
	}

	for _, id := range serviceData.ServiceEntries {
		region := strings.Split(id.ID, ":")[1]
		if !contains(region, totalRegions) {
			totalRegions = append(totalRegions, region)
		}
	}
	return totalRegions, err

}

func (m *AwsServiceMap) GetRegionsForService(reqService string) ([]string, error) {
	regionsForServiceMap := map[string][]string{}

	serviceData, err := m.parseJson()
	if err != nil {
		return regionsForServiceMap[reqService], err
	}
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
	return regionsForServiceMap[reqService], err

}

func (m *AwsServiceMap) GetAllServices() ([]string, error) {
	totalServices := []string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return totalServices, err
	}
	for _, id := range serviceData.ServiceEntries {
		service := strings.Split(id.ID, ":")[0]
		if !contains(service, totalServices) {
			totalServices = append(totalServices, service)
		}
	}
	return totalServices, err

}

func (m *AwsServiceMap) GetServicesForRegion(reqRegion string) ([]string, error) {
	servicesForRegionMap := map[string][]string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return servicesForRegionMap[reqRegion], err
	}

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
	return servicesForRegionMap[reqRegion], err
}

func serviceEntryContains(element serviceEntry, array []serviceEntry) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func (m *AwsServiceMap) IsServiceInRegion(reqService string, reqRegion string) (bool, error) {
	//servicesForRegionMap := map[string][]string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return false, err
	}

	reqPair := serviceEntry{ID: fmt.Sprintf("%s:%s", reqService, reqRegion)}
	if serviceEntryContains(reqPair, serviceData.ServiceEntries) {
		return true, err
	} else {
		return false, err
	}
}
