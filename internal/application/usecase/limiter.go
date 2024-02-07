package usecase

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/infra"
	"io"
	"net/http"
)

type limiter struct {
	maxSizeReqBody         int
	maxSizeMultipartMemory int
	ipRateLimiter          *infra.IpRateLimiter
}

type Limiter interface {
	ValidateRate(ip string) error
	ValidateRequestSize(contentType string, body io.ReadCloser) ([]byte, error)
}

func NewLimiter(maxSizeRequestBody, maxSizeMultipartMemory, maxIpRequestPerSeconds int) Limiter {
	maxRequestPerSeconds := 5
	if helper.IsGreaterThan(maxIpRequestPerSeconds, 0) {
		maxRequestPerSeconds = maxIpRequestPerSeconds
	}
	return limiter{
		maxSizeReqBody:         maxSizeRequestBody,
		maxSizeMultipartMemory: maxSizeMultipartMemory,
		ipRateLimiter:          infra.NewIpRateLimiter(1, maxRequestPerSeconds),
	}
}

func (l limiter) ValidateRate(ip string) (err error) {
	ipLimiter := l.ipRateLimiter.GetLimiter(ip)
	if !ipLimiter.Allow() {
		err = errors.New("too many requests by ip:", ip)
	}
	return err
}

func (l limiter) ValidateRequestSize(headerContentType string, body io.ReadCloser) ([]byte, error) {
	defer l.closeBody(body)
	maxBytesReader := l.maxSizeReqBody
	if helper.ContainsIgnoreCase(headerContentType, "multipart/form-data") {
		maxBytesReader = l.maxSizeMultipartMemory
	}
	read := http.MaxBytesReader(nil, body, int64(maxBytesReader))
	bodyBytes, err := io.ReadAll(read)
	if helper.IsNotNil(err) {
		return nil, errors.New("request body too large! limit:", maxBytesReader, "bytes")
	}
	return bodyBytes, nil
}

func (l limiter) closeBody(body io.ReadCloser) {
	_ = body.Close()
}
