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
	"compress/gzip"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/clbanning/mxj/v2"
	"github.com/iancoleman/strcase"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"strings"
)

// Body represents the content and format of an HTTP httpRequest or httpResponse body.
type Body struct {
	// contentType is an enumeration type that represents the format of the content.
	// It can have the following values: ContentTypeText, ContentTypeJson.
	contentType enum.ContentType
	// value represents the content of an HTTP httpRequest or httpResponse body. It is stored as a bytes.Buffer object.
	value *bytes.Buffer
}

func NewBodyByContentType(contentType string, buffer *bytes.Buffer) *Body {
	return NewBody(contentType, "", buffer)
}

func NewBody(contentType, contentEncoding string, buffer *bytes.Buffer) *Body {
	// se vazio, retornamos vazio
	if helper.IsEmpty(buffer.Bytes()) {
		return nil
	}

	// instanciamos o content type como enum com base na string
	contentTypeEnum := enum.ContentTypeFromString(contentType)
	// instanciamos o content encoding como enum com base na string
	contentEncodingEnum := enum.ContentEncodingFromString(contentEncoding)

	// verificamos se tem encode valido, se tiver, chamamos a func para trabalhar em cima do encode
	if contentEncodingEnum.IsEnumValid() {
		return newBodyByContentEncoding(contentTypeEnum, contentEncodingEnum, buffer)
	}

	// montamos o body com o valor original
	return &Body{
		contentType: contentTypeEnum,
		value:       buffer,
	}
}

func newBodyByContentEncoding(contentType enum.ContentType, contentEncoding enum.ContentEncoding, buffer *bytes.Buffer,
) *Body {
	// verificamos se o encoding é valido, caso não seja apenas retornamos sem trabalhar em cima do encode
	if !contentEncoding.IsEnumValid() {
		return &Body{
			contentType: contentType,
			value:       buffer,
		}
	}
	// criamos o reader gzip a partir do buffer recebido
	reader, err := gzip.NewReader(buffer)
	if helper.IsNotNil(err) {
		logger.Warning("Error creating gzip reader for body:", err)
		return nil
	}
	defer reader.Close()
	// lemos o reader gerando os bytes
	unzipBytes, err := io.ReadAll(reader)
	if helper.IsNotNil(err) {
		logger.Warning("Error read gzip bytes body:", err)
		return nil
	}
	// montamos o body com o valor descompactado
	return &Body{
		contentType: contentType,
		value:       bytes.NewBuffer(unzipBytes),
	}
}

func newBodyByString(s string) *Body {
	if helper.IsEmpty(s) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeText,
		value:       helper.SimpleConvertToBuffer(s),
	}
}

func newBodyByJson(a any) *Body {
	if helper.IsNil(a) || helper.IsEmpty(a) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(a),
	}
}

func newBodyByError(path string, err error) *Body {
	// construímos o errorBody a partir do erro e path
	errBody := newErrorBody(path, err)
	if helper.IsNil(errBody) {
		return nil
	}
	// construímos o body com esse objeto
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(errBody),
	}
}

// newBodyByHttpBackendResponse creates a new instance of Body based on the provided index and backendResponseVO.
// It constructs a default httpResponse body with initial fields: "ok" and "code" populated from backendResponseVO.
// The constructed body has a content type of ContentTypeJson and a value of the JSON representation of the initial fields.
// If backendResponseVO's body is nil, it returns the default body. Otherwise, it aggregates the body with the initial fields.
// The aggregation behavior depends on the value of groupResponse field in backendResponseVO.
// If groupResponse is true, it aggregates the body with the key generated from backendResponseVO's key method and body.
// Otherwise, it aggregates all the JSON fields of the body into the bodyHistory.
// It returns the constructed body.
func newBodyByHttpBackendResponse(index int, backendResponseVO *httpBackendResponse) *Body {
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
	if backendResponseVO.GroupByType() {
		body = body.AggregateByKey(backendResponseVO.Key(index), backendResponseVO.Body())
	} else {
		body = body.Aggregate(backendResponseVO.Body())
	}
	// retornamos o body
	return body
}

