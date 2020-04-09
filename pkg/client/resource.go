// Copyright 2020 Paul Greenberg (greenpau@outlook.com)

package client

// Resource represents raw Redfish object
type Resource struct {
	Raw []byte
}

func (r *Resource) String() string {
	return string(r.Raw)
}

// NewResourceFromString returns Resource instance from an input string.
func NewResourceFromString(s string) (*Resource, error) {
	return NewResourceFromBytes([]byte(s))
}

// NewResourceFromBytes returns Resource instance from an input byte array.
func NewResourceFromBytes(s []byte) (*Resource, error) {
	r := &Resource{}
	r.Raw = s
	return r, nil
}
