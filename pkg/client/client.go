// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// ReceiverDataLimit is the limit of data in bytes the client will read
// from a server.
const ReceiverDataLimit int64 = 1e6

// CliOperation represents supported command line operations
type CliOperation struct {
	Name        string
	Description string
}

// Client is an instance of iDRAC Redfish API client.
type Client struct {
	url                string
	host               string
	port               int
	protocol           string
	username           string
	password           string
	validateServerCert bool
	rootPath           string
	dataLimit          int64
}

// NewClient returns an instance of Client.
func NewClient() *Client {
	return &Client{
		dataLimit: ReceiverDataLimit,
		port:      443,
		protocol:  "https",
		rootPath:  "/redfish/v1/",
	}
}

func (cli *Client) rebaseURL() {
	if (cli.protocol == "https" && cli.port == 443) ||
		(cli.protocol == "http" && cli.port == 80) {
		cli.url = fmt.Sprintf("%s://%s", cli.protocol, cli.host)
		return
	}
	cli.url = fmt.Sprintf("%s://%s:%d", cli.protocol, cli.host, cli.port)
	return
}

// SetHost sets the target host for the API calls.
func (cli *Client) SetHost(s string) error {
	if s == "" {
		return fmt.Errorf("empty hostname or ip address")
	}
	cli.host = s
	cli.rebaseURL()
	return nil
}

// SetPort sets the port number for the API calls.
func (cli *Client) SetPort(p int) error {
	if p == 0 {
		return fmt.Errorf("invalid port: %d", p)
	}
	cli.port = p
	cli.rebaseURL()
	return nil
}

// SetUsername sets the username for the API calls.
func (cli *Client) SetUsername(s string) error {
	if s == "" {
		return fmt.Errorf("empty username")
	}
	cli.username = s
	return nil
}

// SetPassword sets the password for the API calls.
func (cli *Client) SetPassword(s string) error {
	if s == "" {
		return fmt.Errorf("empty password")
	}
	cli.password = s
	return nil
}

// SetProtocol sets the protocol for the API calls.
func (cli *Client) SetProtocol(s string) error {
	switch s {
	case "http":
		cli.protocol = s
	case "https":
		cli.protocol = s
	default:
		return fmt.Errorf("supported protocols: http, https; unsupported protocol: %s", s)
	}
	cli.rebaseURL()
	return nil
}

// SetValidateServerCertificate instructs the client to enforce the validation of certificates
// and check certificate errors.
func (cli *Client) SetValidateServerCertificate() error {
	cli.validateServerCert = true
	return nil
}

// GetOperations returns the names of available operations.
func (cli *Client) GetOperations() map[string]*CliOperation {
	operations := make(map[string]*CliOperation)
	operations["get-info"] = &CliOperation{
		Name:        "get-info",
		Description: "Get basic information about a remote Redfish API endpoint",
	}
	operations["get-systems"] = &CliOperation{
		Name:        "get-info",
		Description: "Get information about computer systems exposed via Redfish API",
	}
	return operations
}

func (cli *Client) callAPI(method string, contentType string, urlPath string, payload []byte) ([]byte, error) {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	url := fmt.Sprintf("%s%s", cli.url, urlPath)
	log.Debugf("%s request to %s", method, url)
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if !cli.validateServerCert {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}

	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(cli.username, cli.password)

	res, err := httpClient.Do(req)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "EOF") {
			return nil, err
		}
	}
	if res == nil {
		return nil, fmt.Errorf("response: <nil>, verify url: %s", url)
	}
	defer res.Body.Close()

	log.Debugf("API Server responded with %s", res.Status)

	dataLimiter := io.LimitReader(res.Body, cli.dataLimit)
	body, err := ioutil.ReadAll(dataLimiter)
	if err != nil {
		return nil, fmt.Errorf("non-EOF error at url %s: %s", url, err)
	}

	switch res.StatusCode {
	case 200:
		return body, nil
	default:
		return nil, fmt.Errorf("error: status code %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
