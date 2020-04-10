// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseResourceJsonOutput(t *testing.T) {
	testFailed := 0
	dataDir := "../../assets/responses"
	for i, test := range []struct {
		input      string
		shouldFail bool // Whether test should result in a failure
		shouldErr  bool // Whether parsing of a response should result in error
	}{
		{
			input:      "root_1",
			shouldFail: false,
			shouldErr:  false,
		},
		{
			input:      "root_2",
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
		resource, err := newResourceFromBytes(content)
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
			resourceFromString, resourceFromStringError := newResourceFromString(string(content))
			if resourceFromStringError != nil {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but got error: %v", i, fp, resourceFromStringError)
				testFailed++
				continue
			}
			if resourceFromString.String() != resource.String() {
				t.Logf("FAIL: Test %d: input '%s', expected to pass, but got value mismatch: newResourceFromString() vs. newResourceFromBytes()",
					i, fp)
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
