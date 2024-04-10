package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/clbanning/mxj/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

type Body struct {
	contentType enum.ContentType
	value       string
}

type CacheBody struct {
	value string
}

// NewBody returns a new Body object with the specified content type and bytes.
// If the bytes are empty, it returns an empty Body object.
// It checks if the content type passed is json and sets the content type enum accordingly.
// The value field of the Body object is set to the string representation of the bytes.
// Returns: Body - A new Body object.
func NewBody(contentType string, bytes []byte) Body {
	// se vazio, retornamos vazio
	if helper.IsEmpty(bytes) {
		return Body{}
	}
	// verificamos se o content-type passado é json
	contentTypeEnum := enum.ContentTypeText
	if helper.ContainsIgnoreCase(contentType, "application/json") {
		contentTypeEnum = enum.ContentTypeJson
	}
	// montamos o body
	return Body{
		contentType: contentTypeEnum,
		value:       string(bytes),
	}
}

func NewBodyFromString(value string) Body {
	// se vazio, retornamos nil
	if helper.IsEmpty(value) {
		return Body{}
	}
	// verificamos se a string passada é json
	contentTypeEnum := enum.ContentTypeText
	if helper.IsJson(value) {
		contentTypeEnum = enum.ContentTypeJson
	}
	// montamos o body
	return Body{
		contentType: contentTypeEnum,
		value:       value,
	}
}

func newCacheBody(body Body) *CacheBody {
	if body.IsEmpty() {
		return nil
	}
	return &CacheBody{
		value: body.value,
	}
}

func newEmptyBody() Body {
	return Body{}
}

func newErrorBody(endpointVO Endpoint, err error) Body {
	detailsErr := errors.Details(err)
	if helper.IsNil(detailsErr) {
		return newEmptyBody()
	}
	errResponseBody := &errorResponseBody{
		File:      detailsErr.GetFile(),
		Line:      detailsErr.GetLine(),
		Endpoint:  endpointVO.Path(),
		Message:   detailsErr.GetMessage(),
		Timestamp: time.Now(),
	}
	bodyBytes := helper.SimpleConvertToBytes(errResponseBody)
	return Body{
		contentType: enum.ContentTypeJson,
		value:       string(bodyBytes),
	}
}

func newBodyFromIndex(index int, backendResponseVO backendResponse) Body {
	// construímos o body padrão de resposta, com os campos iniciais
	bodyJson := "{}"
	bodyJson, _ = sjson.Set(bodyJson, "ok", backendResponseVO.Ok())
	bodyJson, _ = sjson.Set(bodyJson, "code", backendResponseVO.StatusCode())

	// construímos o body com os valores padrões
	body := Body{
		contentType: enum.ContentTypeJson,
		value:       bodyJson,
	}
	// caso seja string ou slice agregamos na chave, caso contrario, iremos agregar todos os campos json no bodyHistory
	if backendResponseVO.GroupResponse() {
		body = body.AggregateByKey(backendResponseVO.Key(index), backendResponseVO.Body())
	} else {
		body = body.Aggregate(backendResponseVO.Body())
	}
	// retornamos o body
	return body
}

func newBodyFromBackendResponse(backendResponseVO backendResponse) Body {
	// verificamos se backendResponse quer ser agrupado com o campo extra-config
	if !backendResponseVO.groupResponse {
		return backendResponseVO.body
	}
	// caso ele queira ser agrupado independente se for json ou não, transformamos ele em json
	bodyJson := "{}"
	// construímos o body vazio para poder agregar logo após
	body := Body{
		contentType: enum.ContentTypeJson,
		value:       bodyJson,
	}
	// retornamos o body agregado com a chave
	return body.AggregateByKey(backendResponseVO.Key(-1), backendResponseVO.Body())
}

func newSliceBody(slice []Body) Body {
	if helper.IsEmpty(slice) {
		return newEmptyBody()
	}
	bodyBytes := helper.SimpleConvertToBytes(slice)
	return Body{
		contentType: enum.ContentTypeJson,
		value:       string(bodyBytes),
	}
}

func (b Body) SetValue(value string) Body {
	return Body{
		contentType: b.contentType,
		value:       value,
	}
}

func (b Body) ContentType() enum.ContentType {
	return b.contentType
}

