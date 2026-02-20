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
	EvalAny(expr string, request *vo.HTTPRequest, history *vo.History) (any, []error)
	EvalGuards(onlyIf, ignoreIf []string, request *vo.HTTPRequest, history *vo.History) (bool, string, []error)
}

func NewDynamicValue(jsonPath domain.JSONPath) DynamicValue {
	return dynamicValueService{
		jsonPath: jsonPath,
	}
}

func (d dynamicValueService) Get(value string, request *vo.HTTPRequest, history *vo.History) (string, []error) {
	var allErrs []error

	apply := func(stage string, fn func(string, *vo.HTTPRequest, *vo.History) (string, []error)) {
		v, errs := fn(value, request, history)
		if checker.IsNotEmpty(errs) {
			for _, e := range errs {
				allErrs = append(allErrs, errors.Inheritf(e, "dynamic-value failed: op=%s", stage))
			}
		}
		value = v

	}

	apply("any-expressions", d.replaceAllAnyExpressions)
	apply("bool-expressions", d.replaceAllBoolExpressions)
	apply("all-expressions", d.replaceAllExpressions)

	return value, allErrs
}

func (d dynamicValueService) GetAsSliceOfString(value string, request *vo.HTTPRequest, history *vo.History) (
	[]string, []error) {
	newValue, errs := d.Get(value, request, history)

	if checker.IsSlice(newValue) {
		var ss []string
		if err := converter.ToDestWithErr(newValue, &ss); checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err,
				"dynamic-value failed: op=parse-slice value=%s", newValue))
		}
		return ss, errs
	}

	return []string{newValue}, errs
}

func (d dynamicValueService) EvalBool(exprs []string, request *vo.HTTPRequest, history *vo.History) (bool, []error) {
	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if checker.IsEmpty(expr) {
			continue
		}

		v, errs := d.evalBoolExpr(expr, request, history)
		if checker.IsNotEmpty(errs) {
			return false, errs
		} else if v {
			return true, nil
		}
	}
	return false, nil
}

