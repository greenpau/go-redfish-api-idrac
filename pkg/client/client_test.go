// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	var timerStartTime time.Time
	logLevel, _ := log.ParseLevel("debug")
	log.SetLevel(logLevel)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var err error
		var fp string
		var fc []byte
		dataDir := "../../assets/responses"
		pathMap := map[string]string{
			"/redfish/v1":  "root_1.json",
			"/redfish/v1/": "root_1.json",
		}

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

		respFileName, respFileExists := pathMap[req.URL.Path]

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
	server := httptest.NewServer(mux)
	defer server.Close()

	srv := strings.Split(server.URL, ":")
	proto := srv[0]
	port, _ := strconv.Atoi(srv[2])

	t.Logf("Server URL: %s", server.URL)

	cli := NewClient()
	cli.SetHost("127.0.0.1")
	cli.SetPort(port)
	cli.SetProtocol(proto)
	cli.SetUsername("admin")
	cli.SetPassword("badSecret")

	t.Logf("client: testing GetInfo() with bad credentials")
	timerStartTime = time.Now()
	if _, err := cli.GetInfo(); err != nil {
		t.Logf("client: %s", err)
	} else {
		t.Fatalf("expected failure, but got non-error response")
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

	t.Logf("client: testing cli.GetResource(/redfish/v1/)")
	timerStartTime = time.Now()
	infoResource, err := cli.GetResource("/redfish/v1/")
	if err != nil {
		t.Fatalf("client: %s", err)
	}
	t.Logf("client: Response: %s\n", infoResource)
	t.Logf("client: took %s", time.Since(timerStartTime))

}
