package main

import (
	"net"
	"strings"
)

// Get - return free open TCP port
func getAvailablePort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return 0, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}

// Slugy takes arry of string as input and makes
// slug for it by replaceing any blank space to dash
func Slugy(inputs []string) string {
	return strings.TrimSpace(strings.ToLower(strings.Join(inputs, "-")))
}