// Value returns the value of the body.
// If the body is an ordered map, it returns the ordered map value.
// If the body is a slice of ordered maps, it returns the slice of ordered maps value.
// Otherwise, it returns the default value of the body.
func (b Body) Value() string {
	return b.value
}

func (b Body) IsEmpty() bool {
	return helper.IsEmpty(b.value)
}

func (b Body) IsNotEmpty() bool {
	return helper.IsNotEmpty(b.value)
}

func (b Body) Aggregate(anotherBody Body) Body {
	mergedBodyStr := b.value
	switch b.contentType {
	case enum.ContentTypeJson:
		mergedBodyStr = mergeJSON(mergedBodyStr, anotherBody.value)
		break
	case enum.ContentTypeText:
		mergedBodyStr = mergeString(mergedBodyStr, anotherBody.value)
		break
	}
	return Body{
		contentType: b.contentType,
		value:       mergedBodyStr,
	}
}

func (b Body) AggregateByKey(key string, anotherBody Body) Body {
	if b.IsNotJson() {
		return b
	}
	mergedBodyStr := setJsonKeyValue(b.value, key, anotherBody)
	return Body{
		contentType: b.contentType,
		value:       mergedBodyStr,
	}
}

func (b Body) Interface() any {
	switch b.contentType {
	case enum.ContentTypeJson:
		return gjson.Parse(b.value).Value()
	}
	return b.value
}

func (b Body) Json() string {
	if b.IsEmpty() {
		return ""
	}

	switch b.contentType {
	case enum.ContentTypeJson:
		return b.value
	}
	return fmt.Sprintf("{\"text\": \"%v\"}", b.value)
}

func (b Body) Xml() string {
	if b.IsEmpty() {
		return ""
	}

	switch b.contentType {
	case enum.ContentTypeJson:
		mapJson, err := mxj.NewMapJson([]byte(b.value))
		if helper.IsNil(err) {
			xmlBytes, err := mapJson.XmlIndent("", "  ", "object")
			if helper.IsNil(err) {
				return string(xmlBytes)
			}
		}
		return "<object></object>"
	default:
		return fmt.Sprintf("<string>%s</string>", b.value)
	}
}

func (b Body) BytesByContentType(contentType enum.ContentType) []byte {
	switch contentType {
	case enum.ContentTypeJson:
		return []byte(b.Json())
	case enum.ContentTypeXml:
		return []byte(b.Xml())
	default:
		return []byte(b.value)
	}
}

func (b Body) Bytes() []byte {
	switch b.contentType {
	case enum.ContentTypeJson:
		return []byte(b.Json())
	default:
		return []byte(b.value)
	}
}

// String returns a string representation of the current Body instance.
// It utilizes the SimpleConvertToString function from the helper package to convert the value of the Body to a string.
// The resulting string representation of the Body is returned.
func (b Body) String() string {
	return b.value
}

func (b Body) MarshalJSON() ([]byte, error) {
	return helper.ConvertToBytes(b.value)
}

func (b Body) IsText() bool {
	return b.contentType == enum.ContentTypeText
}

func (b Body) IsJson() bool {
	return b.contentType == enum.ContentTypeJson
}

func (b Body) IsNotJson() bool {
	return !b.IsJson()
}

func (c *CacheBody) MarshalJSON() ([]byte, error) {
	return helper.ConvertToBytes(c.value)
}

func (c *CacheBody) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) {
		return nil
	}
	*c = CacheBody{
		value: string(data),
	}
	return nil
}

func mergeString(strA, strB string) string {
	return fmt.Sprintf("%s\n%s", strA, strB)
}

func mergeJSON(jsonA, jsonB string) string {
	merged := jsonA
	parsedJsonB := gjson.Parse(jsonB)
	parsedJsonB.ForEach(func(key, value gjson.Result) bool {
		merged = setJsonKeyValue(merged, key.String(), value.Value())
		return true // continue iterando
	})
	return merged
}

func setJsonKeyValue(jsonStr, key string, value any) string {
	// Se a key já existe no JSON A
	if gjson.Get(jsonStr, key).Exists() {
		jsonStr, _ = sjson.Set(jsonStr, key, []any{
			gjson.Get(jsonStr, key).Value(),
			value,
		})
	} else {
		// Caso contrário, apenas adiciona a key
		jsonStr, _ = sjson.Set(jsonStr, key, value)
	}
	return jsonStr
}
