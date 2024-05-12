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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	xj "github.com/basgys/goxml2json"
	"github.com/clbanning/mxj/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
	"strings"
)

// Body represents the content and format of an HTTP request or response body.
type Body struct {
	// contentType is an enumeration type that represents the format of the content.
	// It can have the following values: ContentTypePlainText, ContentTypeJson.
	contentType     ContentType
	contentEncoding ContentEncoding
	// value represents the content of an HTTP httpRequest or httpResponse body. It is stored as a bytes.Buffer object.
	value *bytes.Buffer
}

// NewBody returns a new Body object based on the provided content type, content encoding, and buffer.
// It checks if the buffer is empty, and if it is, it returns nil.
// It converts the content type and content encoding strings to their respective enumeration values using
// the ContentTypeFromString and ContentEncodingFromString functions.
// If the content encoding is a valid enumeration value, it calls the NewBodyByContentEncoding function,
// passing the content type, content encoding, and buffer as arguments.
// If the content encoding is not valid, it creates a new Body object with the provided content type and buffer.
//
// Parameters:
// - contentType: A string representing the content type.
// - contentEncoding: A string representing the content encoding.
// - buffer: A pointer to a bytes.Buffer object.
//
// Returns:
// - A pointer to a Body object.
func NewBody(contentType, contentEncoding string, buffer *bytes.Buffer) *Body {
	if helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType:     NewContentType(contentType),
		contentEncoding: NewContentEncoding(contentEncoding),
		value:           buffer,
	}
}

// NewBodyByString creates a new Body object with a content type set to enum.ContentTypePlainText and a value
// set to the converted buffer of the string parameter. It returns the created Body object.
func NewBodyByString(s string) *Body {
	if helper.IsEmpty(s) {
		return nil
	}
	return &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer(s),
	}
}

// NewBodyByJson constructs a new Body object with the content type set to ContentTypeJson.
// It converts the provided data to a buffer using helper.SimpleConvertToBuffer, and assigns it to the value field.
// If the provided data is nil or empty, it returns nil.
func NewBodyByJson(a any) *Body {
	if helper.IsNil(a) || helper.IsEmpty(a) {
		return nil
	}
	return &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer(a),
	}
}

// NewBodyByError returns a new Body object based on the provided path and error.
// It constructs an errorBody from the given error and path by calling the newErrorBody function.
// If the errorBody is nil, it returns nil.
// Otherwise, it constructs a Body object with the content type set to ContentTypeJson and the value
// set to the buffer representation of the errorBody object.
func NewBodyByError(path string, err error) *Body {
	errBody := newErrorBody(path, err)
	if helper.IsNil(errBody) {
		return nil
	}
	return &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer(errBody),
	}
}

// NewBodyByHttpBackendResponse returns a new Body object based on the provided HttpBackendResponse.
// It constructs the default response body with initial fields "ok" and "code" from the HttpBackendResponse.
// Then it constructs the body with default values.
// If the body of the index is nil, it returns the body with only the default fields.
// If GroupByType() returns true, it aggregates the body based on the key.
// Otherwise, it aggregates all the JSON fields in bodyHistory.
// It returns the modified body.
func NewBodyByHttpBackendResponse(index int, httpBackendResponse *HttpBackendResponse) *Body {
	bodyJson := "{}"
	bodyJson, _ = sjson.Set(bodyJson, "ok", httpBackendResponse.Ok())
	bodyJson, _ = sjson.Set(bodyJson, "code", httpBackendResponse.StatusCode())

	body := &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer(bodyJson),
	}

	if helper.IsNil(httpBackendResponse.Body()) {
		return body
	}

	if httpBackendResponse.GroupByType() {
		return body.AggregateByKey(httpBackendResponse.Key(index), httpBackendResponse.Body())
	}
	return body.Aggregate(httpBackendResponse.Body())
}

// NewBodyAggregateByKey creates a new Body object by aggregating the provided body with the given key.
// It creates a new Body object with a content type of enum.ContentTypeJson and a value of "{}",
// then calls the AggregateByKey method on the new Body object, passing in the key and the provided body.
// If the provided body is nil, it returns nil.
func NewBodyAggregateByKey(key string, anotherBody *Body) *Body {
	if helper.IsNil(anotherBody) {
		return nil
	}

	body := &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer("{}"),
	}
	return body.AggregateByKey(key, anotherBody)
}

