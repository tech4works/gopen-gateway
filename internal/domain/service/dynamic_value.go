package service

import (
	"regexp"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type dynamicValueService struct {
	jsonPath domain.JSONPath
}

type DynamicValue interface {
	Get(value string, request *vo.HTTPRequest, history *vo.History) (string, []error)
	GetAsSliceOfString(value string, request *vo.HTTPRequest, history *vo.History) ([]string, []error)
	EvalBool(exprs []string, request *vo.HTTPRequest, history *vo.History) (bool, []error)
	EvalGuards(onlyIf, ignoreIf []string, request *vo.HTTPRequest, history *vo.History) (bool, string, []error)
}

func NewDynamicValue(jsonPath domain.JSONPath) DynamicValue {
	return dynamicValueService{
		jsonPath: jsonPath,
	}
}

func (d dynamicValueService) Get(value string, request *vo.HTTPRequest, history *vo.History) (string, []error) {
	var errs []error

	// 1) Primeiro resolve expressões $... embutidas na string (ex.: "$equals(...) texto #request...")
	value, errs = d.replaceAllBoolExpressions(value, request, history)
	if checker.IsLengthGreaterThan(errs, 0) {
		return value, errs
	}

	// 2) Depois resolve o restante #...
	return d.replaceAllExpressions(value, request, history)
}

func (d dynamicValueService) GetAsSliceOfString(value string, request *vo.HTTPRequest, history *vo.History) (
	[]string, []error) {
	newValue, errs := d.Get(value, request, history)
	if checker.IsSlice(newValue) {
		var ss []string
		err := converter.ToDestWithErr(newValue, &ss)
		if checker.IsNil(err) {
			return ss, errs
		} else {
			errs = append(errs, err)
		}
	}
	return []string{newValue}, errs
}

func (d dynamicValueService) EvalBool(exprs []string, request *vo.HTTPRequest, history *vo.History) (bool, []error) {
	// MVP: OR entre as expressões do array
	// - array vazio/nil => false
	// - string vazia => ignora (skip)
	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if checker.IsEmpty(expr) {
			continue
		}

		v, errs := d.evalBoolExpr(expr, request, history)
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		} else if v {
			return true, nil
		}
	}
	return false, nil
}

func (d dynamicValueService) EvalGuards(
	onlyIf,
	ignoreIf []string,
	request *vo.HTTPRequest,
	history *vo.History,
) (bool, string, []error) {
	if checker.IsNotEmpty(onlyIf) {
		ok, errs := d.EvalBool(onlyIf, request, history)
		if checker.IsNotEmpty(errs) {
			return false, "", errs
		}
		if !ok {
			return false, "only-if: " + strings.Join(onlyIf, " || "), nil
		}
	}

	if checker.IsNotEmpty(ignoreIf) {
		ignore, errs := d.EvalBool(ignoreIf, request, history)
		if checker.IsNotEmpty(errs) {
			return false, "", errs
		}
		if ignore {
			return false, "ignore-if: " + strings.Join(ignoreIf, " || "), nil
		}
	}

	return true, "", nil
}

func (d dynamicValueService) resolveToString(
	expr string,
	request *vo.HTTPRequest,
	history *vo.History,
	treatNotFoundAsEmpty bool,
) (string, []error) {
	expr = strings.TrimSpace(expr)
	expr = d.stripQuotes(expr)

	// Se o "expr" tem tokens (#...), resolvemos cada token.
	tokens := d.findAllBySyntax(expr)
	for _, tok := range tokens {
		val, err := d.getValueBySyntax(tok, request, history)
		if errors.Is(err, mapper.ErrValueNotFound) {
			if treatNotFoundAsEmpty {
				return "", nil
			}
			return "", []error{err}
		}
		if checker.NonNil(err) {
			return "", []error{err}
		}
		expr = strings.Replace(expr, tok, val, 1)
	}

	return expr, nil
}