func (d dynamicValueService) EvalAny(expr string, request *vo.HTTPRequest, history *vo.History) (any, []error) {
	expr = strings.TrimSpace(expr)

	if !(strings.HasPrefix(expr, "$") && strings.Contains(expr, "(") && strings.HasSuffix(expr, ")")) {
		return nil, errors.NewAsSlicef("dynamic-value failed: unsupported any expression expr=%s", expr)
	}

	name, args, err := d.parseFuncCall(expr)
	if checker.NonNil(err) {
		return nil, errors.InheritAsSlicef(err, "dynamic-value failed: op=parse expr=%s", expr)
	}

	v, errs := d.evalFuncAny(name, args, request, history)
	if checker.IsNotEmpty(errs) {
		for i, e := range errs {
			errs[i] = errors.Inheritf(e, "dynamic-value failed: op=eval-func-any func=%s expr=%s", name, expr)
		}
	}

	return v, errs
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
			for i, e := range errs {
				errs[i] = errors.Inheritf(e, "dynamic-value failed: guard=only-if exprs=%s", strings.Join(onlyIf, " || "))
			}
			return false, "", errs
		} else if !ok {
			return false, "only-if: " + strings.Join(onlyIf, " || "), nil
		}
	}
	if checker.IsNotEmpty(ignoreIf) {
		ignore, errs := d.EvalBool(ignoreIf, request, history)
		if checker.IsNotEmpty(errs) {
			for i, e := range errs {
				errs[i] = errors.Inheritf(e, "dynamic-value failed: guard=ignore-if exprs=%s", strings.Join(onlyIf, " || "))
			}
			return false, "", errs
		} else if ignore {
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

	tokens := d.findAllBySyntax(expr)
	for _, tok := range tokens {
		val, err := d.getValueBySyntax(tok, request, history)
		if checker.NonNil(err) {
			if errors.Is(err, mapper.ErrDynamicValueNotFound) && treatNotFoundAsEmpty {
				return "", nil
			}
			return "", errors.InheritAsSlicef(err, "dynamic-value failed: op=get-value-by-syntax token=%s expr=%s", tok, expr)
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
	if checker.IsNotEmpty(errs) {
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
		if errors.Is(err, mapper.ErrDynamicValueNotFound) {
			continue
		} else if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "dynamic-value failed: op=get-value-by-syntax word=%s", word))
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
	exprs := d.findAllFuncExpressions(value)
	if checker.IsEmpty(exprs) {
		return value, nil
	}

	var errs []error
	for _, expr := range exprs {
		b, es := d.EvalBool([]string{expr}, request, history)
		if checker.IsNotEmpty(es) {
			for _, e := range es {
				errs = append(errs, errors.Inheritf(e, "dynamic-value failed: op=eval-bool expr=%s", expr))
			}
			continue
		}
		value = strings.Replace(value, expr, converter.ToString(b), 1)
	}

	return value, errs
}

func (d dynamicValueService) replaceAllAnyExpressions(
	value string,
	request *vo.HTTPRequest,
	history *vo.History,
) (string, []error) {
	exprs := d.findAllFuncExpressions(value)
	if checker.IsEmpty(exprs) {
		return value, nil
	}

	var errs []error
	for _, expr := range exprs {
		v, es := d.EvalAny(expr, request, history)
		if checker.IsNotEmpty(es) {
			for _, e := range es {
				errs = append(errs, errors.Inheritf(e, "dynamic-value failed: op=eval-bool expr=%s", expr))
			}
			continue
		}

		repl, err := converter.ToStringWithErr(v)
		if checker.NonNil(err) {
			errs = append(errs, errors.Inheritf(err, "dynamic-value failed: op=stringify expr=%s", expr))
			continue
		}

		value = strings.Replace(value, expr, repl, 1)
	}

	return value, errs
}

func (d dynamicValueService) findAllBySyntax(value string) []string {
	// - token simples: #request.body.x
	// - responses por índice: #responses[0].body.x
	// - responses por id:     #responses[ms-auth/v1/validate].body.x
	// - id pode conter ":" (ex.: "id/asoaks/:id")
	// - coalesce: #request.body.x || #request.query.x || #responses[0].body.x || #responses[...].body.x
	//
	// Observação: permite espaços em volta do operador.
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.:\-/\[\]]+(?:\s*\|\|\s*#[a-zA-Z0-9_.:\-/\[\]]+)+|\B#[a-zA-Z0-9_.:\-/\[\]]+`)
	return regex.FindAllString(value, -1)
}

func (d dynamicValueService) findAllFuncExpressions(s string) []string {
	var out []string

	inQuotes := false
	var quote byte
	for i := 0; checker.IsLengthLessThan(s, i); i++ {
		ch := s[i]

		if (checker.Equals(ch, '"') || checker.Equals(ch, '\'')) &&
			(checker.Equals(i, 0) || checker.NotEquals(s[i-1], '\\')) {
			if !inQuotes {
				inQuotes = true
				quote = ch
			} else if checker.Equals(quote, ch) {
				inQuotes = false
				quote = 0
			}
			continue
		}
		if inQuotes {
			continue
		}

		if checker.NotEquals(ch, '$') {
			continue
		}

		start := i

		if checker.IsLessThan(i+1, len(s)) && checker.Equals(s[i+1], '(') {
			end := d.scanBalancedParens(s, i+1)
			if checker.IsGreaterThan(end, 0) {
				out = append(out, s[start:end])
				i = end - 1
			}
			continue
		}

		j := i + 1
		for checker.IsLengthLessThan(s, j) {
			c := s[j]
			if (checker.IsGreaterThanOrEqual(c, 'a') && checker.IsLessThanOrEqual(c, 'z')) ||
				(checker.IsGreaterThanOrEqual(c, 'A') && checker.IsLessThanOrEqual(c, 'Z')) ||
				(checker.IsGreaterThanOrEqual(c, '0') && checker.IsLessThanOrEqual(c, '9')) ||
				checker.Equals(c, '_') {
				j++
				continue
			}
			break
		}
		if checker.IsLessThan(j, len(s)) && checker.IsGreaterThan(j, i+1) && checker.Equals(s[j], '(') {
			end := d.scanBalancedParens(s, j)
			if checker.IsGreaterThan(end, 0) {
				out = append(out, s[start:end])
				i = end - 1
			}
			continue
		}
	}

	return out
}

func (d dynamicValueService) scanBalancedParens(s string, openIdx int) int {
	if checker.IsLessThan(openIdx, 0) ||
		checker.IsGreaterThanOrEqual(openIdx, len(s)) ||
		checker.NotEquals(s[openIdx], '(') {
		return 0
	}

	depth := 0
	inQuotes := false
	var quote byte

	for i := openIdx; checker.IsLengthLessThan(s, i); i++ {
		ch := s[i]

		if (checker.Equals(ch, '"') || checker.Equals(ch, '\'')) &&
			(checker.Equals(i, 0) || checker.NotEquals(s[i-1], '\\')) {
			if !inQuotes {
				inQuotes = true
				quote = ch
			} else if checker.Equals(quote, ch) {
				inQuotes = false
				quote = 0
			}
			continue
		}
		if inQuotes {
			continue
		}

		if checker.Equals(ch, '(') {
			depth++
			continue
		}
		if checker.Equals(ch, ')') {
			depth--
			if checker.Equals(depth, 0) {
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
			if errors.Is(err, mapper.ErrDynamicValueNotFound) {
				lastNotFound = err
				continue
			} else if checker.NonNil(err) {
				return "", errors.Inheritf(err, "dynamic-value failed: op=get-single-value-by-syntax part=%s full=%s", part, word)
			}

			return result, nil
		}

		if checker.NonNil(lastNotFound) {
			return "", lastNotFound
		}

		return "", errors.Newf("dynamic-value failed: invalid full=%s", word)
	}

	return d.getSingleValueBySyntax(word, request, history)
}

func (d dynamicValueService) getSingleValueBySyntax(word string, request *vo.HTTPRequest, history *vo.History) (string,
	error) {
	cleanSintaxe := strings.ReplaceAll(word, "#", "")
	dotSplit := strings.Split(cleanSintaxe, ".")
	if checker.IsEmpty(dotSplit) {
		return "", errors.Newf("dynamic-value failed: op=dot-split invalid token=%s", word)
	}

	prefix := dotSplit[0]
	if checker.Contains(prefix, "request") {
		return d.getRequestValueByJsonPath(cleanSintaxe, request)
	} else if checker.Contains(prefix, "responses") {
		return d.getResponseValueByJsonPath(cleanSintaxe, history)
	} else {
		return "", errors.Newf("dynamic-value failed: op=get-first-dot-split invalid-prefix=%s token=%s", prefix, word)
	}
}

func (d dynamicValueService) getRequestValueByJsonPath(jsonPath string, request *vo.HTTPRequest) (string, error) {
	jsonPath = strings.Replace(jsonPath, "request.", "", 1)

	jsonRequest, err := request.Map()
	if checker.NonNil(err) {
		return "", errors.Inheritf(err, "dynamic-value failed: op=request.map path=%s", jsonPath)
	}

	result := d.jsonPath.Get(jsonRequest, jsonPath)
	if result.Exists() && checker.IsNotEmpty(result.String()) {
		return result.String(), nil
	}

	return "", mapper.NewErrDynamicValueNotFound(jsonPath)
}

func (d dynamicValueService) getResponseValueByJsonPath(jsonPath string, history *vo.History) (string, error) {
	if strings.HasPrefix(jsonPath, "responses[") {
		return d.getResponseValueByBracketJsonPath(jsonPath, history)
	}

	return "", errors.Newf("dynamic-value failed: invalid syntax=%s", jsonPath)
}

func (d dynamicValueService) getResponseValueByBracketJsonPath(jsonPath string, history *vo.History) (string, error) {
	closeIndex := strings.Index(jsonPath, "]")
	if checker.IsLengthLessThan(closeIndex, 0) {
		return "", errors.Newf("dynamic-value failed: invalid syntax=%s missing ']'", jsonPath)
	}

	key := strings.TrimSpace(jsonPath[len("responses["):closeIndex])
	if checker.IsEmpty(key) {
		return "", errors.Newf("dynamic-value failed: invalid syntax=%s empty key", jsonPath)
	}

	rest := jsonPath[closeIndex+1:]
	if strings.HasPrefix(rest, ".") {
		rest = rest[1:]
	}

	if converter.CouldBeInt(key) {
		return d.getResponseValueByIndex(key, rest, history)
	}
	return d.getResponseValueByID(key, rest, history)
}

func (d dynamicValueService) getResponseValueByIndex(index string, rest string, history *vo.History) (string, error) {
	jsonResponses, err := history.ResponsesMap()
	if checker.NonNil(err) {
		return "", errors.Inheritf(err, "dynamic-value failed: op=history.map index=%s rest=%s", index, rest)
	}

	lookupPath := index
	if checker.IsNotEmpty(rest) {
		lookupPath = index + "." + rest
	}

	result := d.jsonPath.Get(jsonResponses, lookupPath)
	if result.Exists() && checker.IsNotEmpty(result.String()) {
		return result.String(), nil
	}

	return "", mapper.NewErrDynamicValueNotFound("responses[" + index + "]." + rest)
}

func (d dynamicValueService) getResponseValueByID(id string, rest string, history *vo.History) (string, error) {
	jsonResponses, err := history.ResponsesMapByID()
	if checker.NonNil(err) {
		return "", errors.Inheritf(err, "dynamic-value failed: op=history.map id=%s rest=%s", id, rest)
	}

	lookupPath := id
	if checker.IsNotEmpty(rest) {
		lookupPath = id + "." + rest
	}

	result := d.jsonPath.Get(jsonResponses, lookupPath)
	if result.Exists() && checker.IsNotEmpty(result.String()) {
		return result.String(), nil
	}

	return "", mapper.NewErrDynamicValueNotFound("responses[" + id + "]." + rest)
}

func (d dynamicValueService) evalBoolExpr(expr string, request *vo.HTTPRequest, history *vo.History) (bool, []error) {
	expr = strings.TrimSpace(expr)
	expr = d.trimOuterParens(expr)

	if parts, ok := d.splitTopLevel(expr, "||"); ok {
		for _, p := range parts {
			v, errs := d.evalBoolExpr(p, request, history)
			if checker.IsNotEmpty(errs) {
				return false, errs
			} else if v {
				return true, nil
			}
		}
		return false, nil
	}

	if parts, ok := d.splitTopLevel(expr, "&&"); ok {
		for _, p := range parts {
			v, errs := d.evalBoolExpr(p, request, history)
			if checker.IsNotEmpty(errs) {
				return false, errs
			} else if !v {
				return false, nil
			}
		}
		return true, nil
	}

	if strings.HasPrefix(expr, "$(") && strings.HasSuffix(expr, ")") {
		inside := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(expr, "$("), ")"))
		return d.evalAsBool(inside, request, history)
	}

	if strings.HasPrefix(expr, "$") && strings.Contains(expr, "(") && strings.HasSuffix(expr, ")") {
		name, args, err := d.parseFuncCall(expr)
		if checker.NonNil(err) {
			return false, errors.InheritAsSlicef(err, "dynamic-value failed: op=parse-func expr=%s", expr)
		}
		return d.evalFuncBool(name, args, request, history)
	}

	if b, ok := d.parseBoolLiteral(expr); ok {
		return b, nil
	}

	return false, errors.NewAsSlicef("dynamic-value failed: unsupported expr=%s", expr)
}

func (d dynamicValueService) evalFuncAny(
	name string,
	args []string,
	request *vo.HTTPRequest,
	history *vo.History,
) (any, []error) {
	n := strings.ToLower(strings.TrimSpace(name))

	need1 := func() (string, []error) {
		if checker.IsLengthNotEquals(args, 1) {
			return "", errors.NewAsSlicef("dynamic-value failed: op=need-1 $%s expects 1 argument, got %d", name,
				len(args))
		}
		return strings.TrimSpace(args[0]), nil
	}

	switch n {
	case "length":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return nil, errs
		}
		return d.evalLength(a, request, history)
	case "distinct":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return nil, errs
		}
		return d.evalDistinct(a, request, history)
	}

	return false, errors.NewAsSlicef("dynamic-value failed: unsupported func=%s", name)
}