func newBodyAggregateByKey(key string, anotherBody *Body) *Body {
	// se o body for nil retornamos nil
	if helper.IsNil(anotherBody) {
		return nil
	}

	// construímos o body vazio
	body := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer("{}"),
	}
	// agregamos na chave indicada o valor do outro body
	return body.AggregateByKey(key, anotherBody)
}

// newBodyBySlice creates a new instance of Body based on the provided slice of *Body.
// If the slice is empty, it returns nil.
// The contentType of the new Body instance is set to ContentTypeJson.
// The value of the new Body instance is obtained by converting the slice to a bytes.Buffer object using helper.SimpleConvertToBuffer.
// It returns a pointer to the constructed Body instance.
func newBodyBySlice(sliceOfBodies []*Body) *Body {
	if helper.IsEmpty(sliceOfBodies) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(sliceOfBodies),
	}
}

// newBodyByCache creates a new instance of Body based on the provided CacheBody.
// If the cacheBody parameter is nil, it returns nil.
// Otherwise, it sets the contentType field of the new Body instance to the ContentType field of the cacheBody parameter.
// It converts the Value field of the cacheBody parameter to a bytes.Buffer pointer and assigns it to the value field of the new Body instance.
// It returns a pointer to the newly created Body instance.
func newBodyByCache(cacheBodyVO *CacheBody) *Body {
	if helper.IsNil(cacheBodyVO) {
		return nil
	}
	return &Body{
		contentType: cacheBodyVO.ContentType,
		value:       (*bytes.Buffer)(cacheBodyVO.Value),
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

	value := anotherBody.Raw()
	mergedBodyRaw := setJsonKeyValue(b.Raw(), key, value)
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mergedBodyRaw),
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
// If the `contentType` is `enum.ContentTypeJson`, it returns the result of `b.Map()`.
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

func (b *Body) Raw() string {
	return parseStringValueToRaw(b.String())
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
func (b *Body) Add(key string, value string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.addString(value), nil
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
func (b *Body) Append(key string, value string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.appendString(value), nil
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
func (b *Body) Set(key string, value string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.setString(value), nil
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
func (b *Body) Replace(key string, value string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.replaceString(key, value), nil
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
func (b *Body) Rename(key string, value string) (*Body, error) {
	switch b.contentType {
	case enum.ContentTypeText:
		return b.replaceString(key, value), nil
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

func (b *Body) Map(mapper *Mapper) *Body {
	// se o mapper for vazio, retornamos o body atual
	if mapper.IsEmpty() {
		return b
	}
	// mapeamos o valor do body com base no tipo de conteúdo
	switch b.contentType {
	case enum.ContentTypeJson:
		return b.mapJson(mapper)
	case enum.ContentTypeText:
		return b.mapText(mapper)
	default:
		return b
	}
}

func (b *Body) Projection(projectionVO *Projection) *Body {
	// se não for do tipo json, ou a projeção for nil, ou vazia, retornamos o body atual, sem modificações
	if helper.IsNotEqualTo(b.contentType, enum.ContentTypeJson) || helper.IsNil(projectionVO) || projectionVO.IsEmpty() {
		return b
	}
	// criamos um novo body com o json projetado a partir do objeto de valor
	return b.projectionJson(projectionVO)
}

// OmitEmpty returns a new instance of Body with empty values omitted.
// The new Body instance will have the same contentType as the original Body instance.
// If the contentType is ContentTypeJson, omitEmptyJson() will be called.
// If the contentType is ContentTypeText, omitEmptyText() will be called.
// If the contentType is neither ContentTypeJson nor ContentTypeText, the original Body instance will be returned.
func (b *Body) OmitEmpty() *Body {
	switch b.contentType {
	case enum.ContentTypeJson:
		return b.omitEmptyJson()
	case enum.ContentTypeText:
		return b.omitEmptyText()
	default:
		return b
	}
}

func (b *Body) ToCase(nomenclature enum.Nomenclature) *Body {
	if helper.IsNotEqualTo(b.contentType, enum.ContentTypeJson) {
		return b
	}
	jsonStr := convertKeysToCase(b.String(), nomenclature)
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(jsonStr),
	}
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
func (b *Body) addJson(key string, value string) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	bodyRaw := b.Raw()

	var modifiedValue string
	var err error

	result := gjson.Get(bodyRaw, key)
	if result.Exists() && helper.IsNotEqualTo(result.Type, gjson.Null) {
		modifiedValue, err = sjson.SetRaw(bodyRaw, key, aggregateJsonValue(result, value))
	} else {
		modifiedValue, err = sjson.SetRaw(bodyRaw, key, parseStringValueToRaw(value))
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
func (b *Body) appendJson(key string, value string) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	bodyRaw := b.Raw()

	result := gjson.Get(bodyRaw, key)
	if !result.Exists() || helper.Equals(result.Type, gjson.Null) {
		return b, nil
	}

	modifiedValue, err := sjson.SetRaw(bodyRaw, key, aggregateJsonValue(result, value))
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
func (b *Body) setJson(key string, value string) (*Body, error) {
	if helper.IsEmpty(key) {
		// todo: futuramente podemos setar o body inteiro de value caso esteja nulo, isso seria
		//  legal pois, poderia mudar totalmente o body a partir de um valor
		return b, nil
	}

	modifiedValue, err := sjson.SetRaw(b.Raw(), key, parseStringValueToRaw(value))
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
		logger.Warning("Replace ignored as the modifier key is empty!")
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
func (b *Body) replaceJson(key string, value string) (*Body, error) {
	if helper.IsEmpty(key) {
		return b, nil
	}

	result := gjson.Get(b.Raw(), key)
	if !result.Exists() || helper.Equals(result.Type, gjson.Null) {
		return b, nil
	}

	modifiedValue, err := sjson.SetRaw(b.Raw(), key, parseStringValueToRaw(value))
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
func (b *Body) renameJson(oldKey, newKey string) (*Body, error) {
	if helper.IsEmpty(oldKey) {
		return b, nil
	}

	bodyRaw := b.Raw()

	result := gjson.Get(bodyRaw, oldKey)
	if !result.Exists() {
		return b, nil
	}

	modifiedValue, err := sjson.Delete(bodyRaw, oldKey)
	if helper.IsNotNil(err) {
		return nil, err
	}

	modifiedValue, err = sjson.SetRaw(modifiedValue, newKey, parseValueToRaw(result))
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
		return b.mergeJSON(anotherBody.Raw())
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

func (b *Body) mergeJSON(jsonStr string) *Body {
	merged := b.Raw()
	parsedJsonStr := gjson.Parse(jsonStr)
	parsedJsonStr.ForEach(func(key, value gjson.Result) bool {
		merged = setJsonKeyValue(merged, key.String(), parseValueToRaw(value))
		return true
	})
	return &Body{
		contentType: b.ContentType(),
		value:       helper.SimpleConvertToBuffer(merged),
	}
}

func (b *Body) mapJson(mapper *Mapper) *Body {
	// damos o parse do json do body
	parsedJson := gjson.Parse(b.String())
	// se for array chamamos o mapJsonArray
	if parsedJson.IsArray() {
		return b.mapJsonArray(mapper, parsedJson)
	}
	// se for um objeto chamamos o mapJsonObject
	return b.mapJsonObject(mapper, parsedJson)
}

func (b *Body) mapJsonArray(mapper *Mapper, jsonArray gjson.Result) *Body {
	// iniciamos o json array vazio
	mappedArray := "[]"
	// iteramos o json array atual para mapear json caso seja um array de objeto
	jsonArray.ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			// caso ele seja objeto, chamamos o mapJsonObject para mapear o json objeto do index atual
			projectedObject := b.mapJsonObject(mapper, value)
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", projectedObject.Raw())
		} else if value.IsArray() {
			// caso ele seja array, chamamos o mesmo novamente para iterar o sub array
			projectedSubArray := b.mapJsonArray(mapper, value)
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", projectedSubArray.Raw())
		} else {
			// caso ele não seja json, não tem nada para mapear, apenas adicionamos o valor
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", parseValueToRaw(value))
		}
		return true
	})
	// retornamos o mappedArray como valor de um novo Body
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedArray),
	}
}

func (b *Body) mapJsonObject(mapper *Mapper, jsonObject gjson.Result) *Body {
	// instanciamos o json mapeado com os valores atuais
	mappedJson := jsonObject.String()
	// iteramos o mapper para renomear os campos
	for _, key := range mapper.Keys() {
		// instanciamos a nova chave a partir da antiga
		newKey := mapper.Get(key)
		// se a chave for igual a chave atual, ignoramos
		if helper.Equals(key, newKey) {
			continue
		}
		// obtemos o valor do json pela chave antiga
		jsonValue := jsonObject.Get(key)
		// caso ele exista, inserimos ele na nova chave
		if jsonValue.Exists() {
			// inserimos o valor na nova chave
			newMappedJson, err := sjson.SetRaw(mappedJson, newKey, parseValueToRaw(jsonValue))
			if helper.IsNil(err) {
				// caso tenha dado certo, removemos a chave antiga
				mappedJson, _ = sjson.Delete(newMappedJson, key)
			} else {
				// caso tenha dado errado, adicionamos a chave antiga no novo mapped json
				mappedJson, _ = sjson.SetRaw(mappedJson, key, parseValueToRaw(jsonValue))
			}
		}
	}
	// retornamos o body com o json mapeado
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedJson),
	}
}

func (b *Body) mapText(mapper *Mapper) *Body {
	// instanciamos o texto mapeado com o valor atual
	mappedText := b.String()
	// iteramos o mapper para renomear os campos
	for _, key := range mapper.Keys() {
		// instanciamos a nova chave a partir da antiga
		newKey := mapper.Get(key)
		// se a chave for igual a chave atual, ignoramos
		if helper.Equals(key, newKey) {
			continue
		}
		// damos o replace da chave antiga para chave atual
		mappedText = strings.ReplaceAll(mappedText, key, newKey)
	}
	// retornamos o texto mapeado
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedText),
	}
}

func (b *Body) projectionJson(projectionVO *Projection) *Body {
	// damos o parse do json do body
	parsedJson := gjson.Parse(b.String())
	// se for array chamamos o projectionJsonArray
	if parsedJson.IsArray() {
		return b.projectionJsonArray(projectionVO, parsedJson)
	}
	// se for um objeto chamamos o projectionJsonObject
	return b.projectionJsonObject(projectionVO, parsedJson)
}

func (b *Body) projectionJsonObject(projectionVO *Projection, jsonObject gjson.Result) *Body {
	// verificamos se o tipo de projeção é Rejection, se for executamos as regras
	if helper.Equals(projectionVO.Type(), enum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonObject(projectionVO, jsonObject)
	}
	// se não for do tipo rejeição, então é adição ou todos, os dois aplicam as mesmas regras
	return b.projectionAdditionJsonObject(projectionVO, jsonObject)
}

func (b *Body) projectionAdditionJsonObject(projectionVO *Projection, jsonObject gjson.Result) *Body {
	// iniciamos o json vazio
	projectedJson := "{}"
	// iteramos o projection vo
	for _, key := range projectionVO.Keys() {
		// se for rejection pulamos e não adicionamos no novo json
		if projectionVO.IsRejection(key) {
			continue
		}
		// obtemos o valor da chave no json
		jsonValue := jsonObject.Get(key)
		// caso ele exista, adicionamos
		if jsonValue.Exists() {
			projectedJson, _ = sjson.SetRaw(projectedJson, key, parseValueToRaw(jsonValue))
		}
	}
	// retornamos o body com o novo json
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectedJson),
	}
}

func (b *Body) projectionRejectionJsonObject(projectionVO *Projection, jsonObject gjson.Result) *Body {
	// iniciamos o json com o valor atual
	projectionJson := jsonObject.String()
	// iteramos o projection vo
	for _, key := range projectionVO.Keys() {
		// removemos o campo com base na chave de projeção
		projectionJson, _ = sjson.Delete(projectionJson, key)
	}
	// retornamos o body com o novo json
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectionJson),
	}
}

