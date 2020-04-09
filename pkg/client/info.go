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
	Sessions apiPath
}

type infoResponse struct {
	AccountService     apiPath
	CertificateService apiPath
	Chassis            apiPath
	Description        string
	EventService       apiPath
	Fabrics            apiPath
	ID                 string `yaml:"Id" json:"Id" xml:"Id"`
	JobService         apiPath
	JSONSchemas        apiPath `yaml:"JsonSchemas" json:"JsonSchemas" xml:"JsonSchemas"`
	Links              linkedAPIPath
	Managers           apiPath
	Name               string
	Product            string
	Features           protocolFeatures
	RedfishVersion     string
	Registries         apiPath
	SessionService     apiPath
	Systems            apiPath
	Tasks              apiPath
	TelemetryService   apiPath
	UpdateService      apiPath
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
	Product           string `yaml:"product" json:"product" xml:"product"`
	ServiceTag        string `yaml:"service_tag" json:"service_tag" xml:"service_tag"`
	ManagerMACAddress string `yaml:"manager_mac_address" json:"manager_mac_address" xml:"manager_mac_address"`
	RedfishVersion    string `yaml:"redfish_version" json:"redfish_version" xml:"redfish_version"`
}

// NewInfoFromString returns Info instance from an input string.
func NewInfoFromString(s string) (*Info, error) {
	return NewInfoFromBytes([]byte(s))
}

// NewInfoFromBytes returns Info instance from an input byte array.
func NewInfoFromBytes(s []byte) (*Info, error) {
	info := &Info{}
	infoResponse := &infoResponse{}
	err := json.Unmarshal(s, infoResponse)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s, server response: %s", err, string(s[:]))
	}
	info.Product = infoResponse.Product
	info.ServiceTag = infoResponse.Oem.Dell.ServiceTag
	info.ManagerMACAddress = infoResponse.Oem.Dell.ManagerMACAddress
	info.RedfishVersion = infoResponse.RedfishVersion

	if infoResponse.RedfishVersion == "" {
		return nil, fmt.Errorf("Error parsing the received response: %s", s)
	}
	return info, nil
}
