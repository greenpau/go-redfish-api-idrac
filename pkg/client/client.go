// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// Client is an instance of iDRAC Redfish API client.
type Client struct {
	url      string
	host     string
	port     int
	protocol string
	username string
	password string
	secure   bool
	rootPath string
}

// NewClient returns an instance of Client.
func NewClient() *Client {
	return &Client{
		port:     443,
		protocol: "https",
		rootPath: "/redfish/v1/",
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

// SetSecure instructs the client to enforce the validation of certificates
// and check certificate errors.
func (cli *Client) SetSecure() error {
	cli.secure = true
	return nil
}

// GetOperations returns the names of available operations.
func (cli *Client) GetOperations() map[string]string {
	operations := make(map[string]string)
	operations["info"] = "Get basic information about a remote API endpoint"
	return operations
}

// GetInfo returns basic information about a system
func (cli *Client) GetInfo() (*Info, error) {
	resp, err := cli.callAPI("GET", "", cli.rootPath, []byte{})
	if err != nil {
		return nil, err
	}
	return NewInfoFromBytes(resp)
}

// GetResource return raw output from Redfish API server for a specific resource (path)
func (cli *Client) GetResource(s string) (*Resource, error) {
	resp, err := cli.callAPI("GET", "", s, []byte{})
	if err != nil {
		return nil, err
	}
	return NewResourceFromBytes(resp)
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
	if !cli.secure {
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
	if len(payload) == 0 {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	}

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if err.Error() != "EOF" {
			return nil, err
		}
	}

	switch res.StatusCode {
	case 200:
		return body, nil
	default:
		return nil, fmt.Errorf("error: status code %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
