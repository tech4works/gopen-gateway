/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vo

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/clbanning/mxj/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

// Body represents the content and format of an HTTP request or response body.
type Body struct {
	// contentType is an enumeration type that represents the format of the content.
	// It can have the following values: ContentTypeText, ContentTypeJson.
	contentType enum.ContentType
	// value represents the content of an HTTP request or response body. It is stored as a bytes.Buffer object.
	value *bytes.Buffer
}

// CacheBody represents the caching value of an HTTP response body.
type CacheBody struct {
	// ContentType represents the format of the content.
	ContentType enum.ContentType `json:"content-type,omitempty"`
	// Value represents the caching content of an HTTP response body.
	// It is a pointer to the CacheBodyValue type, which is an alias for bytes.Buffer.
	// The value is nullable and is omitted in JSON if it is empty.
	Value *CacheBodyValue `json:"value,omitempty"`
}

// CacheBodyValue is an alias for bytes.Buffer type used to represent the caching value
// of an HTTP response body. It contains methods to convert the value to different
// representations, such as string and JSON.
type CacheBodyValue bytes.Buffer

// NewBody creates a new instance of Body based on the provided contentType and buffer.
// If the buffer is empty, it returns nil.
// If the contentType contains the string "application/json", it sets the contentTypeEnum to ContentTypeJson.
// If the contentType contains the string "text/plain", it sets the contentTypeEnum to ContentTypeText.
// Otherwise, the contentTypeEnum remains uninitialized.
// It returns a pointer to the constructed Body instance.
func NewBody(contentType string, buffer *bytes.Buffer) *Body {
	// se vazio, retornamos vazio
	if helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	// montamos o body
	return &Body{
		contentType: enum.ContentTypeFromString(contentType),
		value:       buffer,
	}
}

// newErrorBody creates a new instance of Body as an error response body.
// It takes a path string and an error object as arguments.
// First, it obtains the error details using the errors.Details function from the go-errors library.
// If the detailsErr is nil, it returns nil.
// Then, it constructs the errorResponseBody object using the error details and the provided path.
// After that, it creates the Body instance by setting the contentType to ContentTypeJson and
// converting the errorResponseBody object to a Buffer using the helper.SimpleConvertToBuffer function.
// Finally, it returns a pointer to the constructed Body instance.
func newErrorBody(path string, err error) *Body {
	// obtemos o detalhe do erro usando a lib go-errors
	detailsErr := errors.Details(err)
	if helper.IsNil(detailsErr) {
		return nil
	}

	// com os detalhes, construímos o objeto de retorno padrão de erro da API Gateway
	errResponseBody := errorResponseBody{
		File:      detailsErr.GetFile(),
		Line:      detailsErr.GetLine(),
		Endpoint:  path,
		Message:   detailsErr.GetMessage(),
		Timestamp: time.Now(),
	}

	// construímos o body com esse objeto
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(errResponseBody),
	}
}

// newBodyFromIndexAndBackendResponse creates a new instance of Body based on the provided index and backendResponse.
// It constructs the default response body with initial fields.
// If the backendResponseVO is marked for grouping, it aggregates the body based on the index and body value.
// Otherwise, it aggregates all JSON fields in the bodyHistory.
// It returns the constructed Body instance.
func newBodyFromIndexAndBackendResponse(index int, backendResponseVO *backendResponse) *Body {
	// construímos o body padrão de resposta, com os campos iniciais
	bodyJson := "{}"
	bodyJson, _ = sjson.Set(bodyJson, "ok", backendResponseVO.Ok())
	bodyJson, _ = sjson.Set(bodyJson, "code", backendResponseVO.StatusCode())

	// construímos o body com os valores padrões
	body := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(bodyJson),
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
func newBodyFromBackendResponse(backendResponseVO *backendResponse) *Body {
	// verificamos se backendResponse quer ser agrupado com o campo extra-config
	if !backendResponseVO.group {
		return backendResponseVO.body
	}
	// construímos o body vazio para poder agregar logo após
	body := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer("{}"),
	}
	// retornamos o body agregado com a chave
	return body.AggregateByKey(backendResponseVO.Key(-1), backendResponseVO.Body())
}

