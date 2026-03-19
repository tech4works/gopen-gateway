package factory

import "github.com/tech4works/checker"

func joinErrs(parts ...[]error) []error {
	var out []error
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func fallbackIf[T any](use bool, errs []error, cur, fallback T) T {
	if use && checker.IsNotEmpty(errs) {
		return fallback
	}
	return cur
}

func isDegraded(spec any, errs []error) bool {
	return checker.NonNil(spec) && checker.IsNotEmpty(errs)
}
