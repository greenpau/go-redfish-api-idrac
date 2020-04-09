package client

type apiPath struct {
	Path    string `yaml:"@odata.id" json:"@odata.id" xml:"@odata.id"`
	Type    string `yaml:"@odata.type" json:"@odata.type" xml:"@odata.type"`
	Context string `yaml:"@odata.context" json:"@odata.context" xml:"@odata.context"`
}
