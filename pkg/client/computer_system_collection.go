// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
)

type computerSystemCollectionResponse struct {
	ODataAnnotation
	Name         string
	Description  string
	MembersCount uint64 `yaml:"Members@odata.count" json:"Members@odata.count" xml:"Members@odata.count"`
	Members      []ODataAnnotation
}

type computerSystemCollectionCounters struct {
	ComputerSystems uint64 `yaml:"computer_systems" json:"computer_systems" xml:"computer_systems"`
}

// ComputerSystemCollection represents an instance of Redfish ComputerSystemCollection.
type ComputerSystemCollection struct {
	OData           *ODataAnnotation                 `yaml:"odata" json:"odata" xml:"odata"`
	Name            string                           `yaml:"name" json:"name" xml:"name"`
	Description     string                           `yaml:"description" json:"description" xml:"description"`
	Counters        computerSystemCollectionCounters `yaml:"counters" json:"counters" xml:"counters"`
	Members         []*ODataAnnotation               `yaml:"members" json:"members" xml:"members"`
	ComputerSystems []*ComputerSystem                `yaml:"computer_systems" json:"computer_systems" xml:"computer_systems"`
}

// GetComputerSystemCollection returns an instance of Redfish ComputerSystemCollection.
func (cli *Client) GetComputerSystemCollection() (*ComputerSystemCollection, error) {
	resp, err := cli.callAPI("GET", "", cli.rootPath+"Systems/", []byte{})
	if err != nil {
		return nil, err
	}
	csc, err := newComputerSystemCollectionFromBytes(resp)
	return csc, err
}

// GetComputerSystems returns ComputerSystem instances
func (cli *Client) GetComputerSystems() ([]*ComputerSystem, error) {
	csc, err := cli.GetComputerSystemCollection()
	if err != nil {
		return nil, err
	}
	if err := csc.getComputerSystems(cli); err != nil {
		return nil, err
	}
	return csc.ComputerSystems, nil
}

func (csc *ComputerSystemCollection) getComputerSystems(cli *Client) error {
	if len(csc.ComputerSystems) > 0 {
		return nil
	}
	computerSystems := []*ComputerSystem{}
	for _, member := range csc.Members {
		cs, err := cli.GetComputerSystemByResourceID(member.ID + "/")
		if err != nil {
			return err
		}
		computerSystems = append(computerSystems, cs)
	}
	csc.ComputerSystems = computerSystems
	return nil
}

// newComputerSystemCollectionFromString returns ComputerSystemCollection instance from an input string.
func newComputerSystemCollectionFromString(s string) (*ComputerSystemCollection, error) {
	return newComputerSystemCollectionFromBytes([]byte(s))
}

// newComputerSystemCollectionFromBytes returns ComputerSystemCollection instance from an input byte array.
func newComputerSystemCollectionFromBytes(s []byte) (*ComputerSystemCollection, error) {
	csc := &ComputerSystemCollection{
		ComputerSystems: []*ComputerSystem{},
	}
	cscResponse := &computerSystemCollectionResponse{}
	err := json.Unmarshal(s, cscResponse)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s, server response: %s", err, string(s[:]))
	}
	csc.OData = &ODataAnnotation{
		Context: cscResponse.Context,
		ID:      cscResponse.ID,
		Type:    cscResponse.Type,
	}
	csc.Name = cscResponse.Name
	csc.Description = cscResponse.Description
	csc.Counters.ComputerSystems = cscResponse.MembersCount
	for _, member := range cscResponse.Members {
		csc.Members = append(csc.Members, &member)
	}
	return csc, nil
}
