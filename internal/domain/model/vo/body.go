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
	"strings"
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

// newBodyFromIndexAndBackendResponse creates a new instance of Body based on the provided index and backendResponseVO.
// It constructs a default response body with initial fields: "ok" and "code" populated from backendResponseVO.
// The constructed body has a content type of ContentTypeJson and a value of the JSON representation of the initial fields.
// If backendResponseVO's body is nil, it returns the default body. Otherwise, it aggregates the body with the initial fields.
// The aggregation behavior depends on the value of groupResponse field in backendResponseVO.
// If groupResponse is true, it aggregates the body with the key generated from backendResponseVO's key method and body.
// Otherwise, it aggregates all the JSON fields of the body into the bodyHistory.
// It returns the constructed body.
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

	// caso o body do index seja nil retornamos apenas os campos padrões
	if helper.IsNil(backendResponseVO.Body()) {
		return body
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
// If the backendResponseVO's body is nil, it returns nil.
// If backendResponseVO's group field is false, it returns backendResponseVO's body.
// It constructs an empty body with contentType set to enum.ContentTypeJson and value "{}".
// It returns the aggregated body with the specified key from backendResponseVO and the constructed body.
func newBodyFromBackendResponse(backendResponseVO *backendResponse) *Body {
	// se for nil ja retornamos
	if helper.IsNil(backendResponseVO.Body()) {
		return nil
	}
	// verificamos se backendResponse quer ser agrupado com o campo extra-config
	if !backendResponseVO.Group() {
		return backendResponseVO.Body()
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
	if helper.IsEmpty(value.Bytes()) {
		return nil
	}
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
	// chamamos o merge para agregar os valores do outro body
	return b.merge(anotherBody)
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

// Add adds a new key-value pair to the Body instance.
// The key is a string and the value can be any type.
// If the contentType of the Body is ContentTypeText, the value will be converted to string
// using helper.SimpleConvertToString function and added to the Body.
// If the contentType is ContentTypeJson, the key-value pair will be added to the Body
// using the addJson method.
// If the contentType is neither ContentTypeText nor ContentTypeJson,
// the Body instance will not be modified, and it will be returned as is.
// The method returns the updated Body instance and an error if any.
func (b *Body) Add(key string, value any) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.addString(helper.SimpleConvertToString(value)), nil
	case enum.ContentTypeJson:
		return b.addJson(key, value)
	default:
		return b, nil
	}
}

// Append appends a key-value pair to the Body instance.
// The key-value pair is appended according to the contentType of the Body:
//   - If the contentType is ContentTypeText, the value is converted to a string
//     and appended as the value of the provided key.
//   - If the contentType is ContentTypeJson, the key-value pair is appended as a
//     JSON object to the Body.
//   - For other contentTypes, the Body instance is returned unchanged.
//
// An error is returned if there is any issue with appending the key-value pair.
func (b *Body) Append(key string, value any) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.appendString(helper.SimpleConvertToString(value)), nil
	case enum.ContentTypeJson:
		return b.appendJson(key, value)
	default:
		return b, nil
	}
}

// Set returns a new instance of Body with the provided key-value pair.
// The new Body instance will have the same contentType as the original Body instance.
// If the contentType is ContentTypeText, Set converts the value to a string and calls setString.
// If the contentType is ContentTypeJson, Set calls setJson.
// If the contentType is neither ContentTypeText nor ContentTypeJson, Set returns the original Body instance.
// If an error occurs during the set operation, Set returns an error.
func (b *Body) Set(key string, value any) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.setString(helper.SimpleConvertToString(value)), nil
	case enum.ContentTypeJson:
		return b.setJson(key, value)
	default:
		return b, nil
	}
}

// Replace is a method that replaces the value of a specified key in the Body instance.
// The new value is specified by the provided 'value' argument.
// If the content type of the Body is ContentTypeText, the value is converted to a string before replacing it.
// If the content type is ContentTypeJson, the value is replaced in the JSON object.
// If the content type is not ContentTypeText or ContentTypeJson, the same Body instance is returned.
// The method returns the modified Body instance and an error if any occurred during the replacement process.
func (b *Body) Replace(key string, value any) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.replaceString(key, helper.SimpleConvertToString(value)), nil
	case enum.ContentTypeJson:
		return b.replaceJson(key, value)
	default:
		return b, nil
	}
}