func (d dynamicValueService) resolveToAny(
	expr string,
	request *vo.HTTPRequest,
	history *vo.History,
	treatNotFoundAsEmpty bool,
) (any, []error) {
	s, errs := d.resolveToString(expr, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return nil, errs
	}

	s = strings.TrimSpace(s)
	if checker.IsEmpty(s) {
		return "", nil
	}

	var v any
	if checker.IsNil(converter.ToDestWithErr(s, &v)) {
		return v, nil
	}

	unq := d.stripQuotes(s)

	switch strings.ToLower(strings.TrimSpace(unq)) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	if f, err := converter.ToFloat64WithErr(unq); checker.IsNil(err) {
		return f, nil
	}

	return unq, nil
}

func (d dynamicValueService) replaceAllExpressions(value string, request *vo.HTTPRequest, history *vo.History) (
	string, []error) {
	var errs []error
	for _, word := range d.findAllBySyntax(value) {
		result, err := d.getValueBySyntax(word, request, history)
		if errors.Is(err, mapper.ErrValueNotFound) {
			continue
		} else if checker.NonNil(err) {
			errs = append(errs, err)
		} else {
			value = strings.Replace(value, word, result, 1)
		}
	}
	return value, errs
}

func (d dynamicValueService) replaceAllBoolExpressions(
	value string,
	request *vo.HTTPRequest,
	history *vo.History,
) (string, []error) {
	exprs := d.findAllBoolExpressions(value)
	if len(exprs) == 0 {
		return value, nil
	}

	var errs []error
	for _, expr := range exprs {
		b, es := d.EvalBool([]string{expr}, request, history)
		if checker.IsLengthGreaterThan(es, 0) {
			errs = append(errs, es...)
			continue
		}
		value = strings.Replace(value, expr, converter.ToString(b), 1)
	}

	return value, errs
}

func (d dynamicValueService) findAllBySyntax(value string) []string {
	// Suporta:
	// - token simples: #request.body.x
	// - coalesce: #request.body.x || #request.query.x || #responses.0.body.x
	//
	// Observação: permite espaços em volta do operador.
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\-\[\]]+(?:\s*\|\|\s*#[a-zA-Z0-9_.\-\[\]]+)+|\B#[a-zA-Z0-9_.\-\[\]]+`)
	return regex.FindAllString(value, -1)
}

func (d dynamicValueService) findAllBoolExpressions(s string) []string {
	var out []string

	inQuotes := false
	var quote byte
	for i := 0; i < len(s); i++ {
		ch := s[i]

		if (ch == '"' || ch == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quote = ch
			} else if quote == ch {
				inQuotes = false
				quote = 0
			}
			continue
		}
		if inQuotes {
			continue
		}

		if ch != '$' {
			continue
		}

		start := i

		if i+1 < len(s) && s[i+1] == '(' {
			end := d.scanBalancedParens(s, i+1)
			if end > 0 {
				out = append(out, s[start:end])
				i = end - 1
			}
			continue
		}

		j := i + 1
		for j < len(s) {
			c := s[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				j++
				continue
			}
			break
		}
		if j < len(s) && j > i+1 && s[j] == '(' {
			end := d.scanBalancedParens(s, j)
			if end > 0 {
				out = append(out, s[start:end])
				i = end - 1
			}
			continue
		}
	}

	return out
}

func (d dynamicValueService) scanBalancedParens(s string, openIdx int) int {
	if openIdx < 0 || openIdx >= len(s) || s[openIdx] != '(' {
		return 0
	}

	depth := 0
	inQuotes := false
	var quote byte

	for i := openIdx; i < len(s); i++ {
		ch := s[i]

		if (ch == '"' || ch == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quote = ch
			} else if quote == ch {
				inQuotes = false
				quote = 0
			}
			continue
		}
		if inQuotes {
			continue
		}

		if ch == '(' {
			depth++
			continue
		}
		if ch == ')' {
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}

	return 0
}

