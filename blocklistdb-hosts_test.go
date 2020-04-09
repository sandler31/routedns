package rdns

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestHostsDB(t *testing.T) {
	m, err := NewHostsDB(
		"# some comment",
		"              ",
		"127.0.0.1   domain1.com",
		"0.0.0.0     domain2.com",
		"0.0.0.0     domain3.com domain4.com. ",
		"::          domain5.com",
		"::1         domain6.com",
		"192.168.1.1 domain6.com",
	)
	require.NoError(t, err)

	tests := []struct {
		q     string
		typ   uint16
		match bool
		ip    net.IP
	}{
		{"domain1.com.", dns.TypeA, true, net.ParseIP("127.0.0.1")},
		{"domain2.com.", dns.TypeA, true, nil},
		{"domain4.com.", dns.TypeA, true, nil},
		{"domain5.com.", dns.TypeA, true, nil},
		{"domain6.com.", dns.TypeA, true, net.ParseIP("192.168.1.1")},
		{"domain6.com.", dns.TypeAAAA, true, net.ParseIP("::1")},
		{"domainX.com.", dns.TypeA, false, nil},
	}
	for _, test := range tests {
		q := dns.Question{Name: test.q, Qtype: test.typ, Qclass: dns.ClassINET}
		ip, _, ok := m.Match(q)
		require.Equal(t, test.match, ok, "query: %s", test.q)
		require.Equal(t, test.ip, ip, "query: %s", test.q)
	}
}