func (d dynamicValueService) evalLength(arg string, request *vo.HTTPRequest, history *vo.History) (int, []error) {
	v, errs := d.resolveToAny(arg, request, history, true)
	if checker.IsNotEmpty(errs) {
		return 0, errs
	}

	ln, err := converter.ToLengthWithErr(v)
	if checker.NonNil(err) {
		return 0, errors.InheritAsSlicef(err, "dynamic-value failed: op=to-length arg=%s", arg)
	}

	return ln, nil
}

func (d dynamicValueService) evalDistinct(arg string, request *vo.HTTPRequest, history *vo.History) (any, []error) {
	v, errs := d.resolveToAny(arg, request, history, true)
	if checker.IsNotEmpty(errs) {
		return nil, errs
	} else if !checker.IsSliceOrArrayType(v) {
		return v, nil
	}

	arr, ok := v.([]any)
	if !ok {
		return nil, errors.NewAsSlicef("dynamic-value failed: op=cast-to-slice arg=%s", arg)
	}

	out := make([]any, 0, len(arr))
	for _, it := range arr {
		if checker.Contains(out, it) {
			continue
		}
		out = append(out, it)
	}

	return out, nil
}

func (d dynamicValueService) evalFuncBool(name string, args []string, request *vo.HTTPRequest, history *vo.History) (
	bool, []error) {
	n := strings.ToLower(strings.TrimSpace(name))

	need1 := func() (string, []error) {
		if checker.IsLengthNotEquals(args, 1) {
			return "", errors.NewAsSlicef("dynamic-value failed: op=need-1 $%s expects 1 argument, got %d", name,
				len(args))
		}
		return strings.TrimSpace(args[0]), nil
	}
	need2 := func() (string, string, []error) {
		if checker.IsLengthNotEquals(args, 2) {
			return "", "", errors.NewAsSlicef("dynamic-value failed: op=need-2 $%s expects 2 arguments, got %d",
				name, len(args))
		}
		return strings.TrimSpace(args[0]), strings.TrimSpace(args[1]), nil
	}

	switch n {
	case "isnull":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalIsNull(a, request, history, true)
	case "isnotnull":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		v, es := d.evalIsNull(a, request, history, true)
		if checker.IsNotEmpty(es) {
			return false, es
		}
		return !v, nil
	case "isempty":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalEmpty(a, request, history, true)
	case "isnullorempty":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalNullOrEmpty(a, request, history, true)
	case "isnotempty":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		v, es := d.evalEmpty(a, request, history, true)
		if checker.IsNotEmpty(es) {
			return false, es
		}
		return !v, nil
	case "isnotnullorempty":
		a, errs := need1()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		v, es := d.evalNullOrEmpty(a, request, history, true)
		if checker.IsNotEmpty(es) {
			return false, es
		}
		return !v, nil
	case "isgreaterthan":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool {
			return checker.IsGreaterThan(a, b)
		})
	case "isgreaterthanorequal":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool {
			return checker.IsGreaterThanOrEqual(a, b)
		})
	case "islessthan":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool {
			return checker.IsLessThan(a, b)
		})
	case "islessthanorequal":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalCompareNumber(l, r, request, history, func(a, b float64) bool {
			return checker.IsLessThanOrEqual(a, b)
		})
	case "equals":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalEquals(l, r, request, history, true)
	case "equalsignorecase":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalEqualsIgnoreCase(l, r, request, history, true)
	case "notequals":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		v, es := d.evalEquals(l, r, request, history, true)
		if checker.IsNotEmpty(es) {
			return false, es
		}
		return !v, nil
	case "notequalsignorecase":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		v, es := d.evalEqualsIgnoreCase(l, r, request, history, true)
		if checker.IsNotEmpty(es) {
			return false, es
		}
		return !v, nil
	case "contains":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalContains(l, r, request, history, false)
	case "containsignorecase":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalContainsIgnoreCase(l, r, request, history, false)
	case "notcontains":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalContains(l, r, request, history, true)
	case "notcontainsignorecase":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalContainsIgnoreCase(l, r, request, history, true)
	case "islengthgreaterthan":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "gt")
	case "islengthlessthan":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "lt")
	case "islengthgreaterthanorequal":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "gte")
	case "islengthlessthanorequal":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "lte")
	case "islengthequals":
		l, r, errs := need2()
		if checker.IsNotEmpty(errs) {
			return false, errs
		}
		return d.evalLengthCompare(l, r, request, history, "eq")
	}

	return false, errors.NewAsSlicef("dynamic-value failed: unsupported func=%s", name)
}

