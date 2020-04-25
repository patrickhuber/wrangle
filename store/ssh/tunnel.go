package ssh

import (
	"io"
	"net"

	xssh "golang.org/x/crypto/ssh"
)

type tunnel struct {
	local  Endpoint
	server Endpoint
	remote Endpoint
	config *xssh.ClientConfig
}

// Tunnel defines a ssh tunnel
type Tunnel interface {
	Start() error
}

// NewTunnel creates a new tunnel
func NewTunnel(local Endpoint, server Endpoint, remote Endpoint, auth xssh.AuthMethod) Tunnel {
	return &tunnel{
		local:  local,
		server: server,
		remote: remote,
		config: &xssh.ClientConfig{
			User: server.User(),
			Auth: []xssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key xssh.PublicKey) error {
				// alway accept
				return nil
			},
		},
	}
}

func (t *tunnel) Start() error {
	listener, err := net.Listen("tcp", t.local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go t.forward(conn)
	}
}

func (t *tunnel) forward(localConn net.Conn) {
	serverConn, err := xssh.Dial("tcp", t.server.String(), t.config)
	if err != nil {
		return
	}

	remoteConn, err := serverConn.Dial("tcp", t.remote.String())
	if err != nil {
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			return
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
