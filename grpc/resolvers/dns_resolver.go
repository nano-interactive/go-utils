package resolvers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	grpclbstate "google.golang.org/grpc/balancer/grpclb/state"
	"google.golang.org/grpc/resolver"
)

type (
	DNSResolver struct {
		cc         resolver.ClientConn
		resolver   NetResolver
		resolveNow chan struct{}
		cancel     context.CancelFunc
		ips        []resolver.Address
		wg         sync.WaitGroup
		mu         sync.RWMutex
	}

	DNSResolverBuilder struct {
		ctx     context.Context
		options *option
	}

	NetResolver interface {
		LookupHost(context.Context, string) ([]string, error)
		LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
		LookupSRV(context.Context, string, string, string) (string, []*net.SRV, error)
	}

	option struct {
		resolver         NetResolver
		dnsAuthority     string
		dnsReResolving   time.Duration
		resolvingTimeout time.Duration
		defaultPort      uint16
	}

	Option func(o *option)
)

func WithResolver(r NetResolver) Option {
	return func(o *option) {
		o.resolver = r
	}
}

func WithDNSReResolving(d time.Duration) Option {
	return func(o *option) {
		o.dnsReResolving = d
	}
}

func WithResolvingTimeout(d time.Duration) Option {
	return func(o *option) {
		o.resolvingTimeout = d
	}
}

func WithDnsAuthority(authority string) Option {
	return func(o *option) {
		o.dnsAuthority = authority
	}
}

func WithDefaultPort(p uint16) Option {
	return func(o *option) {
		o.defaultPort = p
	}
}

func NewDNSResolverBuilder(ctx context.Context, opts ...Option) *DNSResolverBuilder {
	cfg := &option{
		dnsReResolving:   30 * time.Second,
		defaultPort:      443,
		resolvingTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &DNSResolverBuilder{
		ctx:     ctx,
		options: cfg,
	}
}

func (d *DNSResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(d.ctx)
	endpoint := target.Endpoint()

	host, port, err := parseTarget(endpoint, strconv.FormatInt(int64(d.options.defaultPort), 10))
	if err != nil {
		cancel()
		return nil, err
	}

	if ipAddr, ok := formatIP(host); ok {
		addr := []resolver.Address{{Addr: ipAddr + ":" + port}}
		_ = cc.UpdateState(resolver.State{Addresses: addr})
		cancel()
		return resolver.Get("passthrough").Build(target, cc, opts)
	}

	if d.options.resolver == nil {
		d.options.resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				authority := address

				if d.options.dnsAuthority != "" {
					authority = d.options.dnsAuthority
				}

				dialer := net.Dialer{
					Timeout:   d.options.resolvingTimeout,
					KeepAlive: 1 * time.Minute,
				}

				return dialer.DialContext(ctx, network, authority)
			},
		}
	}

	resolver := &DNSResolver{
		ips:        make([]resolver.Address, 0),
		resolveNow: make(chan struct{}, 1),
		cancel:     cancel,
		cc:         cc,
		resolver:   d.options.resolver,
	}

	if err := resolver.resolve(ctx, host, port, d.options.resolvingTimeout); err != nil {
		return nil, err
	}

	resolver.wg.Add(1)
	go func() {
		defer resolver.wg.Done()

		go resolver.watcher(ctx, host, port, d.options.dnsReResolving, d.options.resolvingTimeout)
	}()

	return resolver, nil
}

func (d *DNSResolverBuilder) Scheme() string {
	return "nanodns"
}

func (d *DNSResolver) resolve(ctx context.Context, host, port string, timeout time.Duration) error {
	ips, err := resolve(ctx, d.resolver, host, port, timeout)
	if err != nil {
		if errors.Is(err, errNoChange) {
			return nil
		}
		d.cc.ReportError(err)
		return err
	}

	slices.SortStableFunc(ips.Addresses, func(a, b resolver.Address) int {
		return strings.Compare(a.Addr, b.Addr)
	})

	d.mu.RLock()
	{
		eq := slices.EqualFunc(d.ips, ips.Addresses, func(a, b resolver.Address) bool {
			return a.Addr == b.Addr
		})

		if eq {
			d.mu.RUnlock()
			return nil
		}
	}
	d.mu.RUnlock()

	d.mu.Lock()
	{
		d.ips = d.ips[:0]
		d.ips = append(d.ips, ips.Addresses...)
	}
	d.mu.Unlock()

	if err := d.cc.UpdateState(ips); err != nil {
		d.cc.ReportError(err)
	}

	return nil
}

func (d *DNSResolver) watcher(ctx context.Context, host, port string, refresh, timeout time.Duration) {
	ticker := time.NewTicker(refresh)
	reResolveTimeout := time.NewTicker(refresh)

	defer func() {
		ticker.Stop()
		reResolveTimeout.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = d.resolve(ctx, host, port, timeout)
		case <-d.resolveNow:
			select {
			case <-ctx.Done():
				return
			case <-reResolveTimeout.C:
				// Prevents spam calls to the DNS Resolver
				// has to wait for at least 30 seconds
				// between each call
				_ = d.resolve(ctx, host, port, timeout)
			}
		}
	}
}

func resolve(ctx context.Context, r NetResolver, host, port string, timeout time.Duration) (resolver.State, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	addrs, err := r.LookupHost(ctx, host)
	if err != nil {
		return resolver.State{}, err
	}

	if len(addrs) == 0 {
		return resolver.State{}, nil
	}

	addr := make([]resolver.Address, len(addrs))

	for i, a := range addrs {
		addr[i] = resolver.Address{Addr: net.JoinHostPort(a, port)}
	}

	state := resolver.State{Addresses: addr}

	_, srvs, err := r.LookupSRV(ctx, "grpclb", "tcp", host)

	if err == nil && len(srvs) > 0 {
		var newAddrs []resolver.Address

		for _, s := range srvs {
			lbAddrs, err := r.LookupHost(ctx, s.Target)
			if err != nil {
				continue
			}

			for _, a := range lbAddrs {
				ip, ok := formatIP(a)
				if ok {
					addr := ip + ":" + strconv.Itoa(int(s.Port))
					newAddrs = append(newAddrs, resolver.Address{Addr: addr, ServerName: s.Target})
				}
			}
		}

		state = grpclbstate.Set(state, &grpclbstate.State{
			BalancerAddresses: newAddrs,
		})
	}

	return state, nil
}

func (d *DNSResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case d.resolveNow <- struct{}{}:
	default:
	}
}

func (d *DNSResolver) Close() {
	d.cancel()
	d.wg.Wait()
	close(d.resolveNow)
}

func formatIP(addr string) (string, bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}

var (
	ErrEndsWithColon = errors.New("")
	errNoChange      = errors.New("no ip changes")
)

func parseTarget(target, defaultPort string) (string, string, error) {
	if ip := net.ParseIP(target); ip != nil {
		return target, defaultPort, nil
	}
	if host, port, err := net.SplitHostPort(target); err == nil {
		if port == "" {
			return "", "", ErrEndsWithColon
		}
		if host == "" {
			host = "localhost"
		}
		return host, port, nil
	}

	host, port, err := net.SplitHostPort(target + ":" + defaultPort)
	if err != nil {
		return host, port, fmt.Errorf("invalid target address %v, error info: %w", target, err)
	}

	return host, port, nil
}