// Rename renames the value associated with the given key in the Body.
// If the Body has ContentTypeText, it replaces the string value.
// If the Body has ContentTypeJson, it renames the key in the JSON data.
// For any other content type, it returns the original Body.
// The method returns the modified Body and any error encountered during renaming.
func (b *Body) Rename(key string, value any) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.replaceString(key, helper.SimpleConvertToString(value)), nil
	case enum.ContentTypeJson:
		return b.renameJson(key, value)
	default:
		return b, nil
	}
}

// Delete removes the value associated with the specified key from the Body instance.
// If the content type of the Body is ContentTypeText, the resulting Body instance will have the value of the specified
// key replaced by an empty string.
// If the content type of the Body is ContentTypeJson, the resulting Body instance will have the value associated with
// the specified key removed.
// If the content type of the Body is not ContentTypeText or ContentTypeJson, the same Body instance will be returned.
// An error will be returned if there are any issues during the deletion process.
func (b *Body) Delete(key string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.replaceString(key, ""), nil
	case enum.ContentTypeJson:
		return b.deleteJson(key)
	default:
		return b, nil
	}
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

// addString appends the provided value to the existing value of the Body instance.
// The modified value is stored in a new Body instance, while the original Body instance remains unchanged.
// If the value is empty, the method returns the original Body instance without any modifications.
// The new Body instance will have the same contentType as the original Body instance.
// The method returns a pointer to the updated Body instance.
func (b *Body) addString(value string) *Body {
	if helper.IsEmpty(value) {
		return b
	}

	modifiedValue := b.String() + value
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}
}

// addJson adds a key-value pair to the Body instance if the key is not empty.
// It checks if the key already exists in the JSON-encoded Body string.
// If the key exists, it modifies the value by appending the new value to the existing value.
// If the key does not exist, it adds a new key-value pair to the JSON-encoded Body string.
// It returns a new Body instance with the modified value and the contentType of the original Body instance.
// If any error occurs during the modification process, it returns nil and the error.
func (b *Body) addJson(key string, value any) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	bodyStr := b.String()

	var modifiedValue string
	var err error

	result := gjson.Get(bodyStr, key)
	if result.Exists() {
		modifiedValue, err = sjson.Set(bodyStr, key, []any{result.Value(), value})
	} else {
		modifiedValue, err = sjson.Set(bodyStr, key, value)
	}

	if helper.IsNotNil(err) {
		return nil, err
	}
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// appendString appends the provided string value to the content of the Body instance.
// If the current content is empty or the provided value is not empty, it returns the same Body instance.
// Otherwise, it returns a new instance of Body with the concatenated content of the original Body instance and the provided value.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) appendString(value string) *Body {
	bodyStr := b.String()
	if helper.IsEmpty(bodyStr) || helper.IsNotEmpty(value) {
		return b
	}
	modifiedValue := b.String() + value
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}
}

// appendJson appends a key-value pair to the existing JSON content in the Body instance.
// If the provided key is empty, the original Body instance is returned.
// If the provided key does not exist in the JSON content, the original Body instance is returned.
// If the key exists, the value is appended to the existing value for that key.
// Returns a new instance of Body.
// The new Body instance will have the same contentType as the original Body instance.
// If an error occurs during modification, a nil Body instance and the error are returned.
func (b *Body) appendJson(key string, value any) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	bodyStr := b.String()

	result := gjson.Get(bodyStr, key)
	if !result.Exists() {
		return b, nil
	}

	modifiedValue, err := sjson.Set(bodyStr, key, []any{result.Value(), value})
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// setString returns a new instance of Body with the provided string value.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) setString(value string) *Body {
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(value),
	}
}

// setJson sets the specified key to the provided value in the JSON body.
// If the key is empty, it returns the original Body instance.
// The modified JSON body will have the same contentType as the original Body instance.
// If there is an error while setting the key, it returns nil and the error.
// The value parameter can be of any type, it will be converted to a string representation before setting.
// If the value is nil, the entire body will be set to nil.
// The setJson method is used in the Set method, depending on the content type of the Body instance.
// The value can be set into a JSON body or a text body.
// For a text body, the entire body will be set to the converted value.
// If the content type is neither JSON nor text, it returns the original Body instance.
func (b *Body) setJson(key string, value any) (*Body, error) {
	if helper.IsEmpty(key) {
		// todo: futuramente podemos setar o body inteiro de value caso esteja nulo, isso seria
		//  legal pois, poderia mudar totalmente o body a partir de um valor
		return b, nil
	}

	modifiedValue, err := sjson.Set(b.String(), key, value)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// replaceString replaces all occurrences of the key in the Body's string representation
// with the provided value. It returns a new instance of Body with the modified value.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) replaceString(key, value string) *Body {
	if helper.IsEmpty(key) {
		// todo: imprimir log de atenção?
		return b
	}

	modifiedValue := strings.ReplaceAll(b.String(), key, value)
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}
}

