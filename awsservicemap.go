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

type serviceEntry struct {
	ID         string              `json:"id"`
	Attributes *serviceAttributes  `json:"attributes,omitempty"`
}

type serviceAttributes struct {
	Region      string `json:"aws:region"`
	ServiceName string `json:"aws:serviceName"`
	ServiceURL  string `json:"aws:serviceUrl"`
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

// extractServiceSlug extracts the service identifier from AWS service URL
// Example: "https://aws.amazon.com/ec2/" -> "ec2"
//          "https://aws.amazon.com/rds/mysql/" -> "rds"
//          "https://aws.amazon.com/systems-manager/" -> "systems-manager"
func extractServiceSlug(serviceURL string) string {
	if serviceURL == "" {
		return ""
	}

	// Remove trailing slash
	serviceURL = strings.TrimSuffix(serviceURL, "/")

	// Split by / and get the last non-empty part that's not the domain
	parts := strings.Split(serviceURL, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !strings.Contains(parts[i], "aws.amazon.com") && !strings.Contains(parts[i], "http") {
			return parts[i]
		}
	}
	return ""
}

type AwsServiceMap struct {
	JsonFileSource JsonFileSource
	cachedData     *regionalServiceData // Cache the parsed data to avoid repeated HTTP requests
}

// When using "EMBEDDED_IN_PACKAGE" this package does not make any external HTTP requests, but the data might be slightly out of date
// When using "DOWNLOAD_FROM_AWS" this package makes an external HTTP request, but you get real-time data.
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

	// Return cached data if it exists
	if m.cachedData != nil {
		return *m.cachedData, nil
	}

	if m.JsonFileSource == "DOWNLOAD_FROM_AWS" {
		res, err := http.Get("https://api.regional-table.region-services.aws.a2z.com/index.json")
		if err != nil {
			return serviceData, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return serviceData, err
		}

		err = json.Unmarshal([]byte(body), &serviceData)
		if err != nil {
			return serviceData, err
		}

	} else {
		jsonFile, err := awsJson.ReadFile("data/aws-service-regions.json")
		if err != nil {
			return serviceData, err
		}
		err = json.Unmarshal([]byte(jsonFile), &serviceData)
		if err != nil {
			return serviceData, err
		}
	}

	// Cache the data for subsequent calls
	m.cachedData = &serviceData

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
	regionsForService := []string{}

	serviceData, err := m.parseJson()
	if err != nil {
		return regionsForService, err
	}

	for _, entry := range serviceData.ServiceEntries {
		idParts := strings.Split(entry.ID, ":")
		if len(idParts) != 2 {
			continue
		}

		serviceHash := idParts[0]
		region := idParts[1]

		// Check if this is the service we're looking for
		var matches bool

		// Try old format first (direct match)
		if serviceHash == reqService {
			matches = true
		} else if entry.Attributes != nil && entry.Attributes.ServiceURL != "" {
			// New format: extract slug from URL
			serviceSlug := extractServiceSlug(entry.Attributes.ServiceURL)
			if serviceSlug == reqService {
				matches = true
			}
		}

		if matches && !contains(region, regionsForService) {
			regionsForService = append(regionsForService, region)
		}
	}

	return regionsForService, nil
}

// Returns a slice of strings that represent all observed services
func (m *AwsServiceMap) GetAllServices() ([]string, error) {
	totalServices := []string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return totalServices, err
	}

	for _, entry := range serviceData.ServiceEntries {
		var serviceName string

		if entry.Attributes != nil && entry.Attributes.ServiceURL != "" {
			// New format: extract from URL
			serviceName = extractServiceSlug(entry.Attributes.ServiceURL)
		} else {
			// Fallback to hash from ID (old format)
			idParts := strings.Split(entry.ID, ":")
			if len(idParts) > 0 {
				serviceName = idParts[0]
			}
		}

		if serviceName != "" && !contains(serviceName, totalServices) {
			totalServices = append(totalServices, serviceName)
		}
	}

	return totalServices, nil
}

// Returns a slice of strings that represent all service supported in a specific region
func (m *AwsServiceMap) GetServicesForRegion(reqRegion string) ([]string, error) {
	servicesForRegion := []string{}
	serviceData, err := m.parseJson()
	if err != nil {
		return servicesForRegion, err
	}

	for _, entry := range serviceData.ServiceEntries {
		idParts := strings.Split(entry.ID, ":")
		if len(idParts) != 2 {
			continue
		}

		serviceHash := idParts[0]
		region := idParts[1]

		if region != reqRegion {
			continue
		}

		// Determine service name
		var serviceName string
		if entry.Attributes != nil && entry.Attributes.ServiceURL != "" {
			serviceName = extractServiceSlug(entry.Attributes.ServiceURL)
		} else {
			// Fallback to hash (old format)
			serviceName = serviceHash
		}

		if serviceName != "" && !contains(serviceName, servicesForRegion) {
			servicesForRegion = append(servicesForRegion, serviceName)
		}
	}

	return servicesForRegion, nil
}

// Is a specific service supported in a specific region. Returns true/false
// Handles both old format (ec2:us-east-1) and new format (hash:us-east-1 with attributes)
func (m *AwsServiceMap) IsServiceInRegion(reqService string, reqRegion string) (bool, error) {
	serviceData, err := m.parseJson()
	if err != nil {
		return false, err
	}

	// Try direct match first (old format or if hash is provided)
	reqPair := serviceEntry{ID: fmt.Sprintf("%s:%s", reqService, reqRegion)}
	if serviceEntryContains(reqPair, serviceData.ServiceEntries) {
		return true, nil
	}

	// New format: check if any entry with matching region has the service slug
	for _, entry := range serviceData.ServiceEntries {
		// Check if region matches
		idParts := strings.Split(entry.ID, ":")
		if len(idParts) != 2 {
			continue
		}

		if idParts[1] != reqRegion {
			continue
		}

		// Extract service slug from URL if attributes exist
		if entry.Attributes != nil && entry.Attributes.ServiceURL != "" {
			serviceSlug := extractServiceSlug(entry.Attributes.ServiceURL)
			if serviceSlug == reqService {
				return true, nil
			}
		}
	}

	return false, nil
}

func serviceEntryContains(element serviceEntry, array []serviceEntry) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}
