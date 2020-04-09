// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseInfoJsonOutput(t *testing.T) {
	testFailed := 0
	dataDir := "../../assets/responses"
	for i, test := range []struct {
		input      string
		exp        *Info
		shouldFail bool // Whether test should result in a failure
		shouldErr  bool // Whether parsing of a response should result in error
	}{
		{
			input: "root_1",
			exp: &Info{
				Product:           "Integrated Dell Remote Access Controller",
				ServiceTag:        "24A8VC9",
				ManagerMACAddress: "eb:f2:49:84:66:d4",
				RedfishVersion:    "1.6.0",
			},
			shouldFail: false,
			shouldErr:  false,
		},
		{
			input: "root_1",
			exp: &Info{
				Product:           "Integrated Dell Remote Access Controller",
				ServiceTag:        "24A8VC9",
				ManagerMACAddress: "eb:f2:49:84:66:d5",
				RedfishVersion:    "1.6.0",
			},
			shouldFail: true,
			shouldErr:  false,
		},
		{
			input: "root_1",
			exp: &Info{
				Product:           "Integrated Dell Remote Access Controller",
				ServiceTag:        "24A8VC8",
				ManagerMACAddress: "eb:f2:49:84:66:d4",
				RedfishVersion:    "1.6.0",
			},
			shouldFail: true,
			shouldErr:  false,
		},
		{
			input: "root_1",
			exp: &Info{
				Product:           "Non-Integrated Dell Remote Access Controller",
				ServiceTag:        "24A8VC9",
				ManagerMACAddress: "eb:f2:49:84:66:d4",
				RedfishVersion:    "1.6.0",
			},
			shouldFail: true,
			shouldErr:  false,
		},
		{
			input: "root_1",
			exp: &Info{
				Product:           "Integrated Dell Remote Access Controller",
				ServiceTag:        "24A8VC9",
				ManagerMACAddress: "eb:f2:49:84:66:d4",
				RedfishVersion:    "1.6.1",
			},
			shouldFail: true,
			shouldErr:  false,
		},
		{
			input:      "root_2",
			exp:        nil,
			shouldFail: false,
			shouldErr:  true,
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
		info, err := NewInfoFromBytes(content)
		if err != nil {
			if !test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but threw error: %v", i, fp, err)
				testFailed++
				continue
			}
		} else {
			if test.shouldErr {
				t.Logf("FAIL: Test %d: input '%s', expected to throw error, but passed: %v", i, fp, *info)
				testFailed++
				continue
			}
		}

		if err == nil {
			if (info.Product != test.exp.Product) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "Product", info.Product, test.exp.Product)
				testFailed++
				continue
			}

			if (info.ServiceTag != test.exp.ServiceTag) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "ServiceTag", info.ServiceTag, test.exp.ServiceTag)
				testFailed++
				continue
			}

			if (info.ManagerMACAddress != test.exp.ManagerMACAddress) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "ManagerMACAddress", info.ManagerMACAddress, test.exp.ManagerMACAddress)
				testFailed++
				continue
			}
			if (info.RedfishVersion != test.exp.RedfishVersion) && !test.shouldFail {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but failed due to mismatch in '%s' field: '%v' (actual) vs. '%v' (expected)",
					i, fp, "RedfishVersion", info.RedfishVersion, test.exp.RedfishVersion)
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
