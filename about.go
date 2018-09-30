package tyrgin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

func getAboutFieldValue(aboutConfigMap map[string]interface{}, key, aboutFilePath string) string {

	value, ok := aboutConfigMap[key]
	if !ok {
		Log(fmt.Sprintf("Field `%s`  missing from %s.", key, aboutFilePath))
		return AboutFieldNa
	}

	stringValue, ok := value.(string)
	if !ok {
		Log(fmt.Sprintf("Field `%s` is not a string in %s.", key, aboutFilePath))
		return AboutFieldNa
	}

	return stringValue
}

func getAboutFieldValues(aboutConfigMap map[string]interface{}, key, aboutFilePath string) []string {

	value, ok := aboutConfigMap[key]
	if !ok {
		Log(fmt.Sprintf("Field `%s` missing from %s.", key, aboutFilePath))
		return []string{}
	}

	interfaces, ok := value.([]interface{})
	if !ok {
		Log(fmt.Sprintf("Field `%s` is not an array in %s.", key, aboutFilePath))
		return []string{}
	}

	strings := make([]string, len(interfaces))
	for i := range interfaces {

		stringValue, ok := interfaces[i].(string)
		if !ok {
			strings[i] = AboutFieldNa
			Log(fmt.Sprintf("Field[%d] `%s` is not a String in %s.", i, key, aboutFilePath))
		} else {
			strings[i] = stringValue
		}

	}

	return strings
}

func getAboutCustomDataFieldValues(aboutConfigMap map[string]interface{}, aboutFilePath string) map[string]interface{} {

	value, ok := aboutConfigMap["customData"]
	if !ok {
		return nil
	}

	mapValue, ok := value.(map[string]interface{})
	if !ok {
		Log(fmt.Sprintf("Field `customData` is not a valid JSON object in %s.", aboutFilePath))
		return nil
	}

	return mapValue
}

func About(statusEndpoints []StatusEndpoint, protocol, aboutFilePath, versionFilePath string, customData map[string]interface{}) string {

	aboutData, _ := ioutil.ReadFile(aboutFilePath)

	// Initialize ConfigAbout with default values in case we have problems reading from the file
	aboutConfig := ConfigAbout{
		Id:          AboutFieldNa,
		Summary:     AboutFieldNa,
		Description: AboutFieldNa,
		Maintainers: []string{},
		ProjectRepo: AboutFieldNa,
		ProjectHome: AboutFieldNa,
		LogsLinks:   []string{},
		StatsLinks:  []string{},
	}

	// Unmarshal JSON into a generic object so we don't completely fail if one of the fields is invalid or missing
	var aboutConfigMap map[string]interface{}
	err := json.Unmarshal(aboutData, &aboutConfigMap)

	if err == nil {
		// Parse out each value individually
		aboutConfig.Id = getAboutFieldValue(aboutConfigMap, "id", aboutFilePath)
		aboutConfig.Summary = getAboutFieldValue(aboutConfigMap, "summary", aboutFilePath)
		aboutConfig.Description = getAboutFieldValue(aboutConfigMap, "description", aboutFilePath)
		aboutConfig.Maintainers = getAboutFieldValues(aboutConfigMap, "maintainers", aboutFilePath)
		aboutConfig.ProjectRepo = getAboutFieldValue(aboutConfigMap, "projectRepo", aboutFilePath)
		aboutConfig.ProjectHome = getAboutFieldValue(aboutConfigMap, "projectHome", aboutFilePath)
		aboutConfig.LogsLinks = getAboutFieldValues(aboutConfigMap, "logsLinks", aboutFilePath)
		aboutConfig.StatsLinks = getAboutFieldValues(aboutConfigMap, "statsLinks", aboutFilePath)
		aboutConfig.CustomData = getAboutCustomDataFieldValues(aboutConfigMap, aboutFilePath)
	} else {
		ErrorLogger(err, fmt.Sprintf("Error deserializing about data from %s. Error: %s JSON: %s", aboutFilePath, err.Error(), aboutData))
	}

	// Merge custom data from about.json with custom data passed in by client
	// and prefer values passed by client over values in about.json
	if customData != nil {
		if aboutConfig.CustomData == nil {
			aboutConfig.CustomData = make(map[string]interface{})
		}

		for key, value := range customData {
			aboutConfig.CustomData[key] = value
		}
	}

	// Extract version
	var version string
	versionData, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		ErrorLogger(err, fmt.Sprintf("Error reading version from %s. Error: %s", versionFilePath, err.Error()))
		version = VersionNa
	} else {
		version = strings.TrimSpace(string(versionData))
	}

	// Get hostname
	host, err := os.Hostname()
	if err != nil {
		ErrorLogger(err, fmt.Sprintf("Error getting hostname. Error: %s", err.Error()))
		host = "unknown"
	}

	aboutResponse := AboutResponse{
		Id:          aboutConfig.Id,
		Name:        aboutConfig.Summary,
		Description: aboutConfig.Description,
		Protocol:    protocol,
		Owners:      aboutConfig.Maintainers,
		Version:     version,
		Host:        host,
		ProjectRepo: aboutConfig.ProjectRepo,
		ProjectHome: aboutConfig.ProjectHome,
		LogsLinks:   aboutConfig.LogsLinks,
		StatsLinks:  aboutConfig.StatsLinks,
		CustomData:  aboutConfig.CustomData,
	}

	// Execute status checks async
	var wg sync.WaitGroup
	dc := make(chan dependencyPosition)
	wg.Add(len(statusEndpoints))

	for ie, se := range statusEndpoints {
		go func(s StatusEndpoint, i int) {
			start := time.Now()
			dependencyStatus := translateStatusList(s.StatusCheck.CheckStatus(s.Name))
			var elapsed float64 = float64(time.Since(start)) * 0.000000001
			dependency := Dependency{
				Name:           s.Name,
				Status:         dependencyStatus,
				StatusDuration: elapsed,
				StatusPath:     s.Slug,
				Type:           s.Type,
				IsTraversable:  s.IsTraversable,
			}

			dc <- dependencyPosition{
				item:     dependency,
				position: i,
			}
		}(se, ie)
	}

	// Collect our responses and put them in the right spot
	dependencies := make([]Dependency, len(statusEndpoints))
	go func() {
		for dp := range dc {
			dependencies[dp.position] = dp.item
			wg.Done()
		}
	}()

	// Wait until all async status checks are done and collected
	wg.Wait()
	close(dc)

	aboutResponse.Dependencies = dependencies

	aboutResponseJson, err := json.Marshal(aboutResponse)
	if err != nil {
		msg := fmt.Sprintf("Error serializing AboutResponse: %s", err)
		sl := StatusList{
			StatusList: []Status{
				{Description: "Invalid AboutResponse", Result: CRITICAL, Details: msg},
			},
		}
		return SerializeStatusList(sl)
	}

	return string(aboutResponseJson)
}
