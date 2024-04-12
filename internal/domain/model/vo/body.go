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

// NewBody creates a new instance of Body based on the provided contentType and bytes.
// If the bytes are empty, it returns an empty Body.
// If the contentType is "application/json", it sets the contentTypeEnum to ContentTypeJson.
// Otherwise, it sets the contentTypeEnum to ContentTypeText.
// It returns the constructed Body instance.
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

// NewBodyFromString creates a new instance of Body based on the provided value.
// If the value is empty, it returns an empty Body.
// If the value is a valid JSON string, it sets the contentTypeEnum to ContentTypeJson.
// Otherwise, it sets the contentTypeEnum to ContentTypeText.
// It returns the constructed Body instance.
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

// newCacheBody creates a new instance of CacheBody based on the provided body.
// If the body is empty, it returns nil.
// It sets the value of the CacheBody based on the value of the body.
// It returns a pointer to the constructed CacheBody instance.
func newCacheBody(body Body) *CacheBody {
	if body.IsEmpty() {
		return nil
	}
	return &CacheBody{
		value: body.value,
	}
}

// newEmptyBody creates a new instance of the Body struct with empty
// values for contentType and value. It returns the constructed Body instance.
func newEmptyBody() Body {
	return Body{}
}

// newErrorBody creates a new instance of Body based on the provided endpointVO and err.
// It generates a JSON response body containing error details using the errorResponseBody struct.
// If error details are nil or empty, it returns an empty Body.
// Otherwise, it populates the errorResponseBody fields with error details and current timestamp.
// It then converts the errorResponseBody to bytes and constructs the Body instance with ContentTypeJson and bodyBytes.
// It returns the constructed Body instance.
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

// newBodyFromIndex creates a new instance of Body based on the provided index and backendResponse.
// It constructs the default response body with initial fields.
// If the backendResponseVO is marked for grouping, it aggregates the body based on the index and body value.
// Otherwise, it aggregates all JSON fields in the bodyHistory.
// It returns the constructed Body instance.
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

