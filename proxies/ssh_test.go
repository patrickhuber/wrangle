package proxies_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gliderlabs/ssh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/logging"
	"github.com/patrickhuber/wrangle/proxies"
	xssh "golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

var (
	DeadlineTimeout = 30 * time.Second
	IdleTimeout     = 10 * time.Second
	group           errgroup.Group
	username        = "admin"
	password        = "p@$$w0rd"
)

var _ = BeforeSuite(func() {

	forwardHandler := &ssh.ForwardedTCPHandler{}

	server := &ssh.Server{
		Addr:        ":2222",
		MaxTimeout:  DeadlineTimeout,
		IdleTimeout: IdleTimeout,
		LocalPortForwardingCallback: ssh.LocalPortForwardingCallback(func(ctx ssh.Context, dhost string, dport uint32) bool {
			log.Println("Accepted forward", dhost, dport)
			return true
		}),
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		RequestHandlers: map[string]ssh.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"session":      ssh.DefaultSessionHandler,
			"direct-tcpip": ssh.DirectTCPIPHandler,
		},
		Handler: ssh.Handler(func(s ssh.Session) {
			log.Println("new connection")
			i := 0
			for {
				i += 1
				log.Println("active seconds:", i)
				select {
				case <-time.After(time.Second):
					continue
				case <-s.Context().Done():
					log.Println("connection closed")
					return
				}
			}
		}),
		PasswordHandler: func(ctx ssh.Context, pass string) bool {
			return ctx.User() == username && pass == password
		},
	}
	group.Go(server.ListenAndServe)
})

var _ = Describe("Ssh", func() {
	Describe("Serve", func() {
		It("can connect", func() {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			g, ctx := errgroup.WithContext(ctx)

			mux := http.NewServeMux()
			httpHandler := func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "hello world")
			}

			mux.HandleFunc("/", httpHandler)

			listenConfig := net.ListenConfig{}
			listener, err := listenConfig.Listen(ctx, "tcp", ":5555")
			Expect(err).To(BeNil())
			g.Go(func() error {
				return http.Serve(listener, mux)
			})

			local := proxies.NewEndpointFromString(":5000")
			server := proxies.NewEndpointFromString(username + "@localhost:2222")
			remote := proxies.NewEndpointFromString("localhost:5555")
			auth := xssh.Password(password)
			c := proxies.NewSSHConfig(local, server, remote, auth)
			s := proxies.NewSSH(c, logging.Default())

			g.Go(func() error { return s.Serve(ctx) })
			err = g.Wait()
			Expect(err).To(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	group.Wait()
})
