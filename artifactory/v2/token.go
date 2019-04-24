package v2

import (
	"context"
	"fmt"
	"github.com/atlassian/go-artifactory/v2/artifactory/client"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type TokenService Service

type AccessTokenOptions struct {
	GrantType string
	Username string

	ExpiresIn int

	Groups []string

	Refreshable bool
}

type AccessTokenDescriptor struct {
	TokenId string `json:"token_id"`
	Issuer string `json:"issuer"`
	Subject string `json:"subject"`
	Expiry int `json:"expiry"`

	Refreshable bool `json:"refreshable"`
	IssuedAt int `json:"issued_at"`
}

type AccessTokenList struct {
	Tokens []AccessTokenDescriptor `json:"tokens"`
}

type CreateAccessTokenResponse struct {
    Token string `json:"access_token"`
    ExpiresIn int `json:"expires_in"`
    Scope string `json:"scope"`
    Type string `json:"token_type"`
    RefreshToken string `json:"refresh_token"`
}

func (s *TokenService) CreateAccessToken(ctx context.Context, options *AccessTokenOptions) (*CreateAccessTokenResponse, error) {
	path := fmt.Sprintf("api/security/token")
	
	data := url.Values{}
	data.Set("username", options.Username)
	data.Set("expires_in", strconv.Itoa(options.ExpiresIn))
	data.Set("scope", fmt.Sprintf("member-of-groups:\"%s\"", strings.Join(options.Groups, ",")))

	if options.GrantType != "" {
		data.Set("grant_type", options.GrantType)
	}

	if options.Refreshable {
		data.Set("refreshable", "true")
	}
	
	req, err := s.client.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", client.MediaTypeForm)
	req.Header.Set("Accept", client.MediaTypeJson)


	accessToken := CreateAccessTokenResponse{}
	_, err = s.client.Do(ctx, req, &accessToken)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (s *TokenService) GetAccessTokens(ctx context.Context) (*AccessTokenList, error) {
	path := fmt.Sprintf("api/security/token")


	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", client.MediaTypeForm)
	req.Header.Set("Accept", client.MediaTypeJson)

	accessTokenList := AccessTokenList{}
	_, err = s.client.Do(ctx, req, &accessTokenList)
	if err != nil {
		return nil, err
	}

	return &accessTokenList, nil
}

func (s *TokenService) RevokeAccessToken(ctx context.Context, token string) (*http.Response, error) {
	path := fmt.Sprintf("api/security/token/revoke")

	data := url.Values{}
	data.Set("token", token)

	req, err := s.client.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", client.MediaTypeForm)

	return s.client.Do(ctx, req, nil)
}