func (b *Body) projectionJsonArray(projectionVO *Projection, jsonArray gjson.Result) *Body {
	// iniciamos o json array vazio
	projectedArray := "[]"
	// iteramos o json array atual para projetar json caso seja um array de objeto
	jsonArray.ForEach(func(key, value gjson.Result) bool {
		// processamos o index do array
		projectedArray = b.projectionJsonArrayCurrentIndex(projectionVO, projectedArray, value)
		// retorna true para continuar iterando
		return true
	})
	// se ele quer filtrar por index
	projectedArray = b.projectionJsonArrayNumericKeys(projectionVO, projectedArray)
	// retornamos o body com o resultado da projeção
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectedArray),
	}
}

func (b *Body) projectionJsonArrayCurrentIndex(projectionVO *Projection, projectedArray string, value gjson.Result,
) string {
	if value.IsObject() {
		// caso ele seja objeto, chamamos o projectionJsonObject para projetar o json objeto do index atual
		projectedObject := b.projectionJsonObject(projectionVO, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedObject.Raw())
	} else if value.IsArray() {
		// caso ele seja array, chamamos o mesmo novamente para iterar o array filho
		projectedSubArray := b.projectionJsonArray(projectionVO, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedSubArray.Raw())
	} else {
		// caso ele não seja json, não tem nada para projetar, apenas adicionamos o valor
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(value))
	}
	// retornamos o array com o index inserido
	return projectedArray
}

