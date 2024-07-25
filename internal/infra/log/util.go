package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"strconv"
)

func BuildTagText(tag string) string {
	return fmt.Sprint(logger.StyleBold, tag, logger.StyleReset)
}

func BuildTraceIDText(traceID string) string {
	return fmt.Sprint(logger.StyleBold, traceID, logger.StyleReset)
}

func BuildMethodText(method string) string {
	return fmt.Sprint(methodTextStyle(method), " ", method, " ", logger.StyleReset)
}

func BuildUriText(uri string) string {
	return strconv.Quote(uri)
}

func BuildStatusCodeText(statusCode vo.StatusCode) string {
	return fmt.Sprint(statusCodeTextStyle(statusCode), " ", statusCode, " ", logger.StyleReset)
}

func statusCodeTextStyle(statusCode vo.StatusCode) string {
	if helper.IsGreaterThanOrEqual(statusCode, 200) && helper.IsLessThan(statusCode, 299) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundGreen)
	} else if helper.IsGreaterThanOrEqual(statusCode, 300) && helper.IsLessThan(statusCode, 400) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundCyan)
	} else if helper.IsGreaterThanOrEqual(statusCode, 400) && helper.IsLessThan(statusCode, 500) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundYellow)
	} else if helper.IsGreaterThanOrEqual(statusCode, 500) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundRed)
	}
	return logger.StyleBold
}

func methodTextStyle(method string) string {
	switch method {
	case http.MethodPost:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundYellow)
	case http.MethodGet:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundBlue)
	case http.MethodDelete:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundRed)
	case http.MethodPut:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundMagenta)
	case http.MethodPatch:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundCyan)
	default:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundBlack)
	}
}
