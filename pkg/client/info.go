// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
)

type protocolFeatures struct {
	ExcerptQuery    bool
	FilterQuery     bool
	OnlyMemberQuery bool
	SelectQuery     bool
	ExpandQuery     expandQueryProtocolFeatures
}

type expandQueryProtocolFeatures struct {
	ExpandAll bool
	Levels    bool
	Links     bool
	MaxLevels uint64
	NoLinks   bool
}

type linkedAPIPath struct {
	Sessions ODataAnnotation
}

type infoResponse struct {
	ODataAnnotation
	AccountService     ODataAnnotation
	CertificateService ODataAnnotation
	Chassis            ODataAnnotation
	Description        string
	EventService       ODataAnnotation
	Fabrics            ODataAnnotation
	ID                 string `yaml:"Id" json:"Id" xml:"Id"`
	JobService         ODataAnnotation
	JSONSchemas        ODataAnnotation `yaml:"JsonSchemas" json:"JsonSchemas" xml:"JsonSchemas"`
	Links              linkedAPIPath
	Managers           ODataAnnotation
	Name               string
	Product            string
	Features           protocolFeatures
	RedfishVersion     string
	Registries         ODataAnnotation
	SessionService     ODataAnnotation
	Systems            ODataAnnotation
	Tasks              ODataAnnotation
	TelemetryService   ODataAnnotation
	UpdateService      ODataAnnotation
	Oem                struct {
		Dell struct {
			IsBranded         int
			ManagerMACAddress string
			ServiceTag        string
		}
	}
}

// Info contains system information. The information in the structure
// is from querying Root service.
type Info struct {
	OData             *ODataAnnotation `yaml:"odata_annotations" json:"odata_annotations" xml:"odata_annotations"`
	Product           string           `yaml:"product" json:"product" xml:"product"`
	ServiceTag        string           `yaml:"service_tag" json:"service_tag" xml:"service_tag"`
	ManagerMACAddress string           `yaml:"manager_mac_address" json:"manager_mac_address" xml:"manager_mac_address"`
	RedfishVersion    string           `yaml:"redfish_version" json:"redfish_version" xml:"redfish_version"`
}

// GetInfo returns basic information about a system
func (cli *Client) GetInfo() (*Info, error) {
	resp, err := cli.callAPI("GET", "", cli.rootPath, []byte{})
	if err != nil {
		return nil, err
	}
	return newInfoFromBytes(resp)
}

// newInfoFromString returns Info instance from an input string.
func newInfoFromString(s string) (*Info, error) {
	return newInfoFromBytes([]byte(s))
}

// newInfoFromBytes returns Info instance from an input byte array.
func newInfoFromBytes(s []byte) (*Info, error) {
	info := &Info{}
	response := &infoResponse{}
	err := json.Unmarshal(s, response)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s, server response: %s", err, string(s[:]))
	}
	info.Product = response.Product
	info.ServiceTag = response.Oem.Dell.ServiceTag
	info.ManagerMACAddress = response.Oem.Dell.ManagerMACAddress
	info.RedfishVersion = response.RedfishVersion
	info.OData = &ODataAnnotation{
		Context: response.Context,
		ID:      response.ID,
		Type:    response.Type,
	}

	if response.RedfishVersion == "" {
		return nil, fmt.Errorf("Error parsing the received response: %s", s)
	}
	return info, nil
}