// newBodyFromBackendResponse creates a new instance of Body based on the provided backendResponseVO.
// If backendResponseVO.groupResponse is false, it returns backendResponseVO.body.
// If backendResponseVO.groupResponse is true, it creates a new Body with contentType set to ContentTypeJson and value set to "{}".
// It then calls body.AggregateByKey() to aggregate the body with the provided key and anotherBody.
// The aggregated Body instance is returned.
func newBodyFromBackendResponse(backendResponseVO backendResponse) Body {
	// verificamos se backendResponse quer ser agrupado com o campo extra-config
	if !backendResponseVO.group {
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

// newSliceBody creates a new instance of Body based on the provided slice of Body.
// If the slice is empty, it returns an empty Body.
// It converts the provided slice of Body to bytes using the SimpleConvertToBytes function from the helper package.
// It sets the contentTypeEnum to ContentTypeJson.
// It constructs the Body instance with the contentTypeEnum and the string value of the converted bytes.
// It returns the constructed Body instance.
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

// SetValue returns a new instance of Body with the provided value.
// The new Body instance will have the same contentType as the original Body instance.
func (b Body) SetValue(value string) Body {
	return Body{
		contentType: b.contentType,
		value:       value,
	}
}

// ContentType returns the value of the contentType field in the Body struct.
func (b Body) ContentType() enum.ContentType {
	return b.contentType
}

// Value returns the value of the Body instance.
func (b Body) Value() string {
	return b.value
}

// IsEmpty returns a boolean value indicating whether the value of the Body is empty or not.
// It uses the helper.IsEmpty function to determine the emptiness of the value.
func (b Body) IsEmpty() bool {
	return helper.IsEmpty(b.value)
}

// IsNotEmpty returns a boolean value indicating whether the value of the Body is not empty.
// It utilizes the helper.IsNotEmpty function to determine the non-emptiness of the value.
func (b Body) IsNotEmpty() bool {
	return helper.IsNotEmpty(b.value)
}

// Aggregate returns a new instance of Body by merging the value of the current Body
// instance with the value of anotherBody.
// The merging process depends on the contentType of the current Body instance.
// If contentType is enum.ContentTypeJson, the JSON values of the two bodies will be merged.
// If contentType is enum.ContentTypeText, the two string values will be concatenated with a newline separator.
// The new Body instance will have the same contentType as the original Body instance.
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

// AggregateByKey returns a new instance of Body by merging the value of the current Body
// instance with the value of anotherBody. The merging process depends on the contentType
// of the current Body instance. If contentType is enum.ContentTypeJson, the JSON values of
// the two bodies will be merged by adding the key-value pair to the existing JSON. If the
// key already exists in the JSON, the value will be appended to the existing array of values
// under that key. If the key does not exist, it will be added with the provided value.
// If contentType is not enum.ContentTypeJson, the method returns the current Body instance without
// any modifications. The new Body instance will have the same contentType as the original Body instance.
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

// Interface returns the interface representation of the Body object.
// If the contentType of the Body is enum.ContentTypeJson, it parses the value as JSON
// using gjson.Parse and returns the rooted JSON.Value.
// For any other contentType, it returns the value as it is.
// The returned value will have the type `interface{}`.
func (b Body) Interface() any {
	switch b.contentType {
	case enum.ContentTypeJson:
		return gjson.Parse(b.value).Value()
	}
	return b.value
}

// Json returns the JSON representation of the Body instance.
// If the Body is empty, an empty string is returned.
// If the contentType of the Body is enum.ContentTypeJson, the value is returned as it is.
// For any other contentType, the value is formatted as a JSON string with a "text" field.
// The resulting JSON string is returned.
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

// Xml returns the XML representation of the Body instance.
// If the Body is empty, an empty string is returned.
// If the contentType of the Body is enum.ContentTypeJson, the JSON value is converted to XML
// using mxj.NewMapJson and mapJson.XmlIndent functions.
// If the conversion is successful, the XML bytes are returned as a string.
// If the conversion fails, "<object></object>" is returned.
// For any other contentType, the value is formatted as a string with an XML tag wrapping it.
// The resulting XML string is returned.
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

// BytesByContentType returns the byte representation of the Body instance
// based on the provided contentType. If the contentType is ContentTypeJson,
// it returns the byte representation of the JSON value of the Body. If the
// contentType is ContentTypeXml, it returns the byte representation of the
// XML value of the Body. For any other contentType, it returns the byte
// representation of the value of the Body.
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

// Bytes returns the byte representation of the Body instance.
// If the contentType of the Body is enum.ContentTypeJson, the JSON value is returned as bytes.
// For any other contentType, the value is returned as bytes.
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

// MarshalJSON marshals the value of the Body instance into JSON format.
// It uses helper.ConvertToBytes to convert the value to bytes and return it along with no errors.
func (b Body) MarshalJSON() ([]byte, error) {
	return helper.ConvertToBytes(b.value)
}

// IsText returns a boolean value indicating whether the contentType of the Body is ContentTypeText.
func (b Body) IsText() bool {
	return helper.Equals(b.contentType, enum.ContentTypeText)
}

// IsJson returns a boolean value indicating whether the contentType of the Body is ContentTypeJson.
func (b Body) IsJson() bool {
	return helper.Equals(b.contentType, enum.ContentTypeJson)
}

// IsNotJson returns a boolean value indicating whether the contentType of the Body is not ContentTypeJson.
func (b Body) IsNotJson() bool {
	return !b.IsJson()
}

// MarshalJSON returns the JSON encoding of CacheBody's value using helper.ConvertToBytes.
// It implements the json.Marshaler interface.
func (c *CacheBody) MarshalJSON() ([]byte, error) {
	return helper.ConvertToBytes(c.value)
}

// UnmarshalJSON unmarshals the given JSON data into a CacheBody instance.
// If the data is empty, it returns nil.
// Otherwise, it sets the value of the CacheBody instance to the string representation of the data.
// The function always returns nil as error.
func (c *CacheBody) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) {
		return nil
	}
	*c = CacheBody{
		value: string(data),
	}
	return nil
}

// mergeString concatenates two strings with a newline separator and returns the result.
func mergeString(strA, strB string) string {
	return fmt.Sprintf("%s\n%s", strA, strB)
}

// mergeJSON merges two JSON strings together.
// The function iterates over each key-value pair in jsonB and adds it to jsonA.
// If a key already exists in jsonA, the value is appended to the existing value as an array.
// If a key does not exist in jsonA, it is added with the provided value.
// The merged JSON string is returned.
func mergeJSON(jsonA, jsonB string) string {
	merged := jsonA
	parsedJsonB := gjson.Parse(jsonB)
	parsedJsonB.ForEach(func(key, value gjson.Result) bool {
		merged = setJsonKeyValue(merged, key.String(), value.Value())
		return true // continue iterando
	})
	return merged
}

// setJsonKeyValue sets the value of a key in a JSON string.
// If the key already exists in the JSON string, the value is appended to the existing array of values under that key.
// If the key does not exist, it is added with the provided value.
// The updated JSON string is returned.
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