// NewBodyBySlice returns a new Body object based on the provided slice of Body objects.
// If the sliceOfBodies is empty, it returns nil.
// Otherwise, it creates a new Body object with ContentTypeJson and the value obtained by converting
// the sliceOfBodies to a buffer using helper.SimpleConvertToBuffer function.
func NewBodyBySlice(sliceOfBodies []*Body) *Body {
	if helper.IsEmpty(sliceOfBodies) {
		return nil
	}
	return &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer(sliceOfBodies),
	}
}

// NewBodyByCache returns a new Body object based on the provided cacheBody.
// If the cacheBody is nil, it returns nil.
// Otherwise, it creates a new Body object with the contentType and value of the cacheBody.
// The value is converted to the type *bytes.Buffer before assigning to the Body.
func NewBodyByCache(cacheBody *CacheBody) *Body {
	if helper.IsNil(cacheBody) {
		return nil
	}
	return &Body{
		contentType: cacheBody.ContentType,
		value:       (*bytes.Buffer)(cacheBody.Value),
	}
}

// NewEmptyBodyJson returns a new Body object with content type as "JSON" and value as "{}".
func NewEmptyBodyJson() *Body {
	return &Body{
		contentType: NewContentTypeJson(),
		value:       helper.SimpleConvertToBuffer("{}"),
	}
}

// ContentType returns the ContentType value of the Body instance.
// It is used to access the content type of the request body.
// The content type can be used to determine how the request body should be interpreted or processed.
// If the Body instance is nil, it returns an empty ContentType value.
func (b *Body) ContentType() ContentType {
	return b.contentType
}

func (b *Body) ContentEncoding() ContentEncoding {
	return b.contentEncoding
}

func (b *Body) HasContentEncoding() bool {
	return b.contentEncoding.Valid()
}

func (b *Body) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(b.RawBytes())
}

// Aggregate aggregates the values of anotherBody into the current Body instance.
// If anotherBody is nil, it returns the current Body instance.
// Otherwise, it calls the merge method to aggregate the values of anotherBody.
func (b *Body) Aggregate(anotherBody *Body) *Body {
	if helper.IsNil(anotherBody) {
		return b
	}
	return b.merge(anotherBody)
}

func (b *Body) AggregateByKey(key string, anotherBody *Body) *Body {
	if b.ContentType().IsNotJson() || helper.IsNil(anotherBody) {
		return b
	}

	value := anotherBody.Raw()
	mergedBodyRaw := setJsonKeyValue(b.Raw(), key, value)
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mergedBodyRaw),
	}
}

// Interface returns the interface representation of the Body instance.
// It checks the content type of the Body. If the content type is ContentTypeJson,
// it parses the body byte array using gjson.ParseBytes and returns its value.
// Otherwise, it returns the string representation of the Body value.
// The returned interface can be either a map or a string.
func (b *Body) Interface() any {
	if b.ContentType().IsJson() {
		return gjson.ParseBytes(b.Bytes()).Value()
	}
	return b.value.String()
}

func (b *Body) Json() ([]byte, error) {
	if b.ContentType().IsText() {
		return helper.SimpleConvertToBytes(fmt.Sprintf("{\"text\": \"%v\"}", b.value)), nil
	} else if b.ContentType().IsJson() {
		return b.Bytes(), nil
	} else if b.ContentType().IsXml() {
		return convertXmlToJson(b.Bytes())
	}
	return []byte{}, nil
}

func (b *Body) Xml() ([]byte, error) {
	if b.ContentType().IsText() {
		return []byte(fmt.Sprintf("<root>%s</root>", b.value)), nil
	} else if b.ContentType().IsJson() {
		return convertJsonToXml(b.Bytes())
	} else if b.ContentType().IsXml() {
		return b.Bytes(), nil
	}
	return []byte{}, nil
}

