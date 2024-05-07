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
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/enum"
	"github.com/clbanning/mxj/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"strings"
)

// Body represents the content and format of an HTTP request or response body.
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
		return NewBodyByContentEncoding(contentTypeEnum, contentEncodingEnum, buffer)
	}

	// montamos o body com o valor original
	return &Body{
		contentType: contentTypeEnum,
		value:       buffer,
	}
}

func NewBodyByContentEncoding(contentType enum.ContentType, contentEncoding enum.ContentEncoding, buffer *bytes.Buffer,
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

func NewBodyByString(s string) *Body {
	if helper.IsEmpty(s) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeText,
		value:       helper.SimpleConvertToBuffer(s),
	}
}

func NewBodyByJson(a any) *Body {
	if helper.IsNil(a) || helper.IsEmpty(a) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(a),
	}
}

func NewBodyByError(path string, err error) *Body {
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

func NewBodyByHttpBackendResponse(index int, httpBackendResponse *HttpBackendResponse) *Body {
	// construímos o body padrão de resposta, com os campos iniciais
	bodyJson := "{}"
	bodyJson, _ = sjson.Set(bodyJson, "ok", httpBackendResponse.Ok())
	bodyJson, _ = sjson.Set(bodyJson, "code", httpBackendResponse.StatusCode())

	// construímos o body com os valores padrões
	body := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(bodyJson),
	}

	// caso o body do index seja nil retornamos apenas os campos padrões
	if helper.IsNil(httpBackendResponse.Body()) {
		return body
	}

	// caso seja string ou slice agregamos na chave, caso contrario, iremos agregar todos os campos json no bodyHistory
	if httpBackendResponse.GroupByType() {
		body = body.AggregateByKey(httpBackendResponse.Key(index), httpBackendResponse.Body())
	} else {
		body = body.Aggregate(httpBackendResponse.Body())
	}
	// retornamos o body
	return body
}

func NewBodyAggregateByKey(key string, anotherBody *Body) *Body {
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

func NewBodyBySlice(sliceOfBodies []*Body) *Body {
	if helper.IsEmpty(sliceOfBodies) {
		return nil
	}
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer(sliceOfBodies),
	}
}

func NewBodyByCache(cacheBody *CacheBody) *Body {
	if helper.IsNil(cacheBody) {
		return nil
	}
	return &Body{
		contentType: cacheBody.ContentType,
		value:       (*bytes.Buffer)(cacheBody.Value),
	}
}

func NewBodyJson() *Body {
	return &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer("{}"),
	}
}

func (b *Body) ContentType() enum.ContentType {
	return b.contentType
}

func (b *Body) Value() *bytes.Buffer {
	return b.value
}

func (b *Body) Aggregate(anotherBody *Body) *Body {
	if helper.IsNil(anotherBody) {
		return b
	}
	// chamamos o merge para agregar os valores do outro body
	return b.merge(anotherBody)
}

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

func (b *Body) Interface() any {
	switch b.contentType {
	case enum.ContentTypeJson:
		return gjson.ParseBytes(b.value.Bytes()).Value()
	}
	return b.value.String()
}

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

