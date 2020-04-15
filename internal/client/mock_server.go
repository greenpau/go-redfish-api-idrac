package internal

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

// MockTestServerInstance is an instance of a mock web server.
type MockTestServerInstance struct {
	Instance *httptest.Server
	URL      *url.URL
	Hostname string
	Protocol string
	Port     int
}

// MockTestServer is a mock web server. The server supports both HTTPS and HTTP.
type MockTestServer struct {
	NonTLS *MockTestServerInstance
	TLS    *MockTestServerInstance
}

// Close closes running instances of MockTestServerInstance, if any.
func (srv *MockTestServer) Close() {
	if srv.NonTLS.Instance != nil {
		srv.NonTLS.Instance.Close()
	}
	if srv.TLS.Instance != nil {
		srv.TLS.Instance.Close()
	}
}

// NewMockTestServer return an instance of MockTestServer running
// with and without TLS.
func NewMockTestServer(pathMap map[string]string, tlsEnabled bool) (*MockTestServer, error) {
	// Create web server instance
	mts := &MockTestServer{
		NonTLS: &MockTestServerInstance{},
		TLS:    &MockTestServerInstance{},
	}
	serverEndpoints := map[string]string{
		"/redfish/v1/":                           "root_1.json",
		"/redfish/v1/Systems/":                   "computer_system_collection_1.json",
		"/redfish/v1/Systems/System.Embedded.1/": "computer_system_1.json",
		"/redfish/v1/Systems/System.Embedded.1/NetworkInterfaces/":                                                           "network_interface_collection_1.json",
		"/redfish/v1/Systems/System.Embedded.1/NetworkInterfaces/NIC.Integrated.1":                                           "network_interface_integrated_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/":                               "network_port_collection_integrated_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/":                                            "network_adapter_integrated_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkDeviceFunctions/":                     "network_device_function_collection_integrated_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkDeviceFunctions/NIC.Integrated.1-1-1": "network_device_function_integrated_1_1_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkDeviceFunctions/NIC.Integrated.1-2-1": "network_device_function_integrated_1_2_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkDeviceFunctions/NIC.Integrated.1-3-1": "network_device_function_integrated_1_3_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkDeviceFunctions/NIC.Integrated.1-4-1": "network_device_function_integrated_1_4_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Integrated.1-1":             "network_port_integrated_1_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Integrated.1-2":             "network_port_integrated_1_2.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Integrated.1-3":             "network_port_integrated_1_3.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Integrated.1-4":             "network_port_integrated_1_4.json",
		"/redfish/v1/Systems/System.Embedded.1/NetworkInterfaces/NIC.Slot.2":                                                 "network_interface_slot_2.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Slot.2/":                                                  "network_adapter_slot_2.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Slot.2/NetworkDeviceFunctions/":                           "network_device_function_collection_slot_2.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Slot.2/NetworkDeviceFunctions/NIC.Slot.2-1-1":             "network_device_function_slot_2_1_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Slot.2/NetworkDeviceFunctions/NIC.Slot.2-2-1":             "network_device_function_slot_2_2_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Slot.2/NetworkPorts/":                                     "network_port_collection_slot_2.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Slot.2-1":                   "network_port_slot_2_1.json",
		"/redfish/v1/Chassis/System.Embedded.1/NetworkAdapters/NIC.Integrated.1/NetworkPorts/NIC.Slot.2-2":                   "network_port_slot_2_2.json",
	}

	if pathMap != nil {
		for k, v := range pathMap {
			if _, exists := serverEndpoints[k]; !exists {
				serverEndpoints[k] = v
			}
		}
	}

	dataDir := "../../assets/responses"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var err error
		var fp string
		var fc []byte
		isAuthError := true
		authHeader := req.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Basic ") {
			authHeader = strings.TrimLeft(authHeader, "Basic")
			authHeader = strings.TrimSpace(authHeader)
			if b, err := base64.StdEncoding.DecodeString(authHeader); err == nil {
				if string(b) == "admin:secret" {
					isAuthError = false
				}
			}
		}

		if isAuthError {
			fp = fmt.Sprintf("%s/access_denied_error_1.json", dataDir)
			fc, err = ioutil.ReadFile(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(fc), http.StatusNotFound)
			return
		}

		if req.Method != "GET" {
			http.Error(w, "Bad Request, expecting GET", http.StatusBadRequest)
			return
		}

		if strings.HasSuffix(req.URL.Path, "/empty_response") {
			panic("")
		}

		if strings.HasSuffix(req.URL.Path, "/replay_request") {
			reqBlob, _ := httputil.DumpRequest(req, true)
			w.Write(reqBlob)
			return
		}

		respFileName, respFileExists := serverEndpoints[req.URL.Path]
		if !respFileExists {
			fp = fmt.Sprintf("%s/not_found_error_1.json", dataDir)
			fc, err = ioutil.ReadFile(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(fc), http.StatusNotFound)
			return
		}

		fp = fmt.Sprintf("%s/%s", dataDir, respFileName)
		fc, err = ioutil.ReadFile(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(fc)
	})

	// Initialize HTTP server
	mts.NonTLS.Instance = httptest.NewServer(mux)
	log.Debugf("HTTP Server URL: %s", mts.NonTLS.Instance.URL)
	httpServerURL, err := url.Parse(mts.NonTLS.Instance.URL)
	if err != nil {
		return nil, err
	}
	mts.NonTLS.URL = httpServerURL
	mts.NonTLS.Hostname = httpServerURL.Hostname()
	log.Debugf("HTTP Server Hostname: %s", mts.NonTLS.Hostname)
	mts.NonTLS.Protocol = strings.Split(mts.NonTLS.Instance.URL, ":")[0]
	log.Debugf("HTTP Server Protocol: %s", mts.NonTLS.Protocol)
	httpServerPort, _ := strconv.Atoi(httpServerURL.Port())
	mts.NonTLS.Port = httpServerPort
	log.Debugf("HTTP Server Port: %d", mts.NonTLS.Port)

	if tlsEnabled {
		// Initialize HTTPS server
		mts.TLS.Instance = httptest.NewTLSServer(mux)
		log.Debugf("HTTPS Server URL: %s", mts.TLS.Instance.URL)
		httpsServerURL, err := url.Parse(mts.TLS.Instance.URL)
		if err != nil {
			return nil, err
		}
		mts.TLS.URL = httpsServerURL
		mts.TLS.Hostname = httpsServerURL.Hostname()
		log.Debugf("HTTPS Server Hostname: %s", mts.TLS.Hostname)
		mts.TLS.Protocol = strings.Split(mts.TLS.Instance.URL, ":")[0]
		log.Debugf("HTTPS Server Protocol: %s", mts.TLS.Protocol)
		httpsServerPort, _ := strconv.Atoi(httpsServerURL.Port())
		mts.TLS.Port = httpsServerPort
		log.Debugf("HTTPS Server Port: %d", mts.TLS.Port)
	}

	return mts, nil
}