func compressByEncoding(bs []byte, contentEncoding ContentEncoding) ([]byte, error) {
	if contentEncoding.IsGzip() {
		return helper.CompressWithGzip(bs)
	} else if contentEncoding.IsDeflate() {
		return helper.CompressWithDeflate(bs)
	}
	return bs, nil
}

func (b *Body) Bytes() []byte {
	if b.ContentEncoding().IsGzip() {
		decompressBytes, err := helper.DecompressWithGzip(b.RawBytes())
		if helper.IsNil(err) {
			return decompressBytes
		}
	} else if b.ContentEncoding().IsDeflate() {
		decompressBytes, err := helper.DecompressWithDeflate(b.RawBytes())
		if helper.IsNil(err) {
			return decompressBytes
		}
	}
	return b.RawBytes()
}

func (b *Body) RawBytes() []byte {
	return b.value.Bytes()
}

func (b *Body) String() string {
	return string(b.Bytes())
}

// CompactString returns the compact string representation of the body.
// It converts the body to string and passes it to the helper.CompactString function.
// The helper.CompactString function removes leading/trailing spaces and multiple consecutive spaces.
//
// Returns:
// - The compact string representation of the body.
func (b *Body) CompactString() string {
	return helper.CompactString(b.String())
}

// Raw returns the raw string representation of the Body instance.
// If the content type of the Body is ContentTypePlainText,
// the string will be quoted using strconv.Quote function.
// Otherwise, the string will be returned as is.
// The raw string is parsed using parseStringValueToRaw function
// to convert it into a raw JSON string.
func (b *Body) Raw() string {
	s := b.String()
	if b.ContentType().IsText() {
		s = strconv.Quote(s)
	}
	return parseStringValueToRaw(s)
}

// MarshalJSON marshals the value of the Body instance into JSON format.
// It uses helper.ConvertToBytes to convert the value to bytes and return it along with no errors.
func (b *Body) MarshalJSON() ([]byte, error) {
	return b.value.Bytes(), nil
}

func (b *Body) Add(key string, value string) (*Body, error) {
	if b.ContentType().IsText() {
		return b.addString(value), nil
	} else if b.ContentType().IsJson() {
		return b.addJson(key, value)
	}
	return b, nil
}

func (b *Body) Append(key string, value string) (*Body, error) {
	if b.ContentType().IsText() {
		return b.appendString(value), nil
	} else if b.ContentType().IsJson() {
		return b.appendJson(key, value)
	}
	return b, nil
}

func (b *Body) Set(key string, value string) (*Body, error) {
	if b.ContentType().IsText() {
		return b.setString(value), nil
	} else if b.ContentType().IsJson() {
		return b.setJson(key, value)
	}
	return b, nil
}

func (b *Body) Replace(key string, value string) (*Body, error) {
	if b.ContentType().IsText() {
		return b.replaceString(key, value), nil
	} else if b.ContentType().IsJson() {
		return b.replaceJson(key, value)
	}
	return b, nil
}

func (b *Body) Delete(key string) (*Body, error) {
	if b.ContentType().IsText() {
		return b.replaceString(key, ""), nil
	} else if b.ContentType().IsJson() {
		return b.deleteJson(key)
	}
	return b, nil
}

func (b *Body) Map(mapper *Mapper) *Body {
	if helper.IsNotNil(mapper) && mapper.IsNotEmpty() {
		if b.ContentType().IsText() {
			return b.mapText(mapper)
		} else if b.ContentType().IsJson() {
			return b.mapJson(mapper)
		}
	}
	return b
}

// Projection applies a projection to the Body instance.
// It returns the modified Body with the JSON projected based on the Projection object.
// If the Body content type is not JSON, or the Projection object is nil or empty,
// it returns the current Body instance without modifications.
// The JSON projection is created based on the Projection object's keys and values.
func (b *Body) Projection(projection *Projection) *Body {
	if b.ContentType().IsNotJson() || helper.IsNil(projection) || projection.IsEmpty() {
		return b
	}
	return b.projectionJson(projection)
}

