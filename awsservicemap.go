package awsservicemap

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type AwsServiceMap struct {
	JsonFileSource JsonFileSource
}

type JsonFileSource string

// Enum values for JsonFileSource
const (
	JsonFileSourceEmbedded JsonFileSource = "EMBEDDED_IN_PACKAGE"
	JsonFileSourceDownload JsonFileSource = "DOWNLOAD_FROM_AWS"
)

// Values returns all values for JsonFileSource

func (JsonFileSource) Values() []JsonFileSource {
	return []JsonFileSource{
		"EMBEDDED_IN_PACKAGE",
		"DOWNLOAD_FROM_AWS",
	}
}

// Function that uses the constructor pattern. You can use this or instantiate your struct without it.
func NewServiceMap() *AwsServiceMap {
	return &AwsServiceMap{
		JsonFileSource: "EMBEDDED_IN_PACKAGE",
	}
}

// Takes the AWS JSON file with all of the services and regions and parses it into a single array of strings
// The strings in the array look like service:region and the other functions in this package all take that string
// and split it on the semicolon.
// Depending on the value of parseJson can AwsServiceMap.JsonFileSource, this function will either download the
// file from AWS on demand ("DOWNLOAD_FROM_AWS"), or use the file that is embedded in this package
// ("EMBEDDED_IN_PACKAGE"). THe file embedded in the package is kept up to date with a GitHub action that
// submit's a PR every time the AWS file changes.
// When using "EMBEDDED_IN_PACKAGE" this package does not make any external HTTP requests, but the data might be slightly out of date
// When using "DOWNLOAD_FROM_AWS" this package makes an external HTTP request, but you get real-time data.

func (m *AwsServiceMap) parseJson() (regionalServiceData, error) {
	var serviceData regionalServiceData
	var err error

	if m.JsonFileSource == "DOWNLOAD_FROM_AWS" {
		fmt.Println("user selected download")
		log.Fatalln("exit")

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

// Returns a slice of strings that represent all observed regions
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

// Returns a slice of strings that represent all regions that support the specific service
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

// Returns a slice of strings that represent all observed services

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

// Returns a slice of strings that represent all service supported in a specific region
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

// Is a specific service supported in a specific region. Returns true/false
func (m *AwsServiceMap) IsServiceInRegion(reqService string, reqRegion string) (bool, error) {
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

func serviceEntryContains(element serviceEntry, array []serviceEntry) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}
