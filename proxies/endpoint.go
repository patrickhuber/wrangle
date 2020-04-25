package proxies

import (
	"fmt"
	"strconv"
	"strings"
)

// Endpoint represents an ssh endpoint
type Endpoint interface {
	Host() string
	Port() int
	User() string
	fmt.Stringer
}

type endpoint struct {
	host string
	port int
	user string
}

// NewEndpoint parses the ssh string into a new endpoint struct
func NewEndpointFromString(s string) Endpoint {
	endpoint := &endpoint{
		host: s,
	}

	if parts := strings.Split(endpoint.host, "@"); len(parts) > 1 {
		endpoint.user = parts[0]
		endpoint.host = parts[1]
	}
	if parts := strings.Split(endpoint.host, ":"); len(parts) > 1 {
		endpoint.host = parts[0]
		endpoint.port, _ = strconv.Atoi(parts[1])
	}
	return endpoint
}

// NewEndpoint uses structured inputs to create an endpoint
func NewEndpoint(host string, port int, user string) Endpoint {
	return &endpoint{
		host: host,
		port: port,
		user: user,
	}
}

func (e *endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.host, e.port)
}

func (e *endpoint) Host() string {
	return e.host
}

func (e *endpoint) User() string {
	return e.user
}

func (e *endpoint) Port() int {
	return e.port
}