// Modify modifies the body of the request or response based on the provided Modifier.
// It takes a Modifier, an HttpRequest, and an HttpResponse as parameters.
// It returns the modified Body instance or the original Body instance if the modifier action is not recognized.
// If an error occurs during the modification, it logs a warning message using the logger package.
//
// The Modify method first retrieves the new value by calling the ValueAsString method on the provided Modifier.
// Then, it creates variables to hold the modified Body instance and an error.
//
// The method then uses a switch statement to determine the modifier action.
// Based on the action, it calls the corresponding modification method on the Body instance (Add, Append, Set, Replace,
// or Delete), passing the modifier key and the new value as arguments.
//
// If an error is returned from the modification method, it logs a warning message using the logger package.
//
// Finally, the method returns the modified Body instance or the original Body instance.
func (b *Body) Modify(modifier *Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) *Body {
	newValue := modifier.ValueAsString(httpRequest, httpResponse)

	var modifiedBody *Body
	var err error

	switch modifier.Action() {
	case enum.ModifierActionAdd:
		modifiedBody, err = b.Add(modifier.Key(), newValue)
	case enum.ModifierActionApd:
		modifiedBody, err = b.Append(modifier.Key(), newValue)
	case enum.ModifierActionSet:
		modifiedBody, err = b.Set(modifier.Key(), newValue)
	case enum.ModifierActionRpl:
		modifiedBody, err = b.Replace(modifier.Key(), newValue)
	case enum.ModifierActionDel:
		modifiedBody, err = b.Delete(modifier.Key())
	default:
		return b
	}

	if helper.IsNotNil(err) {
		logger.Warning("Error modify body:", err)
	}

	return modifiedBody
}

func (b *Body) ModifyContentType(contentType ContentType, contentEncoding ContentEncoding) *Body {
	var bs []byte
	var err error
	if contentType.IsJson() {
		bs, err = b.Json()
	} else if contentType.IsXml() {
		bs, err = b.Xml()
	} else {
		bs = b.Bytes()
	}

	if helper.IsNotNil(err) {
		logger.Warningf("Error modify content-type: %s err: %s", contentType, err)
		return b
	}

	bs, err = compressByEncoding(bs, contentEncoding)
	if helper.IsNotNil(err) {
		logger.Warningf("Error modify content-encoding: %s err: %s", contentEncoding, err)
		return b
	}

	return &Body{
		contentType:     contentType,
		contentEncoding: contentEncoding,
		value:           helper.SimpleConvertToBuffer(bs),
	}
}

func (b *Body) OmitEmpty() *Body {
	if b.ContentType().IsText() {
		return b.omitEmptyText()
	} else if b.ContentType().IsJson() {
		return b.omitEmptyJson()
	}
	return b
}

