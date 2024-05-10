package resolvers_test

import (
	"context"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/nano-interactive/go-utils/v2/grpc/resolvers"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc/serviceconfig"
)

type (
	MockClientConn struct {
		mock.Mock
	}
	MockNetResolver struct {
		mock.Mock
	}
)

func (m *MockClientConn) UpdateState(p0 resolver.State) error {
	args := m.Called(p0)
	return args.Error(0)
}

func (m *MockClientConn) ReportError(p0 error) {
	_ = m.Called(p0)
}

func (m *MockClientConn) NewAddress(addresses []resolver.Address) {
	panic("SHOULD PANIC: NOT EXPECTED ANYWHERE")
}

func (m *MockClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	panic("SHOULD PANIC: NOT EXPECTED ANYWHERE")
}

func (m *MockNetResolver) LookupHost(p0 context.Context, p1 string) ([]string, error) {
	args := m.Called(p1)

	return args.Get(0).([]string), args.Error(1)
}

func (m *MockNetResolver) LookupIPAddr(p0 context.Context, p1 string) ([]net.IPAddr, error) {
	args := m.Called(p1)

	return args.Get(0).([]net.IPAddr), args.Error(1)
}

func (m *MockNetResolver) LookupSRV(p0 context.Context, p1 string, p2 string, p3 string) (string, []*net.SRV, error) {
	args := m.Called(p1, p2, p3)

	return args.String(0), args.Get(1).([]*net.SRV), args.Error(2)
}

func getDNSBuilder(resolver resolvers.NetResolver, resolving time.Duration) *resolvers.DnsResolverBuilder {
	return resolvers.NewDnsResolverBuilder(
		context.Background(),
		resolvers.WithDnsReResolving(resolving),
		resolvers.WithResolvingTimeout(3*time.Second),
		resolvers.WithDefaultPort(3000),
		resolvers.WithResolver(resolver),
	)
}

func TestDNSResolver(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	conn := &MockClientConn{}
	builder := getDNSBuilder(nil, 1*time.Hour)

	conn.
		On("UpdateState", mock.Anything).
		Once().
		Return(nil)

	r, err := builder.Build(
		resolver.Target{
			URL: *lo.Must(url.Parse("dusanmalusev.dev")),
		},
		conn,
		resolver.BuildOptions{},
	)

	// As it domain has already been resolved, UpdateState should be called only Once
	// on .Build(), and never again on .ResolveNow()
	r.ResolveNow(resolver.ResolveNowOptions{})
	r.ResolveNow(resolver.ResolveNowOptions{})

	assert.Equal("nanodns", builder.Scheme())
	assert.NoError(err)
	assert.NotNil(r)
	conn.AssertExpectations(t)

	r.Close() // Stop the watcher
}

func TestDNSResolver_TriggerIPChange(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	conn := &MockClientConn{}
	netResolver := &MockNetResolver{}
	builder := getDNSBuilder(netResolver, 400*time.Millisecond)

	conn.
		On("UpdateState", resolver.State{
			Addresses: []resolver.Address{
				{Addr: "1.1.1.1:3000"},
				{Addr: "1.1.1.2:3000"},
			},
		}).
		Once().
		Return(nil)

	netResolver.
		On("LookupHost", "google.com").
		Return([]string{"1.1.1.1", "1.1.1.2"}, nil)

	netResolver.
		On("LookupSRV", "grpclb", "tcp", "google.com").
		Return("google.com", []*net.SRV{}, nil)

	r, err := builder.Build(
		resolver.Target{
			URL: *lo.Must(url.Parse("google.com")),
		},
		conn,
		resolver.BuildOptions{},
	)

	assert.NoError(err)
	assert.NotNil(r)

	netResolver.AssertExpectations(t)
	conn.AssertExpectations(t)

	netResolver.ExpectedCalls = nil
	conn.ExpectedCalls = nil

	netResolver.
		On("LookupHost", "google.com").
		Return([]string{"1.1.1.2", "1.1.1.3"}, nil)

	conn.
		On("UpdateState", resolver.State{
			Addresses: []resolver.Address{
				{Addr: "1.1.1.2:3000"},
				{Addr: "1.1.1.3:3000"},
			},
		}).
		Once().
		Return(nil)

	netResolver.
		On("LookupSRV", "grpclb", "tcp", "google.com").
		Return("google.com", []*net.SRV{}, nil)

	r.ResolveNow(resolver.ResolveNowOptions{})

	time.Sleep(1 * time.Second)

	netResolver.AssertExpectations(t)
	conn.AssertExpectations(t)
}
