// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
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
