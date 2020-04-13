// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
	. "github.com/greenpau/go-idrac-redfish-api/internal/client"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseComputerSystemJsonOutput(t *testing.T) {
	testFailed := 0
	dataDir := "../../assets/responses"
	for i, test := range []struct {
		input      string
		exp        *ComputerSystem
		shouldFail bool // Whether test should result in a failure
		shouldErr  bool // Whether parsing of a response should result in error
	}{
		{
			input: "computer_system_1",
			exp: &ComputerSystem{
				ID: "System.Embedded.1",
				OData: NewODataAnnotation(
					"System.Embedded.1",
					"#ComputerSystem.v1_5_1.ComputerSystem",
					"/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
				),
				BiosVersion:  "2.4.8",
				Manufacturer: "Dell Inc.",
				Model:        "PowerEdge R640",
				PartNumber:   "0HG0B7V21",
				SKU:          "24A8VC9",
				Description:  "Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.",
				AssetTag:     "KN23N857Z",
				Name:         "System",
				SerialNumber: "CNCMU00201476Z",
				SystemType:   "Physical",
				UUID:         "73c016b0-76ff-4a58-a6f0-d22346f44046",
				//HostName: "xxx",
				//IndicatorLED: "xxx",
				Counters: &computerSystemCounters{
					PCIeDevices:   11,
					PCIeFunctions: 16,
				},
			},
			shouldFail: false,
			shouldErr:  false,
		},
	} {
		// Read response file
		fp := fmt.Sprintf("%s/%s.json", dataDir, test.input)
		content, err := ioutil.ReadFile(fp)
		if err != nil {
			t.Logf("FAIL: Test %d: failed reading '%s', error: %v", i, fp, err)
			testFailed++
			continue
		}

		// Parse API response
		resource, err := newComputerSystemFromBytes(content)
		if err != nil {
			if !test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but threw error: %v", i, fp, err)
				testFailed++
				continue
			}
		} else {
			if test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to throw error, but passed: %v", i, fp, *resource)
				testFailed++
				continue
			}
			// Parse API response from string
			computerSystemFromString, computerSystemFromStringError := newComputerSystemFromString(string(content))
			if computerSystemFromStringError != nil {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but got error: %v", i, fp, computerSystemFromStringError)
				testFailed++
				continue
			}
			if !reflect.DeepEqual(computerSystemFromString, resource) {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but got value mismatch: '%v' (newComputerSystemFromString) vs. '%v' (newComputerSystemFromBytes)",
					i, fp, *computerSystemFromString, *resource)
				testFailed++
				continue
			}
		}

		if err == nil {
			if (resource.ID != test.exp.ID) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "ID", resource.ID, test.exp.ID)
				testFailed++
				continue
			}
			if (resource.BiosVersion != test.exp.BiosVersion) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "BiosVersion", resource.BiosVersion, test.exp.BiosVersion)
				testFailed++
				continue
			}
			if (resource.Manufacturer != test.exp.Manufacturer) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "Manufacturer", resource.Manufacturer, test.exp.Manufacturer)
				testFailed++
				continue
			}
			if (resource.Model != test.exp.Model) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "Model", resource.Model, test.exp.Model)
				testFailed++
				continue
			}
			if (resource.PartNumber != test.exp.PartNumber) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "PartNumber", resource.PartNumber, test.exp.PartNumber)
				testFailed++
				continue
			}
			if (resource.SKU != test.exp.SKU) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "SKU", resource.SKU, test.exp.SKU)
				testFailed++
				continue
			}
			if (resource.Description != test.exp.Description) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "Description", resource.Description, test.exp.Description)
				testFailed++
				continue
			}
			if (resource.AssetTag != test.exp.AssetTag) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "AssetTag", resource.AssetTag, test.exp.AssetTag)
				testFailed++
				continue
			}
			if (resource.Name != test.exp.Name) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "Name", resource.Name, test.exp.Name)
				testFailed++
				continue
			}
			if (resource.SerialNumber != test.exp.SerialNumber) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "SerialNumber", resource.SerialNumber, test.exp.SerialNumber)
				testFailed++
				continue
			}
			if (resource.SystemType != test.exp.SystemType) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "SystemType", resource.SystemType, test.exp.SystemType)
				testFailed++
				continue
			}
			if (resource.UUID != test.exp.UUID) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "UUID", resource.UUID, test.exp.UUID)
				testFailed++
				continue
			}
			if !reflect.DeepEqual(resource.Counters, test.exp.Counters) && !test.shouldFail {

				//t.Logf("ComputerSystem.Counters.BootOrder: %d", resource.Counters.BootOrder)
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but got value mismatch in '%s' field: ''%v' (actual) vs. '%v' (expected)",
					i, fp, "Counters", *resource.Counters, *test.exp.Counters)
				testFailed++
				continue
			}

		}

		if test.shouldFail {
			t.Logf("PASS: Test %d: input '%s', expected to fail, failed", i, fp)
		} else {
			t.Logf("PASS: Test %d: input '%s', expected to pass, passed", i, fp)
		}
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}

func convertFieldToTag(s string) string {
	s = strings.Replace(s, "OData", "odata", -1)
	s = strcase.ToSnake(s)
	return s
}

