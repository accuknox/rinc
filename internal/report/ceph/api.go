package ceph

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	summaryEndpoint       = "/api/summary"
	healthEndpoint        = "/api/health/full"
	hostListEndpoint      = "/api/host"
	hostInventoryEndpoint = "/ui-api/host/inventory"
	hostDevicesEndpoint   = "/api/host/%s/devices"
	authToken             = "/api/auth"
	authLogout            = "/api/auth/logout"
)

const (
	mediaTypeV10 = "application/vnd.ceph.api.v1.0+json"
	mediaTypeV13 = "application/vnd.ceph.api.v1.3+json"
)

func (r Reporter) call(ctx context.Context, endp, mediaTyp string, v any, q ...url.Values) error {
	endp, err := url.JoinPath(r.conf.DashboardAPI.URL, endp)
	if err != nil {
		return fmt.Errorf("joining url path: %w", err)
	}
	u, err := url.Parse(endp)
	if err != nil {
		return fmt.Errorf("parsing url %q: %w", endp, err)
	}
	if len(q) != 0 {
		u.RawQuery = q[0].Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("creating new http request: %w", err)
	}

	if r.token == nil {
		err := r.fetchTkn(ctx)
		if err != nil {
			return fmt.Errorf("fetching auth token: %w", err)
		}
		if !r.token.validatePerms() {
			return fmt.Errorf("user %q doesn't have sufficient permissions",
				r.conf.DashboardAPI.Username)
		}
	}

	expired, err := r.token.hasExpired()
	if err != nil {
		return fmt.Errorf("validating auth token expiry: %w", err)
	}
	if !expired {
		err := r.fetchTkn(ctx)
		if err != nil {
			return fmt.Errorf("fetching auth token: %w", err)
		}
		if !r.token.validatePerms() {
			return fmt.Errorf("user %q doesn't have sufficient permissions",
				r.conf.DashboardAPI.Username)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token.T))
	req.Header.Set("accept", mediaTyp)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ceph dashboard api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if resp.StatusCode != 200 {
		return fmt.Errorf("decoding json body: %w", err)
	}

	return nil
}