func (d dynamicValueService) getValueBySyntax(word string, request *vo.HTTPRequest, history *vo.History) (string, error) {
	// Operador de fallback/coalesce: tenta da esquerda para a direita
	// Ex.: "#request.body.cpf || #request.query.cpf"
	if strings.Contains(word, "||") {
		parts := strings.Split(word, "||")
		var lastNotFound error

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if checker.IsEmpty(part) || !strings.HasPrefix(part, "#") {
				continue
			}

			result, err := d.getSingleValueBySyntax(part, request, history)
			if errors.Is(err, mapper.ErrValueNotFound) {
				lastNotFound = err
				continue
			} else if checker.NonNil(err) {
				return "", err
			}
			return result, nil
		}
		if checker.NonNil(lastNotFound) {
			return "", lastNotFound
		}
		return "", errors.Newf("Invalid dynamic value syntax! key: %s", word)
	}

	return d.getSingleValueBySyntax(word, request, history)
}

func (d dynamicValueService) getSingleValueBySyntax(word string, request *vo.HTTPRequest, history *vo.History) (string, error) {
	cleanSintaxe := strings.ReplaceAll(word, "#", "")
	dotSplit := strings.Split(cleanSintaxe, ".")
	if checker.IsEmpty(dotSplit) {
		return "", errors.Newf("Invalid dynamic value syntax! key: %s", word)
	}

	prefix := dotSplit[0]
	if checker.Contains(prefix, "request") {
		return d.getRequestValueByJsonPath(cleanSintaxe, request)
	} else if checker.Contains(prefix, "responses") {
		return d.getResponseValueByJsonPath(cleanSintaxe, history)
	} else {
		return "", errors.Newf("Invalid prefix syntax %s!", prefix)
	}
}

func (d dynamicValueService) getRequestValueByJsonPath(jsonPath string, request *vo.HTTPRequest) (string, error) {
	jsonPath = strings.Replace(jsonPath, "request.", "", 1)

	jsonRequest, err := request.Map()
	if checker.NonNil(err) {
		return "", err
	}

	result := d.jsonPath.Get(jsonRequest, jsonPath)
	if result.Exists() && checker.IsNotEmpty(result.String()) {
		return result.String(), nil
	}

	return "", mapper.NewErrValueNotFound(jsonPath)
}

func (d dynamicValueService) getResponseValueByJsonPath(jsonPath string, history *vo.History) (string, error) {
	jsonPath = strings.Replace(jsonPath, "responses.", "", 1)

	jsonResponses, err := history.ResponsesMap()
	if checker.NonNil(err) {
		return "", err
	}

	result := d.jsonPath.Get(jsonResponses, jsonPath)
	if result.Exists() && checker.IsNotEmpty(result.String()) {
		return result.String(), nil
	}

	return "", mapper.NewErrValueNotFound(jsonPath)
}

func (d dynamicValueService) evalBoolExpr(expr string, request *vo.HTTPRequest, history *vo.History) (bool, []error) {
	expr = strings.TrimSpace(expr)
	expr = d.trimOuterParens(expr)

	// OR (||) no topo
	if parts, ok := d.splitTopLevel(expr, "||"); ok {
		for _, p := range parts {
			v, errs := d.evalBoolExpr(p, request, history)
			if checker.IsLengthGreaterThan(errs, 0) {
				return false, errs
			} else if v {
				return true, nil
			}
		}
		return false, nil
	}

	// AND (&&) no topo
	if parts, ok := d.splitTopLevel(expr, "&&"); ok {
		for _, p := range parts {
			v, errs := d.evalBoolExpr(p, request, history)
			if checker.IsLengthGreaterThan(errs, 0) {
				return false, errs
			} else if !v {
				return false, nil
			}
		}
		return true, nil
	}

	// Caso: $(...)
	// Ex.: $(#request.body.enabled)
	if strings.HasPrefix(expr, "$(") && strings.HasSuffix(expr, ")") {
		inside := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(expr, "$("), ")"))
		return d.evalAsBool(inside, request, history)
	}

	// Caso: $func(...)
	// Ex.: $isEmpty(#request.body.x)
	if strings.HasPrefix(expr, "$") && strings.Contains(expr, "(") && strings.HasSuffix(expr, ")") {
		name, args, err := d.parseFuncCall(expr)
		if checker.NonNil(err) {
			return false, []error{err}
		}
		return d.evalFuncBool(name, args, request, history)
	}

	// Se vier um literal direto ("true"/"false") também aceitamos
	if b, ok := d.parseBoolLiteral(expr); ok {
		return b, nil
	}

	return false, []error{errors.Newf("Unsupported boolean expression: %s", expr)}
}