func (b *Body) ToCase(nomenclature enum.Nomenclature) *Body {
	if b.ContentType().IsNotJson() {
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
// If the contentType is ContentTypePlainText, it uses the mergeString function to merge the string values of the bodies.
// The merged Body value is then used to create a new Body instance with the same contentType,
// which is returned as the result.
// If the contentType is neither ContentTypeJson nor ContentTypePlainText, it returns the current Body instance unchanged.
func (b *Body) merge(anotherBody *Body) *Body {
	if helper.IsNotNil(anotherBody) {
		if b.ContentType().IsText() {
			return b.mergeString(anotherBody.String())
		} else if b.ContentType().IsJson() {
			return b.mergeJSON(anotherBody.Raw())
		}
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

// mergeJSON merges the JSON values of a Body instance with the provided JSON string.
// It iterates over the keys and values of the JSON string and merges them into the existing JSON data of the Body.
// The merged result is then used to create a new Body instance with the same contentType,
// which is returned as the result.
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

// mapText applies a mapper to the Body instance and returns a new modified Body instance.
// It replaces all occurrences of keys in the Body's value with their corresponding new values from the mapper.
// The modified Body instance has the same content type as the original Body instance.
// If the mapper is empty or there are no key replacements, it returns the original Body instance.
func (b *Body) mapText(mapper *Mapper) *Body {
	mappedText := b.String()

	for _, key := range mapper.Keys() {
		newKey := mapper.Get(key)
		if helper.IsNotEqualTo(key, newKey) {
			mappedText = strings.ReplaceAll(mappedText, key, newKey)
		}
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedText),
	}
}

// mapJson applies the specified mapper to the JSON content of the Body instance.
// It parses the JSON content into a gjson.Result object and checks if it is an array or object.
// If it is an array, it calls the mapJsonArray method on the Body and returns the modified Body instance.
// If it is an object, it calls the mapJsonObject method on the Body and returns the modified Body instance.
// If the JSON content is neither an array nor an object, it returns the original Body instance without modifications.
func (b *Body) mapJson(mapper *Mapper) *Body {
	parsedJson := gjson.Parse(b.String())
	if parsedJson.IsArray() {
		return b.mapJsonArray(mapper, parsedJson)
	}
	return b.mapJsonObject(mapper, parsedJson)
}

// mapJsonArray maps the JSON array using the specified mapper. It iterates through each element of the array,
// checks if it is an object or an array. If it is an object, it calls the mapJsonObject method on the Body to map
// the object.
// If the element is an array, it recursively calls the mapJsonArray method to map the sub-array.
// If the element is neither an object nor an array, it parses the value and sets it in the mapped array.
// It returns a new Body instance with the mapped JSON array.
func (b *Body) mapJsonArray(mapper *Mapper, jsonArray gjson.Result) *Body {
	mappedArray := "[]"

	jsonArray.ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			projectedObject := b.mapJsonObject(mapper, value)
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", projectedObject.Raw())
		} else if value.IsArray() {
			projectedSubArray := b.mapJsonArray(mapper, value)
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", projectedSubArray.Raw())
		} else {
			mappedArray, _ = sjson.SetRaw(mappedArray, "-1", parseValueToRaw(value))
		}
		return true
	})

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedArray),
	}
}

// mapJsonObject maps the JSON object using the specified mapper. It iterates through each key in the mapper's keys list,
// compares it to its corresponding value, and if they are different, it tries to get the value from the jsonObject
// using the original key. If the value exists, it sets the value in the newMappedJson with the new key, otherwise,
// it proceeds to the next key.
// If the newMappedJson is successfully created, it deletes the original key from the newMappedJson.
// If any error occurs during the process, the original key and value are set in the newMappedJson.
// Finally, it creates a new Body instance with the mapped JSON object and returns it.
func (b *Body) mapJsonObject(mapper *Mapper, jsonObject gjson.Result) *Body {
	mappedJson := jsonObject.String()

	for _, key := range mapper.Keys() {
		newKey := mapper.Get(key)
		if helper.Equals(key, newKey) {
			continue
		}
		jsonValue := jsonObject.Get(key)
		if !jsonValue.Exists() {
			continue
		}
		newMappedJson, err := sjson.SetRaw(mappedJson, newKey, parseValueToRaw(jsonValue))
		if helper.IsNil(err) {
			mappedJson, _ = sjson.Delete(newMappedJson, key)
		} else {
			mappedJson, _ = sjson.SetRaw(mappedJson, key, parseValueToRaw(jsonValue))
		}
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(mappedJson),
	}
}

// projectionJson applies a projection to the Body instance.
// It parses the JSON content of the Body, and depending on whether it is an array or object,
// it applies the projection using either the b.projectionJsonArray or b.projectionJsonObject method.
// The parsed JSON is then returned as a modified Body instance with the projected JSON.
// If the Body content type is not JSON, the provided Projection is nil, or the Projection is empty,
// the method returns the current Body instance without modifications.
//
// Parameters:
// - projection: The Projection object that defines the keys and values for the projection.
//
// Returns:
// - *Body: The modified Body instance with the projected JSON, or the current instance if no modifications were made.
func (b *Body) projectionJson(projection *Projection) *Body {
	parsedJson := gjson.Parse(b.String())
	if parsedJson.IsArray() {
		return b.projectionJsonArray(projection, parsedJson)
	}
	return b.projectionJsonObject(projection, parsedJson)
}

