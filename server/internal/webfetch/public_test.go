package webfetch

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicFetchRejectsInternalAddresses(t *testing.T) {
	var requests atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
	}))
	defer upstream.Close()
	response, err := PublicClient.Get(upstream.URL)
	require.ErrorContains(t, err, "non-public address")
	require.Nil(t, response)
	require.Zero(t, requests.Load())
	for _, address := range []string{
		"127.0.0.1", "10.0.0.1", "172.16.0.1", "192.168.0.1", "169.254.169.254",
		"100.64.0.1", "0.0.0.0", "224.0.0.1", "198.18.0.1", "::1", "fe80::1", "fc00::1",
	} {
		require.False(t, publicIP(netip.MustParseAddr(address)), address)
	}
	require.True(t, publicIP(netip.MustParseAddr("8.8.8.8")))
	require.True(t, publicIP(netip.MustParseAddr("2606:4700:4700::1111")))
}
