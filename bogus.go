package bogus

import (
	"context"
	"net"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/caddyserver/caddy"
	"github.com/miekg/dns"
)

// N implements the plugin interface.
type N struct {
	Next plugin.Handler
	ips  []net.IP
}

func init() { plugin.Register("bogus", setup) }

func setup(c *caddy.Controller) error {
	ips := []net.IP{}
	for c.Next() {
		args := c.RemainingArgs()

		for _, a := range args {
			ips = append(ips, net.ParseIP(a))
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return N{Next: next, ips: ips}
	})

	c.OnStartup(func() error {
		return nil
	})
	return nil
}

// ServeDNS implements the plugin.Handler interface.
func (n N) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wr := NewResponseReverter(w, r, n.ips)

	if len(n.ips) == 0 {
		return plugin.NextOrFailure(n.Name(), n.Next, ctx, w, r)
	}
	return plugin.NextOrFailure(n.Name(), n.Next, ctx, wr, r)
}

// Name implements the Handler interface.
func (n N) Name() string { return "bogus" }
