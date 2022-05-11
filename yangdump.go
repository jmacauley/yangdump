// yangdump will download the current version of YANG or YIN schema files
// from the target device using the NETCONF protocol.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Permissions for directory and files we create.
const FilePermission os.FileMode = 0644
const DirectoryPermission os.FileMode = 0755

// Command line parameters.
var (
	host         = flag.String("h", "", "Hostname")
	username     = flag.String("u", "", "Username")
	password     = flag.String("p", "", "Password")
	keyfile      = flag.String("k", "", "SSH private key file")
	dest         = flag.String("d", ".", "Schema destination directory.")
	version      = flag.Bool("v", false, "Create YANG file with version in filename.")
	yin          = flag.Bool("yin", false, "Download YIN format instead of YANG.")
)

// check - If error then exit.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// main - all the magic gets done in here.
func main() {
	// Get the command line options.
	flag.Parse()

	// Make sure user provided required parameters.
	if *host == "" || *username == "" || *password == "" {
		flag.Usage()
		os.Exit(1)
	}

	var format string = YANG
	if *yin {
		format = YIN
	}

	var vers bool = false
	if *version {
		vers = true
	}

	// Create directory if it doesn't exist.
	var err error
	if _, err = os.Stat(*dest); os.IsNotExist(err) {
		os.Mkdir(*dest, DirectoryPermission)
	}

	// Get a valid NETCONF session to the device.
	var session *Session
	if keyfile != nil && *keyfile != "" {
		fmt.Println("Sorry I am lazy and have not implemented this option just yet...")
		os.Exit(1)
	} else {
		session, err = newSession(*host, *username, *password)
	}
	check(err)

	// Defer session close until we are done with our NETCONF interactions.
	defer session.Close()

	// Get the list of schema from the device.
	var schemaList []SchemaType
	schemaList, err = session.getSchemaList(format)
	check(err)

	// For each schema returned in the list we want to retrieve the instance document from the device.
	for i := 0; i < len(schemaList); i++ {
		// Give them some indication of progress without without writing filenames.
		fmt.Printf("\r%s\rDownloading... %d of %d complete",
			strings.Repeat(" ", 35), i + 1, len(schemaList))

		// Download the YANG schema from device.
		var schema = schemaList[i]
		var source *SourceType
		source, err = session.getSchema(schema)
		check(err)

		// Write the file to specified location.
		var filename string
		if schema.Version == "" || !vers {
			filename = fmt.Sprintf("%s/%s.%s", *dest, schema.Identifier, schema.Format)
		} else {
			filename = fmt.Sprintf("%s/%s@%s.%s", *dest, schema.Identifier, schema.Version, schema.Format)
		}

		err = ioutil.WriteFile(filename, []byte(source.Data), FilePermission)
		check(err)
	}

	fmt.Printf("\nDownload to %s finished\n", *dest)
}