func (b *Body) projectionJsonArrayNumericKeys(projectionVO *Projection, projectedJson string) string {
	// se ele não contém ao menos uma chave de projeção numérica retornamos o json já projetado
	if projectionVO.NotContainsNumericKey() {
		return projectedJson
	}
	// verificamos se o tipo da projeção numérica é rejeição, para seguir com as regras
	if helper.Equals(projectionVO.TypeNumeric(), enum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonArray(projectionVO, projectedJson)
	}
	// caso não seja rejection, ou ele é all, ou addition, usamos as mesmas regras
	return b.projectionAdditionJsonArray(projectionVO, projectedJson)
}

func (b *Body) projectionAdditionJsonArray(projectionVO *Projection, projectedJson string) string {
	// transformamos o projectedJson em um gjson.Result
	parsedProjectedJson := gjson.Parse(projectedJson)
	// instanciamos a projeção de array zerada
	projectedArray := "[]"
	// iteramos a projeção para projetar os index mencionados
	for _, key := range projectionVO.Keys() {
		if helper.IsNumeric(key) && projectionVO.IsAddition(key) {
			jsonValue := parsedProjectedJson.Get(key)
			if jsonValue.Exists() {
				projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(jsonValue))
			}
		}
	}
	// retornamos o array projetado
	return projectedArray
}

