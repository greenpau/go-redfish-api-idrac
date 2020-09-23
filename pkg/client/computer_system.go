// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
)

type computerSystemResponse struct {
	ODataAnnotation
	ID                   string `json:"Id"`
	UUID                 string
	Name                 string
	AssetTag             string
	BiosVersion          string
	Manufacturer         string
	Model                string
	PartNumber           string
	SKU                  string
	SerialNumber         string
	SystemType           string
	Description          string
	PCIeDevicesCounter   uint64 `json:"PCIeDevices@odata.count"`
	PCIeFunctionsCounter uint64 `json:"PCIeFunctions@odata.count"`
	HostingRolesCounter  uint64 `json:"HostingRoles@odata.count"`
	Status               HealthStatus
	HostName             string
	IndicatorLED         string
	PowerState           string
	ProcessorSummary     computerSystemProcessorSummary
	MemorySummary        computerSystemMemorySummary
	Actions              map[string]computerSystemActions

	// TODO: The below attributes are not in ComputerSystem struct

	SecureBoot         ODataAnnotation
	NetworkInterfaces  ODataAnnotation
	Storage            ODataAnnotation
	SimpleStorage      ODataAnnotation
	EthernetInterfaces ODataAnnotation
	Bios               ODataAnnotation
	Memory             ODataAnnotation
	Processors         ODataAnnotation
	PCIeDevices        []ODataAnnotation
	PCIeFunctions      []ODataAnnotation
	TrustedModules     []computerSystemTrustedModules
	HostWatchdogTimer  computerSystemHostWatchdogTimer
	HostingRoles       interface{}
	Boot               computerSystemBoot
	Links              computerSystemLinks

	Oem struct {
		Dell struct {
			DellSystem struct {
				ODataAnnotation
				BIOSReleaseDate                string
				BaseBoardChassisSlot           string
				BatteryRollupStatus            string
				BladeGeometry                  string
				CMCIP                          string
				CPURollupStatus                string
				ChassisModel                   string
				ChassisName                    string
				ChassisServiceTag              string
				ChassisSystemHeightUnit        uint64
				CurrentRollupStatus            string
				EstimatedExhaustTemperatureCel uint64
				EstimatedSystemAirflowCFM      uint64
				ExpressServiceCode             string
				FanRollupStatus                string
				IDSDMRollupStatus              string
				IntrusionRollupStatus          string
				IsOEMBranded                   string
				LastSystemInventoryTime        string
				LastUpdateTime                 string
				LicensingRollupStatus          string
				MaxCPUSockets                  uint64
				MaxDIMMSlots                   uint64
				MaxPCIeSlots                   uint64
				MemoryOperationMode            string
				NodeID                         string
				PSRollupStatus                 string
				PopulatedDIMMSlots             uint64
				PopulatedPCIeSlots             uint64
				PowerCapEnabledState           string
				SDCardRollupStatus             string
				SELRollupStatus                string
				ServerAllocationWatts          interface{}
				StorageRollupStatus            string
				SysMemErrorMethodology         string
				SysMemFailOverState            string
				SysMemLocation                 string
				SysMemPrimaryStatus            string
				SystemGeneration               string
				SystemID                       uint64
				SystemRevision                 string
				TempRollupStatus               string
				TempStatisticsRollupStatus     string
				UUID                           string
				VoltRollupStatus               string
				smbiosGUID                     string
			}
		}
	}
}

type computerSystemLinks struct {
	Chassis          []ODataAnnotation
	ChassisCounter   uint64 `yaml:"Chassis@odata.count" json:"Chassis@odata.count" xml:"Chassis@odata.count"`
	CooledBy         []ODataAnnotation
	CooledByCounter  uint64 `yaml:"CooledBy@odata.count" json:"CooledBy@odata.count" xml:"CooledBy@odata.count"`
	ManagedBy        []ODataAnnotation
	ManagedByCounter uint64 `yaml:"ManagedBy@odata.count" json:"ManagedBy@odata.count" xml:"ManagedBy@odata.count"`
	PoweredBy        []ODataAnnotation
	PoweredByCounter uint64 `yaml:"PoweredBy@odata.count" json:"PoweredBy@odata.count" xml:"PoweredBy@odata.count"`
	Oem              struct {
		Dell map[string]ODataAnnotation
	}
}

