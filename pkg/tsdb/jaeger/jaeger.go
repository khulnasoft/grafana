package jaeger

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/grafana/grafana/pkg/infra/httpclient"
)

var logger = backend.NewLoggerWith("logger", "tsdb.jaeger")

type Service struct {
	im instancemgmt.InstanceManager
}

func ProvideService(httpClientProvider httpclient.Provider) *Service {
	return &Service{
		im: datasource.NewInstanceManager(newInstanceSettings(httpClientProvider)),
	}
}

type datasourceInfo struct {
	JaegerClient JaegerClient
}

func newInstanceSettings(httpClientProvider httpclient.Provider) datasource.InstanceFactoryFunc {
	return func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		httpClientOptions, err := settings.HTTPClientOptions(ctx)
		if err != nil {
			return nil, backend.DownstreamError(fmt.Errorf("error reading settings: %w", err))
		}

		httpClient, err := httpClientProvider.New(httpClientOptions)
		if err != nil {
			return nil, fmt.Errorf("error creating http client: %w", err)
		}

		if settings.URL == "" {
			return nil, backend.DownstreamError(errors.New("error reading settings: url is empty"))
		}

		logger := logger.FromContext(ctx)
		jaegerClient, err := New(settings.URL, httpClient, logger)
		return &datasourceInfo{JaegerClient: jaegerClient}, err
	}
}

func (s *Service) getDSInfo(ctx context.Context, pluginCtx backend.PluginContext) (*datasourceInfo, error) {
	i, err := s.im.Get(ctx, pluginCtx)
	if err != nil {
		return nil, err
	}

	instance, ok := i.(*datasourceInfo)
	if !ok {
		return nil, errors.New("failed to cast datasource info")
	}

	return instance, nil
}

func (s *Service) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	client, err := s.getDSInfo(ctx, backend.PluginConfigFromContext(ctx))
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	if _, err = client.JaegerClient.Services(); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