// newSliceBody creates a new instance of Body based on the provided slice of *Body.
// If the slice is empty, it returns nil.
// The contentType of the new Body instance is set to ContentTypeJson.
// The value of the new Body instance is obtained by converting the slice to a bytes.Buffer object using helper.SimpleConvertToBuffer.
// It returns a pointer to the constructed Body instance.
func newSliceBody(sliceOfBodies []*Body) *Body {
	if helper.IsEmpty(sliceOfBodies) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(sliceOfBodies),
	}
}

// newBodyFromCacheBody creates a new instance of Body based on the provided CacheBody.
// If the cacheBody parameter is nil, it returns nil.
// Otherwise, it sets the contentType field of the new Body instance to the ContentType field of the cacheBody parameter.
// It converts the Value field of the cacheBody parameter to a bytes.Buffer pointer and assigns it to the value field of the new Body instance.
// It returns a pointer to the newly created Body instance.
func newBodyFromCacheBody(cacheBodyVO *CacheBody) *Body {
	if helper.IsNil(cacheBodyVO) {
		return nil
	}
	return &Body{
		contentType: cacheBodyVO.ContentType,
		value:       (*bytes.Buffer)(cacheBodyVO.Value),
	}
}

// newCacheBody creates a new instance of CacheBody based on the provided body.
// If the body is nil, it returns nil.
// Otherwise, it sets the ContentType field of CacheBody based on the ContentType method of body.
// It sets the Value field of CacheBody by calling newCacheBodyValue with the Value method of body as an argument.
// It returns a pointer to the constructed CacheBody instance.
func newCacheBody(bodyVO *Body) *CacheBody {
	if helper.IsNil(bodyVO) {
		return nil
	}
	return &CacheBody{
		ContentType: bodyVO.ContentType(),
		Value:       newCacheBodyValue(bodyVO.Value()),
	}
}

// newCacheBodyValue creates a new instance of CacheBodyValue based on the provided buffer.
// If the buffer is nil, it returns nil.
// Otherwise, it converts the buffer to a pointer of type *CacheBodyValue.
// It returns a pointer to the constructed CacheBodyValue instance.
func newCacheBodyValue(buffer *bytes.Buffer) *CacheBodyValue {
	if helper.IsNil(buffer) {
		return nil
	}
	return (*CacheBodyValue)(buffer)
}

// SetValue returns a new instance of Body with the provided value.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) SetValue(value *bytes.Buffer) *Body {
	return &Body{
		contentType: b.contentType,
		value:       value,
	}
}

// ContentType returns the value of the contentType field in the Body struct.
func (b *Body) ContentType() enum.ContentType {
	return b.contentType
}

// Value returns the value of the `bytes.Buffer` object stored in the `value` field of the Body struct.
func (b *Body) Value() *bytes.Buffer {
	return b.value
}

// Aggregate merges two Body instances into a new Body instance.
// If anotherBody is nil, it returns the current Body instance.
// It converts the values of the bodies to strings and merges them according to the contentType field.
// If the contentType is ContentTypeJson, it uses the mergeJSON function to merge the JSON strings.
// If the contentType is ContentTypeText, it uses the mergeString function to merge the strings.
// The mergedBody value is then used to create a new Body instance with the same contentType,
// which is returned as the result.
func (b *Body) Aggregate(anotherBody *Body) *Body {
	if helper.IsNil(anotherBody) {
		return b
	}

	// instanciamos o valor dos bodies em string para ser manipulado
	bodyStr := b.String()
	anotherBodyStr := anotherBody.String()

	mergedBody := b.value
	switch b.contentType {
	case enum.ContentTypeJson:
		mergedBody = mergeJSON(bodyStr, anotherBodyStr)
		break
	case enum.ContentTypeText:
		mergedBody = mergeString(bodyStr, anotherBodyStr)
		break
	}

	// construímos o novo ponteiro do body merge
	return &Body{
		contentType: b.contentType,
		value:       mergedBody,
	}
}