// projectionJsonObject applies a projection to the Body instance when the JSON content is an object.
// It checks the type of the projection and calls the appropriate method:
// projectionRejectionJsonObject or projectionAdditionJsonObject.
// If the projection.Type() is enum.ProjectionTypeRejection, it calls projectionRejectionJsonObject method.
// Otherwise, it calls projectionAdditionJsonObject method.
// The method returns the modified Body instance based on the projected JSON.
//
// Parameters:
// - projection: The Projection object that defines the keys and values for the projection.
// - jsonObject: The gjson.Result representing the parsed JSON content of the Body.
//
// Returns:
// - *Body: The modified Body instance with the projected JSON.
func (b *Body) projectionJsonObject(projection *Projection, jsonObject gjson.Result) *Body {
	if helper.Equals(projection.Type(), enum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonObject(projection, jsonObject)
	}
	return b.projectionAdditionJsonObject(projection, jsonObject)
}

// projectionAdditionJsonObject projects a subset of JSON properties from the given jsonObject
// based on the provided projection.
// It returns a new instance of Body with a JSON string containing the projected properties.
// The projected properties are determined by the projection's keys and values.
// If a key in the projection exists in the jsonObject and is not rejected, its value will be added to the projected JSON.
// The resulting projected JSON will be set as the value of the Body, with the same content type as the original Body.
// The projectedJson variable is initialized as an empty JSON object.
// The for loop iterates over each key in the projection's keys.
// If the key is a rejection, the loop continues to the next key.
// Otherwise, the value corresponding to the key is retrieved from the jsonObject.
// If the value exists in the jsonObject, it is added to the projectedJson using the sjson.SetRaw function.
// Finally, a new instance of Body is returned with the contentType and the value set to the projectedJson.
func (b *Body) projectionAdditionJsonObject(projection *Projection, jsonObject gjson.Result) *Body {
	projectedJson := "{}"

	for _, key := range projection.Keys() {
		if projection.IsRejection(key) {
			continue
		}
		jsonValue := jsonObject.Get(key)
		if jsonValue.Exists() {
			projectedJson, _ = sjson.SetRaw(projectedJson, key, parseValueToRaw(jsonValue))
		}
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectedJson),
	}
}

// projectionRejectionJsonObject removes the keys specified in the projection from the given JSON object
// and returns a modified Body instance.
//
// Parameters:
// - projection: The Projection object that defines the keys to be removed.
// - jsonObject: The gjson.Result representing the parsed JSON content.
//
// Returns:
// - *Body: The modified Body instance with the projected JSON content.
func (b *Body) projectionRejectionJsonObject(projection *Projection, jsonObject gjson.Result) *Body {
	projectionJson := jsonObject.String()

	for _, key := range projection.Keys() {
		projectionJson, _ = sjson.Delete(projectionJson, key)
	}

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectionJson),
	}
}

// projectionJsonArray applies the projection on a JSON array.
// It iterates over the current JSON array and performs projections according to the provided Projection object.
// If the array element is an object, it calls the projectionJsonObject method to project the object.
// If the array element is another JSON array, it recursively calls the projectionJsonArray method.
// If the array element is not a JSON value, it adds the element as is.
// The projected JSON array is returned as a modified Body instance.
//
// Parameters:
// - projection: The Projection object that defines the keys and values for the projection.
// - jsonArray: The input JSON array on which the projection is to be applied.
//
// Returns:
// - *Body: The modified Body instance with the projected JSON array.
func (b *Body) projectionJsonArray(projection *Projection, jsonArray gjson.Result) *Body {
	projectedArray := "[]"

	jsonArray.ForEach(func(key, value gjson.Result) bool {
		projectedArray = b.projectionJsonArrayCurrentIndex(projection, projectedArray, value)
		return true
	})
	projectedArray = b.projectionJsonArrayNumericKeys(projection, projectedArray)

	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(projectedArray),
	}
}

