// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	. "github.com/greenpau/go-idrac-redfish-api/internal/client"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	var timerStartTime time.Time

	// Set DEBUG logging level
	logLevel, _ := log.ParseLevel("debug")
	log.SetLevel(logLevel)

	// Create web server instance
	endpoints := map[string]string{
		"/redfish/v1/": "root_1.json",
	}
	server, err := NewMockTestServer(endpoints, true)
	if err != nil {
		t.Fatalf("Failed to start mock test server: %s", err)
	}
	defer server.Close()

	// Initialize client
	cli := NewClient()

	t.Logf("client: testing SetHost()")
	if err := cli.SetHost(""); err == nil {
		t.Fatalf("expected failure, but succeeded")
	}
	cli.SetHost(server.NonTLS.Hostname)

	t.Logf("client: testing SetPort()")
	if err := cli.SetPort(0); err == nil {
		t.Fatalf("expected failure, but succeeded")
	}
	cli.SetPort(server.NonTLS.Port)

	t.Logf("client: testing SetProtocol()")
	if err := cli.SetProtocol("ssh"); err == nil {
		t.Fatalf("expected failure, but succeeded")
	}
	if err := cli.SetProtocol("https"); err != nil {
		t.Fatalf("expected success, but failed")
	}
	cli.SetProtocol(server.NonTLS.Protocol)

	t.Logf("client: testing SetUsername()")
	if err := cli.SetUsername(""); err == nil {
		t.Fatalf("expected failure, but succeeded")
	}
	cli.SetUsername("admin")

	t.Logf("client: testing SetPassword()")
	if err := cli.SetPassword(""); err == nil {
		t.Fatalf("expected failure, but succeeded")
	}

	t.Logf("client: testing SetValidateServerCertificate()")
	if err := cli.SetValidateServerCertificate(); err != nil {
		t.Fatalf("expected success, but failed")
	}

	t.Logf("client: testing GetOperations()")
	for opName, opDescr := range cli.GetOperations() {
		t.Logf("client: operation %s => %s\n", opName, opDescr)
	}

	t.Logf("client: testing GetInfo() with bad credentials")
	cli.SetPassword("badSecret")
	timerStartTime = time.Now()
	if _, err := cli.GetInfo(); err != nil {
		t.Logf("client: %s", err)
	} else {
		t.Fatalf("client: expected failure, but got non-error response")
	}
	t.Logf("client: took %s", time.Since(timerStartTime))

	cli.SetPassword("secret")

	t.Logf("client: testing GetInfo()")
	timerStartTime = time.Now()
	info, err := cli.GetInfo()
	if err != nil {
		t.Fatalf("client: %s", err)
	}
	t.Logf("client: Product: %s\n", info.Product)
	t.Logf("client: Service Tag: %s\n", info.ServiceTag)
	t.Logf("client: Manager MAC Address: %s\n", info.ManagerMACAddress)
	t.Logf("client: Redfish API Version: %s\n", info.RedfishVersion)
	t.Logf("client: took %s", time.Since(timerStartTime))

	t.Logf("client: testing cli.GetResource()")
	timerStartTime = time.Now()
	infoResource, err := cli.GetResource("/redfish/v1/")
	if err != nil {
		t.Fatalf("client: expected success, but got error: %s", err)
	}
	t.Logf("client: GetResource() response: %s\n", infoResource)

	t.Logf("client: testing cli.GetResource() on non-existent dummy endpoint")
	if _, err := cli.GetResource("/redfish/v1/dummy"); err == nil {
		t.Fatalf("client: expected failure, but got non-error response")
	}

	t.Logf("client: testing cli.GetResource() on the path without leading slash")
	if _, err := cli.GetResource("redfish/v1/"); err != nil {
		t.Fatalf("client: expected failure, but got non-error response")
	}

	t.Logf("client: testing secure client")
	cli = NewClient()
	cli.SetHost(server.TLS.Hostname)
	cli.SetPort(server.TLS.Port)
	cli.SetProtocol(server.TLS.Protocol)
	cli.SetValidateServerCertificate()
	cli.SetUsername("admin")
	cli.SetPassword("secret")
	if _, err = cli.GetInfo(); err == nil {
		t.Fatalf("client: expected failure, but got non-error response")
	}
	cli.validateServerCert = false
	if _, err = cli.GetInfo(); err != nil {
		t.Fatalf("client: expected success, but got error: %s", err)
	}

	t.Logf("client: testing invalid HTTP method")
	if _, err := cli.callAPI("\\", "application/json", "/redfish/v1/", []byte{}); err == nil {
		t.Fatalf("client: expected failure, but got non-error response")
	}

	t.Logf("client: testing empty HTTP response")
	if _, err := cli.callAPI("GET", "application/json", "/redfish/v1/empty_response", []byte{}); err == nil {
		t.Fatalf("client: expected failure, but got non-error response")
	} else {
		if !strings.HasPrefix(err.Error(), "response: <nil>, verify url") {
			t.Fatalf("client: expected failure with 'response: <nil>, verify url', but got: %s", err)
		}
	}

	t.Logf("client: testing raw API calls")
	resp, err := cli.callAPI("GET", "application/json", "/redfish/v1/", []byte{})
	if err != nil {
		t.Fatalf("client: expected success, but got error: %s", err)
	}
	infoResource, err = newResourceFromBytes(resp)

	t.Logf("client: took %s", time.Since(timerStartTime))
}