func (d dynamicValueService) evalFuncBool(name string, args []string, request *vo.HTTPRequest, history *vo.History) (
	bool, []error) {
	n := strings.ToLower(strings.TrimSpace(name))

	// helpers de aridade
	need1 := func() (string, []error) {
		if checker.IsLengthNotEquals(args, 1) {
			return "", []error{errors.Newf("$%s expects 1 argument, got %d", name, len(args))}
		}
		return strings.TrimSpace(args[0]), nil
	}
	need2 := func() (string, string, []error) {
		if checker.IsLengthNotEquals(args, 2) {
			return "", "", []error{errors.Newf("$%s expects 2 arguments, got %d", name, len(args))}
		}
		return strings.TrimSpace(args[0]), strings.TrimSpace(args[1]), nil
	}

	switch n {
	case "isnull":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalIsNull(a, request, history, true)
	case "isnotnull":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		v, es := d.evalIsNull(a, request, history, true)
		if checker.IsLengthGreaterThan(es, 0) {
			return false, es
		}
		return !v, nil
	case "isempty":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalEmpty(a, request, history, true)
	case "isnullorempty":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalNullOrEmpty(a, request, history, true)
	case "isnotempty":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		v, es := d.evalEmpty(a, request, history, true)
		if len(es) > 0 {
			return false, es
		}
		return !v, nil
	case "isnotnullorempty":
		a, errs := need1()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		v, es := d.evalNullOrEmpty(a, request, history, true)
		if len(es) > 0 {
			return false, es
		}
		return !v, nil
	case "isgreaterthan":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool { return a > b })
	case "isgreaterthanorequal":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool { return a >= b })
	case "islessthan":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool { return a < b })
	case "islessthanorequal":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool { return a <= b })
	case "equals":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalEquals(l, r, request, history, true)
	case "equalsignorecase":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalEqualsIgnoreCase(l, r, request, history, true)
	case "notequals":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		v, es := d.evalEquals(l, r, request, history, true)
		if len(es) > 0 {
			return false, es
		}
		return !v, nil
	case "notequalsignorecase":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		v, es := d.evalEqualsIgnoreCase(l, r, request, history, true)
		if len(es) > 0 {
			return false, es
		}
		return !v, nil
	case "contains":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalContains(l, r, request, history, false)
	case "containsignorecase":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalContainsIgnoreCase(l, r, request, history, false)
	case "notcontains":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalContains(l, r, request, history, true)
	case "notcontainsignorecase":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalContainsIgnoreCase(l, r, request, history, true)
	case "islengthgreaterthan":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "gt")
	case "islengthlessthan":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "lt")
	case "islengthgreaterthanorequal":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "gte")
	case "islengthlessthanorequal":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "lte")
	case "islengthequals":
		l, r, errs := need2()
		if checker.IsLengthGreaterThan(errs, 0) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "eq")
	}

	return false, []error{errors.Newf("Unsupported function: $%s", name)}
}

