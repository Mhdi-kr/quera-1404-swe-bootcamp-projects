package service

import (
	"context"

	"example.com/authorization/pkg"
)

type AnalyticsService struct {
	cache pkg.Cache
}

func NewAnalyticsService(cache pkg.Cache) AnalyticsService {
	return AnalyticsService{
		cache: cache,
	}
}

func (as AnalyticsService) RegisterIP(ctx context.Context, ip string) error {
	err := as.cache.Client.SAdd(ctx, "ip_addresses", ip).Err()
	return err
}
