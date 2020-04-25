package proxies

import "context"

// Proxy defines a proxy for ssh or http
type Proxy interface {
	Serve(ctx context.Context) error
}
