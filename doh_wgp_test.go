package ctrld

import (
	"context"
	"net/http"
	"testing"
)

// Webguard Pro fork: a custom DoH upstream with send_client_info = true must receive the
// client's MAC/IP/hostname, so the portal can attribute queries to devices.
func TestAddHeaderCustomUpstreamSendsClientInfo(t *testing.T) {
	on := true
	uc := &UpstreamConfig{Name: "wgp", Type: ResolverTypeDOH, Endpoint: "https://dns.example.com/dns-query", SendClientInfo: &on}
	uc.Init()
	ci := &ClientInfo{Mac: "aa:bb:cc:dd:ee:ff", IP: "192.168.8.20", Hostname: "kitchen"}
	ctx := context.WithValue(context.Background(), ClientInfoCtxKey{}, ci)
	req, _ := http.NewRequest(http.MethodPost, uc.Endpoint, nil)
	addHeader(ctx, req, uc)
	for h, want := range map[string]string{dohMacHeader: ci.Mac, dohIPHeader: ci.IP, dohHostHeader: ci.Hostname} {
		if got := req.Header.Get(h); got != want {
			t.Errorf("header %s = %q, want %q", h, got, want)
		}
	}
}
