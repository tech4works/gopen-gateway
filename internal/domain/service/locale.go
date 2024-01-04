package service

import (
	"context"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/infra/geolocalization"
)

type locale struct {
	geoClient geolocalization.Client
}

type Locale interface {
	GetLocaleByIP(ctx context.Context, ip string) (*dto.IPLocale, error)
}

func NewLocale(geoClient geolocalization.Client) Locale {
	return locale{
		geoClient: geoClient,
	}
}

func (l locale) GetLocaleByIP(ctx context.Context, ip string) (*dto.IPLocale, error) {
	return l.geoClient.GetLocaleByIP(ctx, ip)
}