func (d dynamicValueService) evalIsNull(arg string, request *vo.HTTPRequest, history *vo.History, treatNotFoundAsEmpty bool) (bool, []error) {
	v, errs := d.resolveToAny(arg, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	if checker.IsNil(v) {
		return true, nil
	}

	if s, ok := v.(string); ok {
		return checker.EqualsIgnoreCase(strings.TrimSpace(d.stripQuotes(s)), "null"), nil
	}

	return false, nil
}

func (d dynamicValueService) evalEmpty(arg string, request *vo.HTTPRequest, history *vo.History, treatNotFoundAsEmpty bool,
) (bool, []error) {
	s, errs := d.resolveToString(arg, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	} else if checker.IsEmpty(s) {
		return true, nil
	}
	return false, nil
}

func (d dynamicValueService) evalNullOrEmpty(arg string, request *vo.HTTPRequest, history *vo.History,
	treatNotFoundAsEmpty bool) (bool, []error) {
	s, errs := d.resolveToString(arg, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	t := strings.TrimSpace(s)
	if checker.IsEmpty(t) {
		return true, nil
	} else if strings.EqualFold(t, "null") {
		return true, nil
	}

	return false, nil
}

func (d dynamicValueService) evalAsBool(arg string, request *vo.HTTPRequest, history *vo.History) (bool, []error) {
	s, errs := d.resolveToString(arg, request, history, false)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}
	if b, ok := d.parseBoolLiteral(strings.TrimSpace(s)); ok {
		return b, nil
	}
	// opcional: aceitar "1"/"0"
	if strings.TrimSpace(s) == "1" {
		return true, nil
	}
	if strings.TrimSpace(s) == "0" {
		return false, nil
	}
	return false, []error{errors.Newf("Value is not boolean: %s (resolved from %s)", s, arg)}
}

func (d dynamicValueService) evalCompareNumber(left string, right string, request *vo.HTTPRequest,
	history *vo.History, cmp func(a, b float64) bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, false)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, false)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	lf, err := converter.ToFloat64WithErr(strings.TrimSpace(d.stripQuotes(ls)))
	if checker.NonNil(err) {
		return false, []error{errors.Newf("Left operand is not a number: %s", ls)}
	}

	rf, err := converter.ToFloat64WithErr(strings.TrimSpace(d.stripQuotes(rs)))
	if checker.NonNil(err) {
		return false, []error{errors.Newf("Right operand is not a number: %s", rs)}
	}

	return cmp(lf, rf), nil
}

func (d dynamicValueService) evalEquals(left string, right string, request *vo.HTTPRequest, history *vo.History,
	treatNotFoundAsEmpty bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	ls = d.stripQuotes(strings.TrimSpace(ls))
	rs = d.stripQuotes(strings.TrimSpace(rs))

	return checker.Equals(ls, rs), nil
}

func (d dynamicValueService) evalEqualsIgnoreCase(left string, right string, request *vo.HTTPRequest,
	history *vo.History, treatNotFoundAsEmpty bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, treatNotFoundAsEmpty)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	ls = d.stripQuotes(strings.TrimSpace(ls))
	rs = d.stripQuotes(strings.TrimSpace(rs))

	return checker.EqualsIgnoreCase(ls, rs), nil
}

func (d dynamicValueService) evalContains(left string, right string, request *vo.HTTPRequest,
	history *vo.History, negate bool) (bool, []error) {
	lv, errs := d.resolveToAny(left, request, history, true)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rv, errs := d.resolveToAny(right, request, history, true)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	found := checker.Contains(lv, rv)
	if negate {
		return !found, nil
	}
	return found, nil
}

func (d dynamicValueService) evalContainsIgnoreCase(left string, right string, request *vo.HTTPRequest,
	history *vo.History, negate bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, true)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, true)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	ls = d.stripQuotes(strings.TrimSpace(ls))
	rs = d.stripQuotes(strings.TrimSpace(rs))

	found := checker.ContainsIgnoreCase(ls, rs)
	if negate {
		return !found, nil
	}
	return found, nil
}