type computerSystemBoot struct {
	BootOptions                             ODataAnnotation
	BootOrder                               []string
	BootOrderCounter                        uint64 `json:"BootOrder@odata.count"`
	BootSourceOverrideEnabled               string
	BootSourceOverrideMode                  string
	BootSourceOverrideTarget                string
	UefiTargetBootSourceOverride            string
	BootSourceOverrideTargetAllowableValues []string `json:"BootSourceOverrideTarget@Redfish.AllowableValues"`
}

type computerSystemActions struct {
	Target        string   `json:"target"`
	AllowedValues []string `json:"ResetType@Redfish.AllowableValues"`
}

type computerSystemHostWatchdogTimer struct {
	FunctionEnabled bool
	TimeoutAction   string
	Status          HealthStatus
}

type computerSystemTrustedModules struct {
	FirmwareVersion string
	InterfaceType   string
	Status          HealthStatus
}

type computerSystemMemorySummary struct {
	MemoryMirroring      string
	TotalSystemMemoryGiB interface{} // uint64
	Status               HealthStatus
}

type computerSystemProcessorSummary struct {
	Count                 uint64
	LogicalProcessorCount uint64
	Model                 string
	Status                HealthStatus
}

type computerSystemCounters struct {
	BootOrder         uint64 `yaml:"boot_order" json:"boot_order" xml:"boot_order"`
	HostingRoles      uint64 `yaml:"hosting_roles" json:"hosting_roles" xml:"hosting_roles"`
	Chassis           uint64 `yaml:"chassis" json:"chassis" xml:"chassis"`
	CooledBy          uint64 `yaml:"cooled_by" json:"cooled_by" xml:"cooled_by"`
	ManagedBy         uint64 `yaml:"managed_by" json:"managed_by" xml:"managed_by"`
	PoweredBy         uint64 `yaml:"powered_by" json:"powered_by" xml:"powered_by"`
	PCIeDevices       uint64 `yaml:"pcie_devices" json:"pcie_devices" xml:"pcie_devices"`
	PCIeFunctions     uint64 `yaml:"pcie_functions" json:"pcie_functions" xml:"pcie_functions"`
	TotalProcessors   uint64 `yaml:"total_processors" json:"total_processors" xml:"total_processors"`
	LogicalProcessors uint64 `yaml:"logical_processors" json:"logical_processors" xml:"logical_processors"`
	TotalSystemMemory uint64 `yaml:"total_system_memory" json:"total_system_memory" xml:"total_system_memory"`
}

