package tailers

import (
	"io"
)

// SSHCommand is used to manage individual commands sent to the SSH client.
type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}