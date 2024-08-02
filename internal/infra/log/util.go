package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"net/http"
	"strconv"
)

func BuildLevelText(lvl level) string {
	return fmt.Sprint(StyleBold, lvl.String(), StyleReset)
}

func BuildTagText(tag string) string {
	return fmt.Sprint(StyleBold, tag, StyleReset)
}

func BuildTraceIDText(traceID string) string {
	return traceID
}

func BuildMethodText(method string) string {
	return fmt.Sprint(methodTextStyle(method), " ", method, " ", StyleReset)
}

func BuildUriText(uri string) string {
	return strconv.Quote(uri)
}

func BuildStatusCodeText(statusCode vo.StatusCode) string {
	return fmt.Sprint(statusCodeTextStyle(statusCode.Code()), " ", statusCode.Code(), " ", StyleReset)
}

func statusCodeTextStyle(code int) string {
	if helper.IsGreaterThanOrEqual(code, 200) && helper.IsLessThan(code, 299) {
		return BackgroundGreen
	} else if helper.IsGreaterThanOrEqual(code, 300) && helper.IsLessThan(code, 400) {
		return BackgroundCyan
	} else if helper.IsGreaterThanOrEqual(code, 400) && helper.IsLessThan(code, 500) {
		return BackgroundYellow
	} else if helper.IsGreaterThanOrEqual(code, 500) {
		return BackgroundRed
	}
	return StyleBold
}

func methodTextStyle(method string) string {
	switch method {
	case http.MethodPost:
		return fmt.Sprint(BackgroundYellow)
	case http.MethodGet:
		return fmt.Sprint(BackgroundBlue)
	case http.MethodDelete:
		return fmt.Sprint(BackgroundRed)
	case http.MethodPut:
		return fmt.Sprint(BackgroundMagenta)
	case http.MethodPatch:
		return fmt.Sprint(BackgroundCyan)
	default:
		return fmt.Sprint(BackgroundBlack)
	}
}