// ComputerSystem represents an instance of Redfish ComputerSystem.
type ComputerSystem struct {
	ID              string                          `yaml:"id" json:"id" xml:"id"`
	OData           *ODataAnnotation                `yaml:"odata" json:"odata" xml:"odata"`
	BiosVersion     string                          `yaml:"bios_version" json:"bios_version" xml:"bios_version"`
	Manufacturer    string                          `yaml:"manufacturer" json:"manufacturer" xml:"manufacturer"`
	Model           string                          `yaml:"model" json:"model" xml:"model"`
	PartNumber      string                          `yaml:"part_number" json:"part_number" xml:"part_number"`
	SKU             string                          `yaml:"sku" json:"sku" xml:"sku"`
	Description     string                          `yaml:"description" json:"description" xml:"description"`
	AssetTag        string                          `yaml:"asset_tag" json:"asset_tag" xml:"asset_tag"`
	Name            string                          `yaml:"name" json:"name" xml:"name"`
	SerialNumber    string                          `yaml:"serial_number" json:"serial_number" xml:"serial_number"`
	SystemType      string                          `yaml:"system_type" json:"system_type" xml:"system_type"`
	UUID            string                          `yaml:"uuid" json:"uuid" xml:"uuid"`
	Counters        *computerSystemCounters         `yaml:"counters" json:"counters" xml:"counters"`
	Status          HealthStatus                    `yaml:"status" json:"status" xml:"status"`
	Hostname        string                          `yaml:"hostname" json:"hostname" xml:"hostname"`
	IndicatorLED    string                          `yaml:"indicator_led" json:"indicator_led" xml:"indicator_led"`
	PowerState      string                          `yaml:"power_state" json:"power_state" xml:"power_state"`
	ProcessorModel  string                          `yaml:"processor_model" json:"processor_model" xml:"processor_model"`
	ProcessorStatus HealthStatus                    `yaml:"processor_status" json:"processor_status" xml:"processor_status"`
	MemoryMirroring string                          `yaml:"memory_mirroring" json:"memory_mirroring" xml:"memory_mirroring"`
	MemoryStatus    HealthStatus                    `yaml:"memory_status" json:"memory_status" xml:"memory_status"`
	ActionEndpoints []*ComputerSystemActionEndpoint `yaml:"action_endpoints" json:"action_endpoints" xml:"action_endpoints"`
}

// ComputerSystemActionEndpoint represents write-capable API endpoint.
type ComputerSystemActionEndpoint struct {
	Action        string   `yaml:"action" json:"action" xml:"action"`
	Target        string   `yaml:"target" json:"target" xml:"target"`
	AllowedValues []string `yaml:"allowed_actions" json:"allowed_actions" xml:"allowed_actions"`
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
	cs.Counters.BootOrder = response.Boot.BootOrderCounter
	cs.Counters.Chassis = response.Links.ChassisCounter
	cs.Counters.CooledBy = response.Links.CooledByCounter
	cs.Counters.ManagedBy = response.Links.ManagedByCounter
	cs.Counters.PoweredBy = response.Links.PoweredByCounter
	cs.Counters.PCIeDevices = response.PCIeDevicesCounter
	cs.Counters.PCIeFunctions = response.PCIeFunctionsCounter
	cs.Counters.HostingRoles = response.HostingRolesCounter

	cs.Counters.TotalProcessors = response.ProcessorSummary.Count
	cs.Counters.LogicalProcessors = response.ProcessorSummary.LogicalProcessorCount
	switch vt := response.MemorySummary.TotalSystemMemoryGiB.(type) {
	case uint64:
		cs.Counters.TotalSystemMemory = response.MemorySummary.TotalSystemMemoryGiB.(uint64)
	case float64:
		cs.Counters.TotalSystemMemory = uint64(response.MemorySummary.TotalSystemMemoryGiB.(float64))
	default:
		return nil, fmt.Errorf("parsing error: %s, server response: %s, unsupported type: %T", err, string(s[:]), vt)
	}

	cs.Status = response.Status

	cs.Hostname = response.HostName
	cs.IndicatorLED = response.IndicatorLED
	cs.PowerState = response.PowerState

	cs.ProcessorModel = response.ProcessorSummary.Model
	cs.ProcessorStatus = response.ProcessorSummary.Status

	cs.MemoryMirroring = response.MemorySummary.MemoryMirroring
	cs.MemoryStatus = response.MemorySummary.Status
	cs.ActionEndpoints = []*ComputerSystemActionEndpoint{}

	if response.Actions != nil {
		for k, v := range response.Actions {
			endpoint := &ComputerSystemActionEndpoint{
				Action:        k,
				Target:        v.Target,
				AllowedValues: v.AllowedValues,
			}
			cs.ActionEndpoints = append(cs.ActionEndpoints, endpoint)
		}
	}

	return cs, nil
}
