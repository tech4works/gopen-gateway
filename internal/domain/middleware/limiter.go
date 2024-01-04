package middleware

import (
	"bytes"
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/handler"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/infra/rate"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type limiter struct {
	maxSizeReqBody         int64
	maxSizeMultipartMemory int64
	ipRateLimiter          *rate.IPRateLimiter
}

type Limiter interface {
	PreHandlerRequest(ctx *gin.Context)
}

func NewLimiter(maxSizeRequestBody, maxSizeMultipartMemory, maxIpRequestPerSeconds int) Limiter {
	maxRequestPerSeconds := 5
	if maxIpRequestPerSeconds > 0 {
		maxRequestPerSeconds = maxIpRequestPerSeconds
	}
	return limiter{
		maxSizeReqBody:         int64(maxSizeRequestBody),
		maxSizeMultipartMemory: int64(maxSizeMultipartMemory),
		ipRateLimiter:          rate.NewIPRateLimiter(1, maxRequestPerSeconds),
	}
}

func (l limiter) PreHandlerRequest(ctx *gin.Context) {
	ipLimiter := l.ipRateLimiter.GetLimiter(ctx.ClientIP())
	if !ipLimiter.Allow() {
		handler.RespondCode(ctx, http.StatusTooManyRequests)
		return
	}
	maxBytesReader := l.maxSizeReqBody
	if strings.Contains(ctx.GetHeader("Content-Type"), "multipart/form-data") {
		maxBytesReader = l.maxSizeMultipartMemory
	}
	r := http.MaxBytesReader(nil, ctx.Request.Body, maxBytesReader)
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		handler.RespondCodeWithError(ctx, http.StatusRequestEntityTooLarge, errors.New(
			"request body too large! limit: ", strconv.Itoa(int(maxBytesReader))+"B",
		))
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	ctx.Next()
}