func (d dynamicValueService) evalIsNull(arg string, request *vo.HTTPRequest, history *vo.History, treatNotFoundAsEmpty bool) (bool, []error) {
	v, errs := d.resolveToAny(arg, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
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
	if checker.IsNotEmpty(errs) {
		return false, errs
	} else if checker.IsEmpty(s) {
		return true, nil
	}
	return false, nil
}

func (d dynamicValueService) evalNullOrEmpty(arg string, request *vo.HTTPRequest, history *vo.History,
	treatNotFoundAsEmpty bool) (bool, []error) {
	s, errs := d.resolveToString(arg, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
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
	if checker.IsNotEmpty(errs) {
		return false, errs
	} else if b, ok := d.parseBoolLiteral(strings.TrimSpace(s)); ok {
		return b, nil
	} else if checker.Equals(strings.TrimSpace(s), "1") {
		return true, nil
	} else if checker.Equals(strings.TrimSpace(s), "0") {
		return false, nil
	} else {
		return false, errors.NewAsSlicef("dynamic-value failed: value is not boolean: %s (resolved from %s)", s, arg)
	}
}

func (d dynamicValueService) evalCompareNumber(left string, right string, request *vo.HTTPRequest,
	history *vo.History, cmp func(a, b float64) bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, false)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, false)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	lf, err := converter.ToFloat64WithErr(strings.TrimSpace(d.stripQuotes(ls)))
	if checker.NonNil(err) {
		return false, errors.InheritAsSlicef(err, "dynamic-value failed: evalCompareNumber left-not-number value=%s", ls)
	}

	rf, err := converter.ToFloat64WithErr(strings.TrimSpace(d.stripQuotes(rs)))
	if checker.NonNil(err) {
		return false, errors.InheritAsSlicef(err, "dynamic-value failed: evalCompareNumber left-not-number value=%s", ls)
	}

	return cmp(lf, rf), nil
}