func (b *Body) CompactString() string {
	return helper.CompactString(b.String())
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

func (b *Body) Map(mapper *configVO.Mapper) *Body {
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

func (b *Body) Projection(projection *configVO.Projection) *Body {
	// se não for do tipo json, ou a projeção for nil, ou vazia, retornamos o body atual, sem modificações
	if helper.IsNotEqualTo(b.contentType, enum.ContentTypeJson) || helper.IsNil(projection) || projection.IsEmpty() {
		return b
	}
	// criamos um novo body com o json projetado a partir do objeto de valor
	return b.projectionJson(projection)
}

func (b *Body) Modify(modify *Modify) *Body {
	// instanciamos o valor a ser usado na modificação
	newValue := modify.ValueAsString()

	// instanciamos o body modificado e o erro caso aconteça
	var modifiedBody *Body
	var err error

	// abaixo verificamos qual ação desejada para modificar o valor do body
	switch modify.Action() {
	case configEnum.ModifierActionAdd:
		modifiedBody, err = b.Add(modify.Key(), newValue)
	case configEnum.ModifierActionApd:
		modifiedBody, err = b.Append(modify.Key(), newValue)
	case configEnum.ModifierActionSet:
		modifiedBody, err = b.Set(modify.Key(), newValue)
	case configEnum.ModifierActionRpl:
		modifiedBody, err = b.Replace(modify.Key(), newValue)
	case configEnum.ModifierActionDel:
		modifiedBody, err = b.Delete(modify.Key())
	default:
		return b
	}

	// tratamos o erro e retornamos o próprio body
	if helper.IsNotNil(err) {
		logger.Warning("Error modify body:", err)
	}

	return modifiedBody
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

func (b *Body) ToCase(nomenclature configEnum.Nomenclature) *Body {
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

func (b *Body) mapJson(mapper *configVO.Mapper) *Body {
	// damos o parse do json do body
	parsedJson := gjson.Parse(b.String())
	// se for array chamamos o mapJsonArray
	if parsedJson.IsArray() {
		return b.mapJsonArray(mapper, parsedJson)
	}
	// se for um objeto chamamos o mapJsonObject
	return b.mapJsonObject(mapper, parsedJson)
}

func (b *Body) mapJsonArray(mapper *configVO.Mapper, jsonArray gjson.Result) *Body {
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

func (b *Body) mapJsonObject(mapper *configVO.Mapper, jsonObject gjson.Result) *Body {
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

func (b *Body) mapText(mapper *configVO.Mapper) *Body {
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

func (b *Body) projectionJson(projection *configVO.Projection) *Body {
	// damos o parse do json do body
	parsedJson := gjson.Parse(b.String())
	// se for array chamamos o projectionJsonArray
	if parsedJson.IsArray() {
		return b.projectionJsonArray(projection, parsedJson)
	}
	// se for um objeto chamamos o projectionJsonObject
	return b.projectionJsonObject(projection, parsedJson)
}

func (b *Body) projectionJsonObject(projection *configVO.Projection, jsonObject gjson.Result) *Body {
	// verificamos se o tipo de projeção é Rejection, se for executamos as regras
	if helper.Equals(projection.Type(), configEnum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonObject(projection, jsonObject)
	}
	// se não for do tipo rejeição, então é adição ou todos, os dois aplicam as mesmas regras
	return b.projectionAdditionJsonObject(projection, jsonObject)
}

func (b *Body) projectionAdditionJsonObject(projection *configVO.Projection, jsonObject gjson.Result) *Body {
	// iniciamos o json vazio
	projectedJson := "{}"
	// iteramos o projection vo
	for _, key := range projection.Keys() {
		// se for rejection pulamos e não adicionamos no novo json
		if projection.IsRejection(key) {
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

func (b *Body) projectionRejectionJsonObject(projection *configVO.Projection, jsonObject gjson.Result) *Body {
	// iniciamos o json com o valor atual
	projectionJson := jsonObject.String()
	// iteramos o projection vo
	for _, key := range projection.Keys() {
		// removemos o campo com base na chave de projeção
		projectionJson, _ = sjson.Delete(projectionJson, key)
	}
	// retornamos o body com o novo json
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectionJson),
	}
}

func (b *Body) projectionJsonArray(projection *configVO.Projection, jsonArray gjson.Result) *Body {
	// iniciamos o json array vazio
	projectedArray := "[]"
	// iteramos o json array atual para projetar json caso seja um array de objeto
	jsonArray.ForEach(func(key, value gjson.Result) bool {
		// processamos o index do array
		projectedArray = b.projectionJsonArrayCurrentIndex(projection, projectedArray, value)
		// retorna true para continuar iterando
		return true
	})
	// se ele quer filtrar por index
	projectedArray = b.projectionJsonArrayNumericKeys(projection, projectedArray)
	// retornamos o body com o resultado da projeção
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectedArray),
	}
}

func (b *Body) projectionJsonArrayCurrentIndex(projection *configVO.Projection, projectedArray string, value gjson.Result,
) string {
	if value.IsObject() {
		// caso ele seja objeto, chamamos o projectionJsonObject para projetar o json objeto do index atual
		projectedObject := b.projectionJsonObject(projection, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedObject.Raw())
	} else if value.IsArray() {
		// caso ele seja array, chamamos o mesmo novamente para iterar o array filho
		projectedSubArray := b.projectionJsonArray(projection, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedSubArray.Raw())
	} else {
		// caso ele não seja json, não tem nada para projetar, apenas adicionamos o valor
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(value))
	}
	// retornamos o array com o index inserido
	return projectedArray
}

func (b *Body) projectionJsonArrayNumericKeys(projection *configVO.Projection, projectedJson string) string {
	// se ele não contém ao menos uma chave de projeção numérica retornamos o json já projetado
	if projection.NotContainsNumericKey() {
		return projectedJson
	}
	// verificamos se o tipo da projeção numérica é rejeição, para seguir com as regras
	if helper.Equals(projection.TypeNumeric(), configEnum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonArray(projection, projectedJson)
	}
	// caso não seja rejection, ou ele é all, ou addition, usamos as mesmas regras
	return b.projectionAdditionJsonArray(projection, projectedJson)
}

func (b *Body) projectionAdditionJsonArray(projection *configVO.Projection, projectedJson string) string {
	// transformamos o projectedJson em um gjson.Result
	parsedProjectedJson := gjson.Parse(projectedJson)
	// instanciamos a projeção de array zerada
	projectedArray := "[]"
	// iteramos a projeção para projetar os index mencionados
	for _, key := range projection.Keys() {
		if helper.IsNumeric(key) && projection.IsAddition(key) {
			jsonValue := parsedProjectedJson.Get(key)
			if jsonValue.Exists() {
				projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(jsonValue))
			}
		}
	}
	// retornamos o array projetado
	return projectedArray
}

func (b *Body) projectionRejectionJsonArray(projection *configVO.Projection, projectedJson string) string {
	// transformamos o projectedJson em um gjson.Result
	parsedProjectedJson := gjson.Parse(projectedJson)
	// instanciamos a projeção de array vazia
	projectedArray := "[]"
	// iteramos a lista
	parsedProjectedJson.ForEach(func(key, value gjson.Result) bool {
		if helper.NotContains(projection.Keys(), key.String()) {
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

func convertKeysToCase(jsonStr string, nomenclature configEnum.Nomenclature) string {
	parsedJson := gjson.Parse(jsonStr)

	jsonStrCase := "{}"
	if parsedJson.IsArray() {
		jsonStrCase = "[]"
	}

	parsedJson.ForEach(func(key, value gjson.Result) bool {
		newKey := nomenclature.Parse(key.String())
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
