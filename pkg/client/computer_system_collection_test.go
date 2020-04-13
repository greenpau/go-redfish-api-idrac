// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	. "github.com/greenpau/go-idrac-redfish-api/internal/client"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestGetComputerSystems(t *testing.T) {
	testFailed := 0
	var timerStartTime time.Time

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

	computerSystemCollection, err := cli.GetComputerSystemCollection()
	if err != nil {
		t.Fatalf("%s", err)
	}

	complianceMessages, compliant := isStructCompliant(computerSystemCollection)
	if !compliant {
		testFailed++
	}

	for _, entry := range complianceMessages {
		t.Logf("%s", entry)
	}

	computerSystems, err := cli.GetComputerSystems()
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("Number of Computer Systems: %d\n", len(computerSystems))
	t.Logf("---------------------------------\n")
	for _, cs := range computerSystems {
		if cs.Manufacturer != "" {
			t.Logf("System: %s | Manufacturer: %s\n", cs.ID, cs.Manufacturer)
		}
		if cs.Model != "" {
			t.Logf("System: %s | Model: %s\n", cs.ID, cs.Model)
		}
		if cs.SKU != "" {
			t.Logf("System: %s | SKU: %s\n", cs.ID, cs.SKU)
		}
		if cs.PartNumber != "" {
			t.Logf("System: %s | SKU: %s\n", cs.ID, cs.PartNumber)
		}
		if cs.BiosVersion != "" {
			t.Logf("System: %s | BIOS Version: %s\n", cs.ID, cs.BiosVersion)
		}
	}
	t.Logf("client: took %s", time.Since(timerStartTime))
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
