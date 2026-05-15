// 2026 Jan Provaznik (jan@provaznik.pro)
//
// Inspired by the original tftp example.

package main

import "io"
import "os"
import "fmt"

import "github.com/pin/tftp/v3"

func main () {
	if 3 != len(os.Args) {
		fmt.Println("dhcpee [port] [root]")
		fmt.Println("... serves root over tftp listening on port")
		os.Exit(1)
	}

	port := os.Args[1]
	path := os.Args[2]

	// Use traversal-resistant paths
	// https://go.dev/blog/osroot
	root, err := os.OpenRoot(path)
	if err != nil {
		fmt.Println("Could not open the root path.")
		os.Exit(1)
	}
	defer root.Close()

	reader := func (path string, rf io.ReaderFrom) error {
		addr := rf.(tftp.OutgoingTransfer).RemoteAddr()
		fmt.Println("addr", addr.String(), "pull", path)

		file, err := root.Open(path)
		if err != nil {
			return err
		}
		_, err = rf.ReadFrom(file)
		if err != nil {
			return err
		}
		return nil
	}

	server := tftp.NewServer(reader, nil)

	// The beacons are lit
	err = server.ListenAndServe(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Could not start the server.")
		os.Exit(1)
	}

}