// projectionJsonArrayCurrentIndex applies the projection on a JSON array at a specific index.
// It iterates over the current JSON array element at the given index and performs projections according to the provided
// Projection object.
// If the array element is an object, it calls the projectionJsonObject method to project the object.
// If the array element is another JSON array, it recursively calls the projectionJsonArray method.
// If the array element is not a JSON value, it adds the element as is.
// The projected JSON array is returned as a modified string representation.
//
// Parameters:
// - projection: The Projection object that defines the keys and values for the projection.
// - projectedArray: The string representation of the current JSON array on which the projection is to be applied.
// - value: The JSON value at the given index for projection.
//
// Returns:
// - string: The modified string representation of the projected JSON array.
func (b *Body) projectionJsonArrayCurrentIndex(projection *Projection, projectedArray string, value gjson.Result,
) string {
	if value.IsObject() {
		projectedObject := b.projectionJsonObject(projection, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedObject.Raw())
	} else if value.IsArray() {
		projectedSubArray := b.projectionJsonArray(projection, value)
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", projectedSubArray.Raw())
	} else {
		projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(value))
	}
	return projectedArray
}

// projectionJsonArrayNumericKeys projects the JSON array with numeric keys based on the given projection.
// If the projection does not contain any numeric key, it returns the input JSON string.
// If the projection type is numeric rejection, it applies rejection projection to the JSON array.
// Otherwise, it applies addition projection to the JSON array.
func (b *Body) projectionJsonArrayNumericKeys(projection *Projection, projectedJson string) string {
	if projection.NotContainsNumericKey() {
		return projectedJson
	} else if helper.Equals(projection.TypeNumeric(), enum.ProjectionTypeRejection) {
		return b.projectionRejectionJsonArray(projection, projectedJson)
	}
	return b.projectionAdditionJsonArray(projection, projectedJson)
}

// projectionAdditionJsonArray returns a JSON array string containing the projected values from
// the provided JSON string, based on the given projection.
// It iterates through the projection keys and adds the corresponding values from the parsed JSON
// into the array string if they exist.
// The array string is returned as a result.
func (b *Body) projectionAdditionJsonArray(projection *Projection, projectedJson string) string {
	parsedProjectedJson := gjson.Parse(projectedJson)
	projectedArray := "[]"

	for _, key := range projection.Keys() {
		if !helper.IsNumeric(key) || projection.IsRejection(key) {
			continue
		}
		jsonValue := parsedProjectedJson.Get(key)
		if jsonValue.Exists() {
			projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(jsonValue))
		}
	}

	return projectedArray
}