// replaceJson replaces the value of the specified key in the JSON body with the provided value.
// If the key is empty or the key does not exist in the JSON body, it returns the original body instance.
// If an error occurs while modifying the JSON body, it returns nil and the error.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) replaceJson(key string, value any) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	result := gjson.Get(b.String(), key)
	if !result.Exists() {
		return b, nil
	}

	modifiedValue, err := sjson.Set(b.String(), key, value)
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// renameJson renames a key in the JSON body.
// The oldKey parameter specifies the key to be renamed.
// The newKey parameter is the new name for the key.
// If the oldKey is empty, the method returns the original Body instance.
// If the oldKey does not exist in the JSON body, the method returns the original Body instance.
// The method modifies the JSON body by deleting the oldKey and setting the newKey with the same value.
// The method returns a new instance of Body with the modified JSON body and the same contentType as the original Body instance.
// If there is an error during the modification of the JSON body, the method returns nil and the error.
func (b *Body) renameJson(oldKey string, newKey any) (*Body, error) {
	if helper.IsEmpty(oldKey) {
		return b, nil
	}

	newKeyStr := helper.SimpleConvertToString(newKey)
	bodyStr := b.String()

	result := gjson.Get(bodyStr, oldKey)
	if !result.Exists() {
		return b, nil
	}

	modifiedValue, err := sjson.Delete(bodyStr, oldKey)
	if helper.IsNotNil(err) {
		return nil, err
	}

	modifiedValue, err = sjson.Set(modifiedValue, newKeyStr, result.Value())
	if helper.IsNotNil(err) {
		return nil, err
	}
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// deleteJson deletes the specified key from the Body's JSON value.
// If the key is empty, the method returns the original Body instance and a nil error.
// If the deletion is successful, the method returns a new Body instance with the modified JSON value.
// The new Body instance will have the same contentType as the original Body instance.
// If an error occurs during the deletion process, the method returns a nil Body instance and the error.
func (b *Body) deleteJson(key string) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	modifiedValue, err := sjson.Delete(b.String(), key)
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(modifiedValue),
	}, nil
}

// merge merges the values of two Body instances into a new Body instance.
// If anotherBody is nil, it returns the current Body instance.
// It checks the contentType of the Body instance and performs the merge operation accordingly.
// If the contentType is ContentTypeJson, it uses the mergeJSON function to merge the JSON values of the bodies.
// If the contentType is ContentTypeText, it uses the mergeString function to merge the string values of the bodies.
// The merged Body value is then used to create a new Body instance with the same contentType,
// which is returned as the result.
// If the contentType is neither ContentTypeJson nor ContentTypeText, it returns the current Body instance unchanged.
func (b *Body) merge(anotherBody *Body) *Body {
	if helper.IsNil(anotherBody) {
		return b
	}
	switch b.contentType {
	case enum.ContentTypeJson:
		return b.mergeJSON(anotherBody.String())
	case enum.ContentTypeText:
		return b.mergeString(anotherBody.String())
	}
	return b
}

// mergeString merges the provided string with the current Body instance and returns a new
// instance of Body. The merged string is formed by concatenating the string representation
// of the current Body instance and the provided string, separated by a newline character.
// The new Body instance will have the same contentType as the original Body instance.
func (b *Body) mergeString(str string) *Body {
	merged := fmt.Sprintf("%s\n%s", b.String(), str)
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(merged),
	}
}

// mergeJSON merges the provided JSON string into the existing JSON string of the Body instance.
// The value of each key in the provided JSON string is set in the existing JSON string.
// If the key already exists in the existing JSON string, the value is appended to the existing array of values under that key.
// If the key does not exist, it is added with the provided value.
// The updated JSON string is returned in a new Body instance with the same contentType as the original Body instance.
func (b *Body) mergeJSON(jsonStr string) *Body {
	merged := b.String()
	parsedJsonB := gjson.Parse(jsonStr)
	parsedJsonB.ForEach(func(key, value gjson.Result) bool {
		merged = setJsonKeyValue(merged, key.String(), value.Value())
		return true // continue iterando
	})
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(merged),
	}
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
