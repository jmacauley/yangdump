// schema contains a manual XML mapping for the ietf-netconf-monitoring
// YANG schema needed to download YANG/YIN schema stored on a device.
package main

import "encoding/xml"

type DataType struct {
	XMLName      xml.Name `xml:"data"`
	NetconfState NetconfStateType `xml:"netconf-state"`
}

type NetconfStateType struct {
	Name        string   `xml:"xmlns,attr"`
	Schemas 	SchemasType `xml:"schemas"`
}

type SchemasType struct {
	SchemaList []SchemaType `xml:"schema"`
}

type SchemaType struct {
	Identifier 	string `xml:"identifier"`
	Version 	string `xml:"version"`
	Format 		string `xml:"format"`
	Namespace   string `xml:"namespace"`
	Location    string `xml:"location"`
}

type SourceType struct {
	XMLName    xml.Name `xml:"data"`
	Name	   string `xml:"xmlns,attr"`
	Data       string `xml:",chardata"`
}
