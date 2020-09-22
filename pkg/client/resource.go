// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

import (
	"encoding/json"
	"fmt"
)

type resourceResponse struct {
	ODataAnnotation
}

// Resource represents Redfish resource object
type Resource struct {
	OData *ODataAnnotation `yaml:"odata_annotations" json:"odata_annotations" xml:"odata_annotations"`
	Raw   []byte           `yaml:"-" json:"-" xml:"-"`
}

func (r *Resource) String() string {
	return string(r.Raw)
}

// GetResource return raw output from Redfish API server for a specific resource (path)
func (cli *Client) GetResource(s string) (*Resource, error) {
	resp, err := cli.callAPI("GET", "", s, []byte{})
	if err != nil {
		return nil, err
	}
	return newResourceFromBytes(resp)
}

// newResourceFromString returns Resource instance from an input string.
func newResourceFromString(s string) (*Resource, error) {
	return newResourceFromBytes([]byte(s))
}

// newResourceFromBytes returns Resource instance from an input byte array.
func newResourceFromBytes(s []byte) (*Resource, error) {
	r := &Resource{
		Raw: s,
	}
	response := &infoResponse{}
	err := json.Unmarshal(s, response)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %s, server response: %s", err, string(s[:]))
	}

	r.OData = &ODataAnnotation{
		Context: response.Context,
		ID:      response.ID,
		Type:    response.Type,
	}

	return r, nil
}
