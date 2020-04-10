package client

// ODataAnnotation represent OData Annotations. See Instance Annotations
// in OData JSON Format Version 4.01 Specification.
type ODataAnnotation struct {
	ID      string `yaml:"@odata.id" json:"@odata.id" xml:"@odata.id"`
	Type    string `yaml:"@odata.type" json:"@odata.type" xml:"@odata.type"`
	Context string `yaml:"@odata.context" json:"@odata.context" xml:"@odata.context"`
}

// NewODataAnnotation returns an instance of ODataAnnotation
func NewODataAnnotation(odaID, odaType, odaContext string) *ODataAnnotation {
	return &ODataAnnotation{
		ID:      odaID,
		Type:    odaType,
		Context: odaContext,
	}
}
