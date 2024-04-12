package consts

const (
	// XForwardedFor represents the name of the "X-Forwarded-For" HTTP header.
	// It is used to indicate the original IP address of a client connecting to a web server
	// through an HTTP proxy or a load balancer.
	XForwardedFor = "X-Forwarded-For"
	// XTraceId represents the name of the "X-Trace-Id" HTTP header.
	// It is used to uniquely identify a request or a set of related requests
	// for tracing and debugging purposes.
	XTraceId = "X-Trace-Id"
	// XGOpenCache represents the name of the "X-Gopen-Cache" HTTP header.
	// It is used to indicate whether a cache is being used for the request.
	XGOpenCache = "X-Gopen-Cache"
	// XGOpenCacheTTL represents the name of the "X-Gopen-Cache-TTL" HTTP header.
	// It is used to indicate the time-to-live (TTL) value of a cache response,
	// which specifies how long the response should be considered valid and can be reused.
	// The value of X-Gopen-Cache-TTL header is typically a duration in seconds.
	XGOpenCacheTTL = "X-Gopen-Cache-TTL"
	// XGOpenComplete represents the name of the "X-Gopen-Complete" HTTP header. It is used to indicate the completion
	// status of a request.
	XGOpenComplete = "X-Gopen-Complete"
	// XGOpenSuccess represents the name of the "X-Gopen-Success" HTTP header.
	// It is used to indicate the success status of a request.
	XGOpenSuccess = "X-Gopen-Success"
)