// projectionRejectionJsonArray applies rejection projection to the given JSON array based on the given projection.
// It returns a new JSON array with the rejected elements removed.
// The projection specifies the keys to be rejected from the JSON array.
// The rejected elements are identified by keys that are not present in the projection.
//
// If the projection is empty or does not have any keys, the original JSON array is returned.
// The function takes a projection and a projected JSON string as input.
//
// The function parses the projected JSON string and iterates over each key-value pair.
// If a key is not present in the projection, the value is added to the new JSON array.
// Otherwise, the value is ignored.
//
// The function uses the `helper.NotContains` function to check if a key is not present in the projection.
// It uses the `sjson.SetRaw` function to add the value to the new JSON array.
// The final JSON array is returned as a string.
//
// An example usage of this method can be found in the `projectionJsonArrayNumericKeys` method of the `Body` struct.
// The `projectionJsonArrayNumericKeys` method applies numeric rejection projection to a JSON array with numeric keys,
// using the `projectionRejectionJsonArray` method to perform the rejection.
// If the projection does not contain any numeric key, the method returns the input JSON string.
func (b *Body) projectionRejectionJsonArray(projection *Projection, projectedJson string) string {
	parsedProjectedJson := gjson.Parse(projectedJson)
	projectedArray := "[]"

	parsedProjectedJson.ForEach(func(key, value gjson.Result) bool {
		if helper.NotContains(projection.Keys(), key.String()) {
			projectedArray, _ = sjson.SetRaw(projectedArray, "-1", parseValueToRaw(value))
		}
		return true
	})

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
// If the contentType is ContentTypePlainText, omitEmptyText() will be called.
// If the contentType is neither ContentTypeJson nor ContentTypePlainText, the original Body instance will be returned.
func (b *Body) omitEmptyText() *Body {
	s := helper.CleanAllRepeatSpaces(b.String())
	return &Body{
		contentType: b.contentType,
		value:       helper.SimpleConvertToBuffer(s),
	}
}

func (b *Body) Length() int {
	return len(b.RawBytes())
}

func (b *Body) LengthStr() string {
	return helper.SimpleConvertToString(b.Length())
}

func convertJsonToXml(bs []byte) ([]byte, error) {
	mapJson, err := mxj.NewMapJson(bs)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return mapJson.Xml()
}

func convertXmlToJson(bs []byte) ([]byte, error) {
	jsonData, err := xj.Convert(bytes.NewBuffer(bs))
	if helper.IsNotNil(err) {
		return nil, err
	}
	return jsonData.Bytes(), nil
}

// setJsonKeyValue sets the value of a given key in a JSON string. If the key already exists in the JSON string and
// its value is not null, the function aggregates the existing value with the new value. If the key does not exist,
// the function adds the key-value pair to the JSON string. The updated JSON string is returned.
// If the key already exists but its value is null, it sets the value of the key to the new value.
//
// Parameters:
// - jsonStr: The original JSON string.
// - key: The key to set or add.
// - value: The new value for the key.
//
// Returns:
//
// The updated JSON string.
func setJsonKeyValue(jsonStr, key string, value string) string {
	result := gjson.Get(jsonStr, key)

	if result.Exists() && helper.IsNotEqualTo(result.Type, gjson.Null) {
		jsonStr, _ = sjson.SetRaw(jsonStr, key, aggregateJsonValue(result, value))
	} else {
		jsonStr, _ = sjson.SetRaw(jsonStr, key, parseStringValueToRaw(value))
	}

	return jsonStr
}

// parseValueToRaw parses the value to a raw string representation.
// If the value is of type Null, it returns the string "null". Otherwise, it returns the raw value as is.
func parseValueToRaw(value gjson.Result) string {
	if helper.Equals(value.Type, gjson.Null) {
		return "null"
	}
	return value.Raw
}

// parseStringValueToRaw parses the provided string value into a raw JSON string.
// It uses the gjson.Parse() function from the gjson package to parse the value.
// If the parsed value is of type Null, it returns the string "null".
// Otherwise, it returns the raw value as is.
func parseStringValueToRaw(value string) string {
	parse := gjson.Parse(value)
	if helper.Equals(parse.Type, gjson.Null) {
		return "null"
	}
	return parse.Raw
}

// aggregateJsonValue aggregates the provided JSON value with a new value.
// It takes a gjson.Result value and a new value as inputs.
// If the provided value is an array, it appends the new value to the existing array.
// If the provided value is not an array, it creates a new array with the existing value and appends the new value.
// The function then constructs a JSON string representation of the aggregated array.
// It skips null and empty values, and returns the final JSON string.
//
// Parameters:
// - value: The original JSON value.
// - newValue: The new value to aggregate.
//
// Returns:
// The JSON string representation of the aggregated array.
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

// removeAllEmptyFields removes all empty fields from the JSON string.
// It recursively traverses the JSON structure and checks if a field is empty.
// If a field is empty, it is deleted from the JSON string using the sjson.Delete function.
// If a field is not empty, its value is updated using the sjson.SetRaw function with the raw value.
// The function uses the parseValueToRaw function to parse the value to a raw string representation.
// The final modified JSON string is returned as the result.
//
// Please note that the original JSON string is modified in-place.
// Also, the parseValueToRaw function is used to parse the value to a raw string representation.
// If the value is of type Null, it returns the string "null". Otherwise, it returns the raw value as is.
func removeAllEmptyFields(jsonStr string) string {
	gjson.Parse(jsonStr).ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() || value.IsArray() {
			subJsonStr := removeAllEmptyFields(parseValueToRaw(value))
			value = gjson.Parse(subJsonStr)
		}
		if helper.IsEmpty(value.Value()) {
			jsonStr, _ = sjson.Delete(jsonStr, key.String())
		} else {
			jsonStr, _ = sjson.SetRaw(jsonStr, key.String(), parseValueToRaw(value))
		}
		return true
	})
	return jsonStr
}

// convertKeysToCase converts the keys of a JSON string to the specified case format
// based on the provided Nomenclature. It recursively processes nested objects and arrays.
// The function returns the modified JSON string with the converted keys.
func convertKeysToCase(jsonStr string, nomenclature enum.Nomenclature) string {
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
