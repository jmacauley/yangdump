// session contains functionality needed to download YANG/YIN
// schema as defined in ietf-netconf-monitoring.
package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/Juniper/go-netconf/netconf"
)

const YIN string = "yin"
const YANG string = "yang"

type Session struct {
	handle *netconf.Session
}

// Get a NETCONF session to the target device.
func newSession(host string, username string, password string) (session *Session, err error) {

	// Connect to target device.
	var s *netconf.Session
	s, err = netconf.DialSSH(host, netconf.SSHConfigPassword(username, password))
	if err != nil {
		return nil, err
	}

	session = &Session{}
	session.handle = s
	return session, nil
}

// Close is used to close and end a transport session
func (s *Session) Close() error {
	return s.handle.Close()
}

// Retrieve list of available schema from device.
func (s *Session) getSchemaList(format string) (result []SchemaType, err error) {
	var buffer *DataType = &DataType{}

	// Get the schema list
	var reply *netconf.RPCReply
	reply, err = s.handle.Exec(netconf.RawMethod(GetSchemaList))
	if err != nil {
		fmt.Printf("%v", reply.Errors)
		return nil, err
	}

	fmt.Printf("%v", reply.Data)

	// Unmarshal the XML into our had built type structures.
	err = xml.Unmarshal([]byte(reply.Data), buffer)
	if err != nil {
		fmt.Printf("%v", reply.Errors)
		return nil, err
	}

	// For each schema returned in the list we want to retrieve the YANG instance from the device.
	var in []SchemaType = buffer.NetconfState.Schemas.SchemaList
	for i := 0; i < len(in); i++ {
		var list []string = strings.Split(in[i].Format, ":")
		if len(list) > 1 {
			in[i].Format = list[1]
		} else {
			in[i].Format = list[0]
		}
		if  strings.Contains(in[i].Format, format) &&
			strings.Compare(in[i].Location, "NETCONF") == 0 {
			result = append(result, in[i])
		}
	}

	return result, nil
}

// Retrieve an individual schema instance from device.
func (s *Session) getSchema(schema SchemaType) (source *SourceType, err error) {
	// Build the NETCONF get message for the schema document.
	var request string
	request = fmt.Sprintf(GetSchemaInstance, schema.Identifier, schema.Version, schema.Format)

	var reply *netconf.RPCReply
	reply, err = s.handle.Exec(netconf.RawMethod(request))
	if err != nil {
		fmt.Printf("%v", reply.Errors)
		return nil, err
	}

	// Unmarshal the document.
	source = &SourceType{}
	err = xml.Unmarshal([]byte(reply.Data), source)
	if err != nil {
		return nil, err
	}

	// YIN is an XML document so our simple SourceType "chardata" will not work so manually parse.
	if schema.Format == YIN {
		s1 := strings.Index(reply.Data, "<module")
		s2 := strings.Index(reply.Data[s1:], "</data>")
		source.Data = xml.Header + reply.Data[s1:s1 + s2]
	}

	return source, nil
}