func (b *Body) projectionRejectionJsonArray(projectionVO *Projection, projectedJson string) string {
	// transformamos o projectedJson em um gjson.Result
	parsedProjectedJson := gjson.Parse(projectedJson)
	// instanciamos a projeção de array vazia
	projectedArray := "[]"
	// iteramos a lista
	parsedProjectedJson.ForEach(func(key, value gjson.Result) bool {
		if helper.NotContains(projectionVO.Keys(), key.String()) {
			projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(value))
		}
		return true
	})
	// retornamos o array projetado
	return projectedArray
}

// omitEmptyJson returns a new instance of Body with all empty fields removed from the JSON string.
// It uses the removeAllEmptyFields function to remove empty fields recursively.
// The new Body instance will have the same contentType as the original Body instance.
// The original JSON string is modified in-place to remove the empty fields.
func (b *Body) omitEmptyJson() *Body {
	jsonStr := removeAllEmptyFields(b.Raw())
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(jsonStr),
	}
}

// omitEmptyText returns a new instance of Body with empty text omitted.
// The new Body instance will have the same contentType as the original Body instance.
//
// If the contentType is ContentTypeJson, omitEmptyJson() will be called.
// If the contentType is ContentTypeText, omitEmptyText() will be called.
// If the contentType is neither ContentTypeJson nor ContentTypeText, the original Body instance will be returned.
func (b *Body) omitEmptyText() *Body {
	s := helper.CleanAllRepeatSpaces(b.String())
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(s),
	}
}

