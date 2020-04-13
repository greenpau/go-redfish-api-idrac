// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
)

type computerSystemResponse struct {
	ODataAnnotation
	ID           string `yaml:"Id" json:"Id" xml:"Id"`
	UUID         string
	Name         string
	AssetTag     string
	BiosVersion  string
	Manufacturer string
	Model        string
	PartNumber   string
	SKU          string
	SerialNumber string
	SystemType   string
	Description  string
	//BootOrderCounter     uint64 `yaml:"BootOrder@odata.count" json:"BootOrder@odata.count" xml:"BootOrder@odata.count"`
	//HostingRolesCounter  uint64 `yaml:"HostingRoles@odata.count" json:"HostingRoles@odata.count" xml:"HostingRoles@odata.count"`
	//ChassisCounter       uint64 `yaml:"Chassis@odata.count" json:"Chassis@odata.count" xml:"Chassis@odata.count"`
	//CooledByCounter      uint64 `yaml:"CooledBy@odata.count" json:"CooledBy@odata.count" xml:"CooledBy@odata.count"`
	//ManagedByCounter     uint64 `yaml:"ManagedBy@odata.count" json:"ManagedBy@odata.count" xml:"ManagedBy@odata.count"`
	//PoweredByCounter     uint64 `yaml:"PoweredBy@odata.count" json:"PoweredBy@odata.count" xml:"PoweredBy@odata.count"`
	PCIeDevicesCounter   uint64 `yaml:"PCIeDevices@odata.count" json:"PCIeDevices@odata.count" xml:"PCIeDevices@odata.count"`
	PCIeFunctionsCounter uint64 `yaml:"PCIeFunctions@odata.count" json:"PCIeFunctions@odata.count" xml:"PCIeFunctions@odata.count"`
}

type computerSystemCounters struct {
	//BootOrder     uint64 `yaml:"boot_order" json:"boot_order" xml:"boot_order"`
	//HostingRoles  uint64 `yaml:"hosting_roles" json:"hosting_roles" xml:"hosting_roles"`
	//Chassis       uint64 `yaml:"chassis" json:"chassis" xml:"chassis"`
	//CooledBy      uint64 `yaml:"cooled_by" json:"cooled_by" xml:"cooled_by"`
	//ManagedBy     uint64 `yaml:"managed_by" json:"managed_by" xml:"managed_by"`
	//PoweredBy     uint64 `yaml:"powered_by" json:"powered_by" xml:"powered_by"`
	PCIeDevices   uint64 `yaml:"pcie_devices" json:"pcie_devices" xml:"pcie_devices"`
	PCIeFunctions uint64 `yaml:"pcie_functions" json:"pcie_functions" xml:"pcie_functions"`
}

// ComputerSystem represents an instance of Redfish ComputerSystem.
type ComputerSystem struct {
	ID           string                  `yaml:"id" json:"id" xml:"id"`
	OData        *ODataAnnotation        `yaml:"odata" json:"odata" xml:"odata"`
	BiosVersion  string                  `yaml:"bios_version" json:"bios_version" xml:"bios_version"`
	Manufacturer string                  `yaml:"manufacturer" json:"manufacturer" xml:"manufacturer"`
	Model        string                  `yaml:"model" json:"model" xml:"model"`
	PartNumber   string                  `yaml:"part_number" json:"part_number" xml:"part_number"`
	SKU          string                  `yaml:"sku" json:"sku" xml:"sku"`
	Description  string                  `yaml:"description" json:"description" xml:"description"`
	AssetTag     string                  `yaml:"asset_tag" json:"asset_tag" xml:"asset_tag"`
	Name         string                  `yaml:"name" json:"name" xml:"name"`
	SerialNumber string                  `yaml:"serial_number" json:"serial_number" xml:"serial_number"`
	SystemType   string                  `yaml:"system_type" json:"system_type" xml:"system_type"`
	UUID         string                  `yaml:"uuid" json:"uuid" xml:"uuid"`
	Counters     *computerSystemCounters `yaml:"counters" json:"counters" xml:"counters"`
}

// GetComputerSystemByResourceID returns an instance of Redfish ComputerSystem.
func (cli *Client) GetComputerSystemByResourceID(s string) (*ComputerSystem, error) {
	resp, err := cli.callAPI("GET", "", s, []byte{})
	if err != nil {
		return nil, err
	}
	return newComputerSystemFromBytes(resp)
}

// newComputerSystemFromString returns ComputerSystem instance from an input string.
func newComputerSystemFromString(s string) (*ComputerSystem, error) {
	return newComputerSystemFromBytes([]byte(s))
}

// newComputerSystemFromBytes returns ComputerSystem instance from an input byte array.
func newComputerSystemFromBytes(s []byte) (*ComputerSystem, error) {
	cs := &ComputerSystem{}
	response := &computerSystemResponse{}
	err := json.Unmarshal(s, response)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s, server response: %s", err, string(s[:]))
	}
	cs.OData = &ODataAnnotation{
		Context: response.Context,
		ID:      response.ID,
		Type:    response.Type,
	}
	cs.ID = response.ID
	cs.BiosVersion = response.BiosVersion
	cs.Manufacturer = response.Manufacturer
	cs.Model = response.Model
	cs.PartNumber = response.PartNumber
	cs.SKU = response.SKU
	cs.Description = response.Description
	cs.AssetTag = response.AssetTag
	cs.Name = response.Name
	cs.SerialNumber = response.SerialNumber
	cs.SystemType = response.SystemType
	cs.UUID = response.UUID
	cs.Counters = &computerSystemCounters{}
	//cs.Counters.BootOrder = response.BootOrderCounter
	//cs.Counters.HostingRoles = response.HostingRolesCounter
	//cs.Counters.Chassis = response.ChassisCounter
	//cs.Counters.CooledBy = response.CooledByCounter
	//cs.Counters.ManagedBy = response.ManagedByCounter
	//cs.Counters.PoweredBy = response.PoweredByCounter
	cs.Counters.PCIeDevices = response.PCIeDevicesCounter
	cs.Counters.PCIeFunctions = response.PCIeFunctionsCounter
	return cs, nil
}
