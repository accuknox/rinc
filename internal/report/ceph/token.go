package ceph

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	types "github.com/accuknox/rinc/types/ceph"

	"github.com/golang-jwt/jwt/v5"
)

func (r *Reporter) fetchTkn(ctx context.Context) error {
	endp, err := url.JoinPath(r.conf.DashboardAPI.URL, authToken)
	if err != nil {
		return fmt.Errorf("joining url path: %w", err)
	}

	cred := types.Credential{
		Username: r.conf.DashboardAPI.Username,
		Password: r.conf.DashboardAPI.Password,
	}
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(cred)
	if err != nil {
		return fmt.Errorf("marshaling credentials to json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endp, body)
	if err != nil {
		return fmt.Errorf("creating new http request: %w", err)
	}
	req.Header.Set("accept", mediaTypeV10)
	req.Header.Set("Content-Type", "application/json")

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

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("non-201 status: %s", resp.Status)
	}

	tkn := new(types.Token)
	err = json.NewDecoder(resp.Body).Decode(tkn)
	if err != nil {
		return fmt.Errorf("decoding json body: %w", err)
	}
	r.token = &token{Token: *tkn}

	return nil
}

type token struct {
	types.Token
}

func (t token) hasExpired() (bool, error) {
	jwt, _, err := new(jwt.Parser).ParseUnverified(t.T, jwt.MapClaims{})
	if err != nil {
		return false, fmt.Errorf("parsing jwt token: %w", err)
	}
	exp, err := jwt.Claims.GetExpirationTime()
	if err != nil {
		return false, fmt.Errorf("getting token expiry: %w", err)
	}
	if time.Since(exp.Time) <= 0 {
		return true, nil
	}
	return false, nil
}

func (t token) validatePerms() bool {
	p := t.Permissions
	if p[types.PermissionHosts] == nil {
		return false
	}
	for _, action := range p[types.PermissionHosts] {
		if action == types.PermissionActionRead {
			return true
		}
	}
	return false
}