func setJsonKeyValue(jsonStr, key string, value string) string {
	result := gjson.Get(jsonStr, key)

	// Se a key já existe no JSON A
	if result.Exists() && helper.IsNotEqualTo(result.Type, gjson.Null) {
		jsonStr, _ = sjson.SetRaw(jsonStr, key, aggregateJsonValue(result, value))
	} else {
		jsonStr, _ = sjson.SetRaw(jsonStr, key, parseStringValueToRaw(value))
	}

	return jsonStr
}

func parseValueToRaw(value gjson.Result) string {
	if helper.Equals(value.Type, gjson.Null) {
		return "null"
	}
	return value.Raw
}

func parseStringValueToRaw(value string) string {
	parse := gjson.Parse(value)
	if helper.Equals(parse.Type, gjson.Null) {
		return "null"
	}
	return parse.Raw
}

func aggregateJsonValue(value gjson.Result, newValue string) string {
	var newArray []gjson.Result

	if value.IsArray() {
		newArray = value.Array()
	} else {
		newArray = []gjson.Result{value}
	}

	newParsedValue := gjson.Parse(newValue)

	if newParsedValue.IsArray() {
		newArray = append(newArray, newParsedValue.Array()...)
	} else {
		newArray = append(newArray, newParsedValue)
	}

	newArrayJson := "["
	for i, v := range newArray {
		if helper.Equals(v.Type, gjson.Null) || helper.IsEmpty(v.String()) {
			continue
		}
		if i != 0 {
			newArrayJson += ","
		}
		newArrayJson += parseValueToRaw(v)
	}
	newArrayJson += "]"
	return newArrayJson
}

// removeAllEmptyFields removes all empty fields from a JSON string recursively.
// It iterates over each line in the JSON string using gjson.ForEachLine, and for each line,
// it iterates over each key-value pair using line.ForEach.
// If the value is empty, it deletes the corresponding key-value pair using sjson.Delete.
// If the value is an object or an array, it recursively calls removeAllEmptyFields on that value.
// Note: The input JSON string is modified in-place.
func removeAllEmptyFields(jsonStr string) string {
	gjson.Parse(jsonStr).ForEach(func(key, value gjson.Result) bool {
		// se for objeto, chamamos novamente este método passando o value
		if value.IsObject() || value.IsArray() {
			subJsonStr := removeAllEmptyFields(parseValueToRaw(value))
			value = gjson.Parse(subJsonStr)
		}
		// verificamos se o valor esta vazio
		if helper.IsEmpty(value.Value()) {
			// caso esteja vazio, removemos a chave
			jsonStr, _ = sjson.Delete(jsonStr, key.String())
		} else {
			// caso não esteja vazio, inserimos o valor possivelmente modificado
			jsonStr, _ = sjson.SetRaw(jsonStr, key.String(), parseValueToRaw(value))
		}
		return true
	})
	return jsonStr
}

func convertKeysToCase(jsonStr string, nomenclature enum.Nomenclature) string {
	parsedJson := gjson.Parse(jsonStr)

	jsonStrCase := "{}"
	if parsedJson.IsArray() {
		jsonStrCase = "[]"
	}

	parsedJson.ForEach(func(key, value gjson.Result) bool {
		newKey := key.String()
		switch nomenclature {
		case enum.NomenclatureCamel:
			newKey = strcase.ToCamel(newKey)
		case enum.NomenclatureLowerCamel:
			newKey = strcase.ToLowerCamel(newKey)
		case enum.NomenclatureSnake:
			newKey = strcase.ToSnake(newKey)
		case enum.NomenclatureScreamingSnake:
			newKey = strcase.ToScreamingSnake(newKey)
		case enum.NomenclatureKebab:
			newKey = strcase.ToKebab(newKey)
		case enum.NomenclatureScreamingKebab:
			newKey = strcase.ToScreamingKebab(newKey)
		}

		if value.IsObject() || value.IsArray() {
			subJsonStr := convertKeysToCase(parseValueToRaw(value), nomenclature)
			jsonStrCase, _ = sjson.SetRaw(jsonStrCase, newKey, subJsonStr)
		} else {
			jsonStrCase, _ = sjson.SetRaw(jsonStrCase, newKey, parseValueToRaw(value))
		}

		return true
	})

	return jsonStrCase
}