func (d dynamicValueService) evalLengthCompare(left string, right string, request *vo.HTTPRequest,
	history *vo.History, op string) (bool, []error) {
	lv, errs := d.resolveToAny(left, request, history, true)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	rv, errs := d.resolveToAny(right, request, history, false)
	if checker.IsLengthGreaterThan(errs, 0) {
		return false, errs
	}

	var n int
	switch v := rv.(type) {
	case string:
		i, err := converter.ToIntWithErr(strings.TrimSpace(d.stripQuotes(v)))
		if checker.NonNil(err) {
			return false, []error{errors.Newf("Right operand must be int: %v", rv)}
		}
		n = i
	default:
		i, err := converter.ToIntWithErr(rv)
		if checker.NonNil(err) {
			return false, []error{errors.Newf("Right operand must be int: %v", rv)}
		}
		n = i
	}

	switch op {
	case "gt":
		return checker.IsLengthGreaterThan(lv, n), nil
	case "lt":
		return checker.IsLengthLessThan(lv, n), nil
	case "gte":
		return checker.IsLengthGreaterThanOrEqual(lv, n), nil
	case "lte":
		return checker.IsLengthLessThanOrEqual(lv, n), nil
	case "eq":
		return checker.IsLengthEquals(lv, n), nil
	default:
		return false, []error{errors.Newf("Invalid length operator: %s", op)}
	}
}

func (d dynamicValueService) parseFuncCall(expr string) (string, []string, error) {
	// expr: "$name(arg1, arg2)"
	expr = strings.TrimSpace(expr)
	if !strings.HasPrefix(expr, "$") {
		return "", nil, errors.Newf("Invalid function expression: %s", expr)
	}
	open := strings.Index(expr, "(")
	if open < 0 || !strings.HasSuffix(expr, ")") {
		return "", nil, errors.Newf("Invalid function call syntax: %s", expr)
	}
	name := strings.TrimSpace(expr[1:open])
	inside := strings.TrimSpace(expr[open+1 : len(expr)-1])

	args := d.splitArgsTopLevel(inside)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	if checker.IsEmpty(name) {
		return "", nil, errors.Newf("Invalid function name: %s", expr)
	}
	return name, args, nil
}

func (d dynamicValueService) splitArgsTopLevel(s string) []string {
	var out []string
	var cur strings.Builder
	var depth int
	var inQuotes bool
	var quoteChar byte

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if (ch == '"' || ch == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if quoteChar == ch {
				inQuotes = false
				quoteChar = 0
			}
			cur.WriteByte(ch)
			continue
		}

		if !inQuotes {
			if ch == '(' {
				depth++
			} else if ch == ')' && depth > 0 {
				depth--
			} else if ch == ',' && depth == 0 {
				out = append(out, cur.String())
				cur.Reset()
				continue
			}
		}

		cur.WriteByte(ch)
	}

	if strings.TrimSpace(cur.String()) != "" {
		out = append(out, cur.String())
	}
	if strings.TrimSpace(s) == "" {
		return []string{}
	}

	return out
}

func (d dynamicValueService) splitTopLevel(expr string, op string) ([]string, bool) {
	// split por op no "nível 0" (fora de parênteses e aspas)
	var parts []string
	depth := 0
	inQuotes := false
	start := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '"' && (i == 0 || expr[i-1] != '\\') {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			continue
		}
		if ch == '(' {
			depth++
			continue
		}
		if ch == ')' && depth > 0 {
			depth--
			continue
		}
		if depth == 0 && strings.HasPrefix(expr[i:], op) {
			parts = append(parts, strings.TrimSpace(expr[start:i]))
			i += len(op) - 1
			start = i + 1
		}
	}
	if len(parts) == 0 {
		return nil, false
	}
	parts = append(parts, strings.TrimSpace(expr[start:]))
	return parts, true
}

func (d dynamicValueService) trimOuterParens(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return s
	}

	// remove apenas se os parênteses externas "abraçam" tudo
	depth := 0
	inQuotes := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' && (i == 0 || s[i-1] != '\\') {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			continue
		}
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			// se fechou antes do fim, não é par externo
			if depth == 0 && i != len(s)-1 {
				return s
			}
		}
	}
	if depth == 0 {
		return strings.TrimSpace(s[1 : len(s)-1])
	}
	return s
}

func (d dynamicValueService) parseBoolLiteral(s string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(d.stripQuotes(s))) {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

func (d dynamicValueService) stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if checker.IsLengthLessThan(s, 2) {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	if s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}

	return s
}
