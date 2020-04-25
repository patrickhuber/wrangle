package proxies

import (
	"context"
	"io"
	"net"

	"github.com/patrickhuber/wrangle/logging"
	xssh "golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

type ssh struct {
	cfg    *SSHConfig
	cancel bool
	logger logging.Logger
}

// SSHConfig creates a ssh configuration
type SSHConfig struct {
	Local        Endpoint
	Server       Endpoint
	Remote       Endpoint
	ClientConfig *xssh.ClientConfig
}

// NewSSHConfig returns a new ssh config for the given inputs
func NewSSHConfig(local Endpoint, server Endpoint, remote Endpoint, auth xssh.AuthMethod) *SSHConfig {
	return &SSHConfig{
		Local:  local,
		Server: server,
		Remote: remote,
		ClientConfig: &xssh.ClientConfig{
			Auth: []xssh.AuthMethod{
				auth,
			},
			HostKeyCallback: xssh.InsecureIgnoreHostKey(),
			User:            server.User(),
		},
	}
}

// NewSSH creates a new SSH Proxy
func NewSSH(cfg *SSHConfig, logger logging.Logger) Proxy {
	return &ssh{
		cfg:    cfg,
		logger: logger,
		cancel: false,
	}
}

func (s *ssh) Serve(ctx context.Context) error {

	// establish a connection with the ssh server
	server, err := xssh.Dial("tcp", s.cfg.Server.String(), s.cfg.ClientConfig)
	if err != nil {
		s.logError(err)
		return err
	}
	defer server.Close()

	// establish a connection with the remote server
	remote, err := server.Dial("tcp", s.cfg.Remote.String())
	if err != nil {
		s.logError(err)
		return err
	}
	defer remote.Close()

	// start a local server to forward traffic to the remote location
	listenConfig := &net.ListenConfig{}
	local, err := listenConfig.Listen(ctx, "tcp", s.cfg.Local.String())
	if err != nil {
		s.logError(err)
		return err
	}
	defer local.Close()

	// handle incomming connections
	for {
		// if cancel signal reached, terminate the loop
		if s.cancel {
			return nil
		}

		g, ctx := errgroup.WithContext(ctx)
		clientChan := make(chan net.Conn)
		g.Go(func() error {
			client, err := local.Accept()
			if err != nil {
				return err
			}
			clientChan <- client
			close(clientChan)
			return nil
		})

		copy := func(source net.Conn, destination net.Conn) error {
			_, err := io.Copy(source, destination)
			return err
		}

		select {
		// if we have a client, use the client
		case client := <-clientChan:
			g.Go(func() error {
				return copy(client, remote)
			})
			g.Go(func() error {
				return copy(remote, client)
			})
			// block until both routines are done, return any error
			err = g.Wait()
			if err != nil {
				return err
			}
		// the context was canceled so return the error
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *ssh) logError(err error) {
	if s.logger == nil {
		return
	}
	s.logger.Fatalln(err)
}

func (s *ssh) Close() error {
	s.cancel = true
	return nil
}
