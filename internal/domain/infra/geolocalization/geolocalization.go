package geolocalization

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"io"
	"net/http"
	"os"
	"strconv"
)

type client struct {
}

type Client interface {
	GetLocaleByIP(ctx context.Context, ip string) (*dto.IPLocale, error)
}

func NewClient() Client {
	return client{}
}

func (g client) GetLocaleByIP(ctx context.Context, ip string) (*dto.IPLocale, error) {
	url := os.Getenv("IP_GEOLOCALIZATION_API_URL")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("api_key", os.Getenv("IP_GEOLOCALIZATION_API_KEY"))
	if !helper.IsPrivateIP(ip) {
		q.Add("ip_address", ip)
	}
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Error("GetLocaleByIP error close body:", err)
		}
	}(res.Body)

	var locale dto.IPLocale
	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("response api status code != OK code: " + strconv.Itoa(res.StatusCode))
	} else if res.StatusCode == http.StatusOK {
		bytesBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytesBody, &locale)
		if err != nil {
			return nil, err
		}
	}
	return &locale, nil
}