// AggregateByKey merges the value of the current Body instance with the value of anotherBody.
// It only performs the merging operation if the current Body instance has a contentType of ContentTypeJson and anotherBody is not nil.
// If either of the conditions is not satisfied, it returns the current Body instance.
// The merging operation is done by setting the provided key to the value of anotherBody in the JSON representation of the current Body instance.
// The resulting merged JSON string is converted to a buffer and used to create a new Body instance with the same contentType as the current instance.
// The new Body instance is then returned.
func (b *Body) AggregateByKey(key string, anotherBody *Body) *Body {
	if b.IsNotJson() || helper.IsNil(anotherBody) {
		return b
	}

	var value any = anotherBody
	if anotherBody.ContentType() == enum.ContentTypeText {
		value = anotherBody.String()
	}

	mergedBodyStr := setJsonKeyValue(b.String(), key, value)
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mergedBodyStr),
	}
}

// Interface returns the interface representation of the Body object.
// If the contentType of the Body is enum.ContentTypeJson, it parses the value as JSON
// using gjson.Parse and returns the rooted JSON.Value.
// For any other contentType, it returns the value as it is.
// The returned value will have the type `interface{}`.
func (b *Body) Interface() any {
	switch b.contentType {
	case enum.ContentTypeJson:
		return gjson.ParseBytes(b.value.Bytes()).Value()
	}
	return b.value.String()
}

// Json returns the byte representation of the Body instance in JSON format.
// If the contentType of the Body is ContentTypeText, it converts the value to a JSON string
// with the `text` field containing the value. If the contentType is ContentTypeJson,
// it returns the same value as the Bytes method. For any other contentType, it returns an empty byte array.
func (b *Body) Json() []byte {
	switch b.contentType {
	case enum.ContentTypeText:
		return helper.SimpleConvertToBytes(fmt.Sprintf("{\"text\": \"%v\"}", b.value))
	case enum.ContentTypeJson:
		return b.Bytes()
	default:
		return []byte{}
	}
}

// Xml returns the XML representation of the Body instance.
// If the Body is empty, an empty string is returned.
// If the contentType of the Body is enum.ContentTypeJson, the JSON value is converted to XML
// using mxj.NewMapJson and mapJson.XmlIndent functions.
// If the conversion is successful, the XML bytes are returned as a string.
// If the conversion fails, "<object></object>" is returned.
// For any other contentType, the value is formatted as a string with an XML tag wrapping it.
// The resulting XML string is returned.
func (b *Body) Xml() []byte {
	switch b.contentType {
	case enum.ContentTypeJson:
		mapJson, err := mxj.NewMapJson(b.Bytes())
		if helper.IsNil(err) {
			xmlBytes, err := mapJson.XmlIndent("", "  ", "object")
			if helper.IsNil(err) {
				return xmlBytes
			}
		}
		return []byte("<object></object>")
	default:
		return []byte(fmt.Sprintf("<string>%s</string>", b.value))
	}
}

// BytesByContentType returns the byte representation of the `Body` instance
// based on the provided `contentType`.
// If the `contentType` is `enum.ContentTypeJson`, it returns the result of `b.Json()`.
// If the `contentType` is `enum.ContentTypeXml`, it returns the result of `b.Xml()`.
// For any other `contentType`, it returns the result of `b.Bytes()`.
func (b *Body) BytesByContentType(contentType enum.ContentType) []byte {
	switch contentType {
	case enum.ContentTypeJson:
		return b.Json()
	case enum.ContentTypeXml:
		return b.Xml()
	default:
		return b.Bytes()
	}
}

