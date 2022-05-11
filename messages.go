// messages are the NETCONF protocol messages sent to a device to
// retrieve YANG/YIN schema as defined in ietf-netconf-monitoring.
package main

// RFC 6022 - NETCONF message to get the list of supported schema.
const GetSchemaList =
	"<get>" +
		"<filter type=\"subtree\">" +
		"<netconf-state xmlns=\"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring\">" +
		"<schemas/>" +
		"</netconf-state>" +
		"</filter>" +
		"</get>"

// RFC 6022 - NETCONF message to get an individual schema file.
const GetSchemaInstance =
	"<get-schema xmlns=\"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring\">" +
		"<identifier>%s</identifier>" +
		"<version>%s</version>" +
		"<format>%s</format>" +
		"</get-schema>"
