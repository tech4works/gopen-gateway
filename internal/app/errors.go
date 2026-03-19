package app

import "github.com/tech4works/errors"

const (
	codeErrBackendConcurrentCancelled     = "BACKEND_CONCURRENT_CANCELLED"
	codeErrBackendBadGateway              = "BACKEND_BAD_GATEWAY"
	codeErrBackendGatewayTimeout          = "BACKEND_GATEWAY_TIMEOUT"
	codeErrBackendDependenciesNotExecuted = "BACKEND_DEPENDENCIES_NOT_EXECUTED"
	codeErrBackendBrokerNotConfigured     = "BACKEND_BROKER_NOT_CONFIGURED"
	codeErrBackendBrokerNotImplemented    = "BACKEND_BROKER_NOT_IMPLEMENTED"
)
const (
	msgErrBackendConcurrentCancelled     = "backend failed: concurrent context cancelled"
	msgErrBackendDependenciesNotExecuted = "backend failed: dependencies=%v not executed"
	msgErrBackendBadGateway              = "backend failed: bad gateway err=%s"
	msgErrBackendGatewayTimeout          = "backend failed: gateway timeout err=%s"
	msgErrBackendBrokerNotConfigured     = "backend failed: broker=%s not configured"
	msgErrBackendBrokerNotImplemented    = "backend failed: broker=%s not implemented"
)

var (
	ErrBackendBrokerNotConfigured     = errors.TargetWithCode(codeErrBackendBrokerNotConfigured)
	ErrBackendBrokerNotImplemented    = errors.TargetWithCode(codeErrBackendBrokerNotImplemented)
	ErrBackendDependenciesNotExecuted = errors.TargetWithCode(codeErrBackendDependenciesNotExecuted)
	ErrBackendBadGateway              = errors.TargetWithCode(codeErrBackendBadGateway)
	ErrBackendGatewayTimeout          = errors.TargetWithCode(codeErrBackendGatewayTimeout)
	ErrBackendConcurrentCancelled     = errors.TargetWithCode(codeErrBackendConcurrentCancelled)
)

func NewErrBackendConcurrentCancelled() error {
	return errors.NewWithSkipCallerAndCode(2, codeErrBackendConcurrentCancelled, msgErrBackendConcurrentCancelled)
}

func NewErrBackendBadGateway(err error) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrBackendBadGateway, msgErrBackendBadGateway, err)
}

func NewErrBackendGatewayTimeout(err error) error {
	return errors.NewWithSkipCallerAndCodef(2, codeErrBackendGatewayTimeout, msgErrBackendGatewayTimeout, err)
}

func NewErrBackendBrokerNotConfigured(broker string) error {
	return errors.NewWithSkipCallerAndCodef(
		2,
		codeErrBackendBrokerNotConfigured,
		msgErrBackendBrokerNotConfigured,
		broker,
	)
}

func NewErrBackendBrokerNotImplemented(broker string) error {
	return errors.NewWithSkipCallerAndCodef(
		2,
		codeErrBackendBrokerNotImplemented,
		msgErrBackendBrokerNotImplemented,
		broker,
	)
}

func NewErrBackendDependenciesNotExecuted(dependencies []string) error {
	return errors.NewWithSkipCallerAndCode(
		2,
		codeErrBackendDependenciesNotExecuted,
		msgErrBackendDependenciesNotExecuted,
		dependencies,
	)
}