// Bytes returns the byte representation of the `Body` instance by returning the byte array from the `bytes.Buffer` value.
func (b *Body) Bytes() []byte {
	return b.value.Bytes()
}

// String returns a string representation of the current Body instance.
// It utilizes the SimpleConvertToString function from the helper package to convert the value of the Body to a string.
// The resulting string representation of the Body is returned.
func (b *Body) String() string {
	return b.value.String()
}

// MarshalJSON marshals the value of the Body instance into JSON format.
// It uses helper.ConvertToBytes to convert the value to bytes and return it along with no errors.
func (b *Body) MarshalJSON() ([]byte, error) {
	return b.value.Bytes(), nil
}

// IsText returns a boolean value indicating whether the contentType of the Body is ContentTypeText.
func (b *Body) IsText() bool {
	return helper.Equals(b.contentType, enum.ContentTypeText)
}

// IsJson returns a boolean value indicating whether the contentType of the Body is ContentTypeJson.
func (b *Body) IsJson() bool {
	return helper.Equals(b.contentType, enum.ContentTypeJson)
}

// IsNotJson returns a boolean value indicating whether the contentType of the Body is not ContentTypeJson.
func (b *Body) IsNotJson() bool {
	return !b.IsJson()
}

// String returns the string representation of the CacheBodyValue instance.
// It calls the String method of the underlying bytes.Buffer type to get the string representation.
func (c *CacheBodyValue) String() string {
	return (*bytes.Buffer)(c).String()
}

// Bytes returns the byte slice representation of the CacheBodyValue instance.
// It calls the Bytes method of the underlying bytes.Buffer type to get the byte slice representation.
func (c *CacheBodyValue) Bytes() []byte {
	return (*bytes.Buffer)(c).Bytes()
}

// MarshalJSON returns the JSON encoding of the CacheBodyValue instance.
// The JSON encoding is obtained by calling the String method of the
// underlying bytes.Buffer type to retrieve the string representation,
// and then encoding it using json.Marshal.
// It returns a byte slice representing the JSON encoding and an error,
// if any occurred during the encoding process.
func (c *CacheBodyValue) MarshalJSON() ([]byte, error) {
	return c.Bytes(), nil
}

// UnmarshalJSON decodes the JSON data into a string and writes
// the string to the underlying bytes.Buffer type.
// It returns an error if there is an issue with decoding or writing
// the string to the buffer.
func (c *CacheBodyValue) UnmarshalJSON(data []byte) error {
	_, err := (*bytes.Buffer)(c).Write(data)
	return err
}

// mergeString concatenates `str` and `anotherStr` with a newline character,
// then converts the merged string to a `bytes.Buffer` using `helper.SimpleConvertToBuffer`
// and returns the converted buffer.
func mergeString(str, anotherStr string) *bytes.Buffer {
	merged := fmt.Sprintf("%s\n%s", str, anotherStr)
	return helper.SimpleConvertToBuffer(merged)
}

// mergeJSON merges the provided JSON strings.
// It iterates through each key-value pair in the second JSON string,
// and sets the value of the corresponding key in the first JSON string.
// If the key already exists in the first JSON string,
// the value is appended to the existing array of values under that key.
// If the key does not exist, it is added with the provided value.
// The merged JSON string is then converted to a *bytes.Buffer and returned.
//
// Parameters:
// - jsonStr: The first JSON string to merge.
// - anotherJsonStr: The second JSON string to merge.
//
// Returns:
// - A *bytes.Buffer containing the merged JSON string.
func mergeJSON(jsonStr, anotherJsonStr string) *bytes.Buffer {
	merged := jsonStr
	parsedJsonB := gjson.Parse(anotherJsonStr)
	parsedJsonB.ForEach(func(key, value gjson.Result) bool {
		merged = setJsonKeyValue(merged, key.String(), value.Value())
		return true // continue iterando
	})
	return helper.SimpleConvertToBuffer(merged)
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