func isParserCompliant(resource interface{}, expResourceMap map[string]interface{}) ([]string, bool) {
	result := true
	output := []string{}

	resourceInstance := reflect.TypeOf(resource).Elem()
	resourceName := fmt.Sprintf("%s", resourceInstance.Name())
	resourceMap := make(map[string]bool)

	for i := 0; i < resourceInstance.NumField(); i++ {
		tagName := "json"
		resourceField := resourceInstance.Field(i)
		tagValue := resourceField.Tag.Get(tagName)
		if tagValue == "" {
			tagValue = fmt.Sprintf("%s", resourceField.Name)
		}
		resourceMap[tagValue] = true
	}

	for k, v := range expResourceMap {
		if _, exists := resourceMap[k]; !exists {
			output = append(output, fmt.Sprintf(
				"    %s %s", k, v,
			))
			result = false
		}
	}

	if len(output) > 0 {
		output = append([]string{fmt.Sprintf(
			"Fix the following struct:\n\ntype %s struct {",
			resourceName,
		)}, output...)
	}

	if len(output) > 0 {
		return []string{strings.Join(output, "\n")}, result
	}

	return output, result
}

func isStructCompliant(resource interface{}) ([]string, bool) {
	result := true
	output := []string{}

	resourceInstance := reflect.TypeOf(resource).Elem()
	resourceType := fmt.Sprintf("%s", resourceInstance.Name())
	resourceKind := fmt.Sprintf("%s", resourceInstance.Kind())

	if resourceKind != "struct" {
		output = append(output, fmt.Sprintf(
			"FAIL: %s resource kind is unsupported",
			resourceKind,
		))
		return output, false
	}

	for i := 0; i < resourceInstance.NumField(); i++ {
		for _, tagName := range []string{"json", "xml", "yaml"} {
			resourceField := resourceInstance.Field(i)
			tagValue := resourceField.Tag.Get(tagName)
			if tagValue == "" {
				result = false
				output = append(output, fmt.Sprintf(
					"FAIL: %s tag not found in %s.%s.%s (%v)",
					tagName,
					resourceType,
					resourceInstance.Name(),
					resourceField.Name,
					resourceField.Type,
				))
				continue
			}
			expTagValue := convertFieldToTag(resourceField.Name)
			if tagValue != expTagValue {
				result = false
				output = append(output, fmt.Sprintf(
					"FAIL: %s tag mismatch found in %s.%s.%s (%v): %s (actual) vs. %s (expected)",
					tagName,
					resourceType,
					resourceInstance.Name(),
					resourceField.Name,
					resourceField.Type,
					tagValue,
					expTagValue,
				))
				continue

			}

			output = append(output, fmt.Sprintf(
				"PASS: %s tag is compliant %s.%s (%v), %s: %v",
				tagName,
				resourceType,
				resourceField.Name,
				resourceField.Type,
				tagName,
				tagValue,
			))
		}
	}

	return output, result
}

func TestComputerSystemStruct(t *testing.T) {
	var isFailedTest bool
	var complianceMessages []string
	var compliant bool
	var timerStartTime time.Time
	timerStartTime = time.Now()

	// Set DEBUG logging level
	logLevel, _ := log.ParseLevel("debug")
	log.SetLevel(logLevel)

	// Create web server instance
	endpoints := map[string]string{
		"/redfish/v1/":                           "root_1.json",
		"/redfish/v1/Systems/":                   "computer_system_collection_1.json",
		"/redfish/v1/Systems/System.Embedded.1/": "computer_system_1.json",
	}
	server, err := NewMockTestServer(endpoints, false)
	if err != nil {
		t.Fatalf("Failed to start mock test server: %s", err)
	}
	defer server.Close()

	// Initialize client
	cli := NewClient()
	cli.SetHost(server.NonTLS.Hostname)
	cli.SetPort(server.NonTLS.Port)
	cli.SetProtocol(server.NonTLS.Protocol)
	cli.SetUsername("admin")
	cli.SetPassword("secret")

	// Define path to ComputerSystem resource.
	resourcePath := "/redfish/v1/Systems/System.Embedded.1/"

	// Test tag compliance of ComputerSystem struct
	resource, err := cli.GetComputerSystemByResourceID(resourcePath)
	if err != nil {
		t.Fatalf("%s", err)
	}

	complianceMessages, compliant = isStructCompliant(resource)
	if !compliant {
		isFailedTest = true
	}
	for _, entry := range complianceMessages {
		t.Logf("%s", entry)
	}

	// Test tag compliance of computerSystemResponse struct
	resourceBytes, err := cli.callAPI("GET", "", resourcePath, []byte{})
	if err != nil {
		t.Fatalf("%s", err)
	}
	//t.Logf("%s", resourceBytes)
	rawResource := &computerSystemResponse{}
	if err := json.Unmarshal(resourceBytes, rawResource); err != nil {
		t.Fatalf("parsing computerSystemResponse error: %s", err)
	}

	var rawResourceMap map[string]interface{}
	if err := json.Unmarshal(resourceBytes, &rawResourceMap); err != nil {
		t.Fatalf("unmarshal computerSystemResponse error: %s", err)
	}

	complianceMessages, compliant = isParserCompliant(rawResource, rawResourceMap)
	if !compliant {
		isFailedTest = true
	}
	for _, entry := range complianceMessages {
		t.Logf("%s", entry)
	}

	t.Logf("client: took %s", time.Since(timerStartTime))
	if isFailedTest {
		t.Fatalf("Failed")
	}

}
