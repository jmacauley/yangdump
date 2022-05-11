# YANG Schema Dump Tool (yangdump)

This simple golang tool will download the current version of YANG or YIN schema files from the target device using the NETCONF protocol.

## Usage

yangdump -h *\<host\>* -u *\<userid\>* -p *\<password\>* -d *\<target directory\>* -yin

```
$ yangdump 
Usage of yangdump:
  -d string
    	Schema destination directory. (default ".")
  -h string
    	Hostname
  -p string
    	Password
  -u string
    	Username
  -k string
        SSH private key file.
  -yin
    	Download YIN format instead of YANG.
```

## Building

A simple Makefile is included to build `yangdump`, otherwise, `go build` will get you there.

## How does it work?

Using the `ietf-netconf-monitoring` schema as defined in RFC 6022 we issue a NETCONF `get`	 request to a device for a list of supported schema.

```
<get>
	<filter type="subtree">
		<netconf-state xmlns="urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring">
			<schemas/>
		</netconf-state>
	</filter>
</get>
```

The device will return a list of schema available for download. 

```
  <data>
    <netconf-state xmlns="urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring">
      <schemas>
        <schema>
          <identifier>openroadm-augment</identifier>
          <version>2020-01-03</version>
          <format>yin</format>
          <namespace>http://coriant.com/yang/os/openroadm-augment</namespace>
          <location>NETCONF</location>
        </schema>
        <schema>
          <identifier>openroadm-augment</identifier>
          <version>2020-01-03</version>
          <format>yang</format>
          <namespace>http://coriant.com/yang/os/openroadm-augment</namespace>
          <location>NETCONF</location>
        </schema>
		...
      </schemas>
    </netconf-state>
  </data>
```

We then use `get-schema` operations defined in RFC 6022 to get the individual schema files listed on the device of format `yang` and location `NETCONF`.

```
<get-schema xmlns="urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring">
	<identifier>openroadm-augment</identifier>
	<version>2020-01-03</version>
	<format>yang</format>
</get-schema>
```

If the `-yin` option is specified `yangdump` will download the XML YIN schema representation.

## Notes
Structures for parsing the `ietf-netconf-monitoring` schema have been hand created using `xml:` annotations and not using a compiler such as that provide by tools such as `ygot`.  The reason for this is that there are no golang tools for building XML parsers for NETCONF/XML, however is good support for NESTCONF/JSON.  Unfortunately, Infinera only supports NETCONF/XML so XML parsers are required.  Infinera does support RESTCONF/JSON, but does not support `ietf-netconf-monitoring` over RESTCONF.