func (d dynamicValueService) evalEquals(left string, right string, request *vo.HTTPRequest, history *vo.History,
	treatNotFoundAsEmpty bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	ls = d.stripQuotes(strings.TrimSpace(ls))
	rs = d.stripQuotes(strings.TrimSpace(rs))

	return checker.Equals(ls, rs), nil
}

func (d dynamicValueService) evalEqualsIgnoreCase(left string, right string, request *vo.HTTPRequest,
	history *vo.History, treatNotFoundAsEmpty bool) (bool, []error) {
	ls, errs := d.resolveToString(left, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, treatNotFoundAsEmpty)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	ls = d.stripQuotes(strings.TrimSpace(ls))
	rs = d.stripQuotes(strings.TrimSpace(rs))

	return checker.EqualsIgnoreCase(ls, rs), nil
}

func (d dynamicValueService) evalContains(left string, right string, request *vo.HTTPRequest,
	history *vo.History, negate bool) (bool, []error) {
	lv, errs := d.resolveToAny(left, request, history, true)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rv, errs := d.resolveToAny(right, request, history, true)
	if checker.IsNotEmpty(errs) {
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
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rs, errs := d.resolveToString(right, request, history, true)
	if checker.IsNotEmpty(errs) {
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
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	rv, errs := d.resolveToAny(right, request, history, false)
	if checker.IsNotEmpty(errs) {
		return false, errs
	}

	var n int
	switch v := rv.(type) {
	case string:
		i, err := converter.ToIntWithErr(strings.TrimSpace(d.stripQuotes(v)))
		if checker.NonNil(err) {
			return false, errors.InheritAsSlicef(err,
				"dynamic-value failed: evalLengthCompare right-not-int op=%s value=%v", op, rv)
		}
		n = i
	default:
		i, err := converter.ToIntWithErr(rv)
		if checker.NonNil(err) {
			return false, errors.InheritAsSlicef(err,
				"dynamic-value failed: evalLengthCompare right-not-int op=%s value=%v", op, rv)
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
		return false, errors.NewAsSlicef("dynamic-value failed: unsupported operator=%s", op)
	}
}

func (d dynamicValueService) parseFuncCall(expr string) (string, []string, error) {
	expr = strings.TrimSpace(expr)
	if !strings.HasPrefix(expr, "$") {
		return "", nil, errors.Newf("dynamic-value failed: unsupported expr=%s", expr)
	}

	open := strings.Index(expr, "(")
	if checker.IsLessThan(open, 0) || !strings.HasSuffix(expr, ")") {
		return "", nil, errors.Newf("dynamic-value failed: op=strings-index invalid function call expr=%s", expr)
	}

	name := strings.TrimSpace(expr[1:open])
	inside := strings.TrimSpace(expr[open+1 : len(expr)-1])

	args := d.splitArgsTopLevel(inside)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	if checker.IsEmpty(name) {
		return "", nil, errors.Newf("dynamic-value failed: op=strings-index invalid function name expr=%s", expr)
	}

	return name, args, nil
}

func (d dynamicValueService) splitArgsTopLevel(s string) []string {
	var out []string
	var cur strings.Builder
	var depth int
	var inQuotes bool
	var quoteChar byte

	for i := 0; checker.IsLengthLessThan(s, i); i++ {
		ch := s[i]

		if (checker.Equals(ch, '"') || checker.Equals(ch, '\'')) &&
			(checker.Equals(i, 0) || checker.NotEquals(s[i-1], '\\')) {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if checker.Equals(quoteChar, ch) {
				inQuotes = false
				quoteChar = 0
			}
			cur.WriteByte(ch)
			continue
		}

		if !inQuotes {
			if checker.Equals(ch, '(') {
				depth++
			} else if checker.Equals(ch, ')') && checker.IsGreaterThan(depth, 0) {
				depth--
			} else if checker.Equals(ch, ',') && checker.Equals(depth, 0) {
				out = append(out, cur.String())
				cur.Reset()
				continue
			}
		}

		cur.WriteByte(ch)
	}
	if checker.NotEquals(strings.TrimSpace(cur.String()), "") {
		out = append(out, cur.String())
	}

	if checker.Equals(strings.TrimSpace(s), "") {
		return []string{}
	}
	return out
}

func (d dynamicValueService) splitTopLevel(expr string, op string) ([]string, bool) {
	var parts []string
	depth := 0
	inQuotes := false
	start := 0

	for i := 0; checker.IsLengthLessThan(expr, i); i++ {
		ch := expr[i]
		if checker.Equals(ch, '"') && (checker.Equals(i, 0) || checker.NotEquals(expr[i-1], '\\')) {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			continue
		}
		if checker.Equals(ch, '(') {
			depth++
			continue
		}
		if checker.Equals(ch, ')') && checker.IsGreaterThan(depth, 0) {
			depth--
			continue
		}
		if checker.Equals(depth, 0) && strings.HasPrefix(expr[i:], op) {
			parts = append(parts, strings.TrimSpace(expr[start:i]))
			i += len(op) - 1
			start = i + 1
		}
	}
	if checker.IsEmpty(parts) {
		return nil, false
	}

	parts = append(parts, strings.TrimSpace(expr[start:]))

	return parts, true
}

func (d dynamicValueService) trimOuterParens(s string) string {
	s = strings.TrimSpace(s)
	if checker.IsLengthLessThan(s, 2) || checker.NotEquals(s[0], '(') || checker.NotEquals(s[len(s)-1], ')') {
		return s
	}

	depth := 0
	inQuotes := false
	for i := 0; checker.IsLengthLessThan(s, i); i++ {
		ch := s[i]
		if checker.Equals(ch, '"') && (checker.Equals(i, 0) || checker.NotEquals(s[i-1], '\\')) {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			continue
		}
		if checker.Equals(ch, '(') {
			depth++
		} else if checker.Equals(ch, ')') {
			depth--
			if checker.Equals(depth, 0) && checker.NotEquals(i, len(s)-1) {
				return s
			}
		}
	}

	if checker.Equals(depth, 0) {
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
	} else if checker.Equals(s[0], '"') && checker.Equals(s[len(s)-1], '"') {
		return s[1 : len(s)-1]
	} else if checker.Equals(s[0], '\'') && checker.Equals(s[len(s)-1], '\'') {
		return s[1 : len(s)-1]
	} else {
		return s
	}
}
