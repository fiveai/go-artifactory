package v2

import (
	"context"
	"io"
	"net/http"
)

type ConfigurationService Service

func (s *ConfigurationService) ApplyConfiguration(ctx context.Context, body io.Reader) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPatch, "/api/system/configuration", body)

	if err != nil || req == nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/yaml")

	return s.client.Do(ctx, req, nil)
}