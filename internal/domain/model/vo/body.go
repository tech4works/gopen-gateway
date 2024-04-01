package vo

import (
	"encoding/json"
	"encoding/xml"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/iancoleman/orderedmap"
	"reflect"
	"strings"
)

type Body struct {
	value any
}

// NewBodyByContentType returns a new Body based on the content type and bytes.
// If the bytes slice is empty, it returns an empty Body.
// If the content type is "application/json", it converts the bytes into a structured object.
// Otherwise, it returns a Body with the string value of the bytes.
func NewBodyByContentType(contentType string, bytes []byte) Body {
	// se vazio, retornamos nil
	if helper.IsEmpty(bytes) {
		return Body{}
	}

	// se ele for json, verificamos se body é um map ou slice para manter ordenado
	if helper.ContainsIgnoreCase(contentType, "application/json") {
		// convertemos os bytes do body em uma interface de objeto
		return newBody(bytes)
	}
	//todo: futuramente podemos trabalhar com o XML e o FORM-DATA com o modifier e envio

	return Body{value: string(bytes)}
}

// newBodyByErr returns a new Body with the error value.
// The error value is set as the value of the Body.
func newBodyByErr(err error) Body {
	return Body{value: err}
}

// newBodyByAny returns a new Body based on the value type.
// If the value is an ordered map, it returns a Body with the ordered map as the value.
// If the value is a slice of ordered maps, it returns a Body with the slice of ordered maps as the value.
// If the value is an error or a string, it returns a Body with the value as it is.
// Otherwise, it converts the value to bytes using helper.SimpleConvertToBytes() and returns a new Body with the
// converted bytes as the value.
func newBodyByAny(value any) Body {
	if orderedMap, isOrderedMap := value.(orderedmap.OrderedMap); isOrderedMap {
		return Body{value: orderedMap}
	} else if sliceOfOrderedMap, isSliceOfOrderedMap := value.([]orderedmap.OrderedMap); isSliceOfOrderedMap {
		return Body{value: sliceOfOrderedMap}
	} else if helper.IsErrorType(value) || helper.IsStringType(value) {
		return Body{value: value}
	}
	return newBody(helper.SimpleConvertToBytes(value))
}

// newBody accepts a byte slice and returns a new Body object by unmarshalling
// the byte slice as JSON. If there is an error during unmarshalling, an empty
// Body object is returned.
func newBody(bytes []byte) (b Body) {
	_ = json.Unmarshal(bytes, &b)
	return b
}

// copyOrderedMap returns a new copy of the given ordered map.
// It iterates over each key in the ordered map and retrieves its corresponding value.
// If the value exists, it sets the key-value pair in the new ordered map.
// Finally, it returns the new copy of the ordered map.
func (b *Body) copyOrderedMap(orderedMap orderedmap.OrderedMap) orderedmap.OrderedMap {
	c := orderedmap.New()
	for _, key := range orderedMap.Keys() {
		valueByKey, exists := orderedMap.Get(key)
		if exists {
			c.Set(key, valueByKey)
		}
	}
	return *c
}

// Value returns the value of the body.
// If the body is an ordered map, it returns the ordered map value.
// If the body is a slice of ordered maps, it returns the slice of ordered maps value.
// Otherwise, it returns the default value of the body.
func (b *Body) Value() any {
	if b.isOrderedMap() {
		return b.OrderedMap()
	} else if b.isSliceOfOrderedMaps() {
		return b.SliceOfOrderedMaps()
	}
	return b.value
}

// OrderedMap returns the ordered map value of the `Body` object.
// It returns a new copy of the ordered map by calling the `copyOrderedMap` method.
// The method cast the value of the `Body` object to an `OrderedMap` type
// and passes it as an argument to the `copyOrderedMap` method.
// Finally, it returns the new copy of the ordered map.
func (b *Body) OrderedMap() orderedmap.OrderedMap {
	return b.copyOrderedMap(b.value.(orderedmap.OrderedMap))
}

// SliceOfOrderedMaps returns a new slice containing copies of the ordered maps stored in b.value.
// It iterates over each ordered map in slicesOfOrderedMaps and calls the copyOrderedMap function to create a copy of each map.
// The copies are then appended to the copy slice.
// Finally, it returns the copy slice with the copied ordered maps.
func (b *Body) SliceOfOrderedMaps() (copy []orderedmap.OrderedMap) {
	slicesOfOrderedMaps := b.value.([]orderedmap.OrderedMap)
	for _, orderedMap := range slicesOfOrderedMaps {
		copy = append(copy, b.copyOrderedMap(orderedMap))
	}
	return copy
}

// Interface returns the value of the Body as an interface{}.
// It checks the type of the value to convert it into an interface{}.
// If the value is an ordered map, it returns the values of the ordered map.
// If the value is a slice of ordered maps, it converts each ordered map into a map[string]interface{} and returns a slice of those maps.
// If the value is not an ordered map or a slice of ordered maps, it returns the value as is.
// In other words, if there is no ordered map to be converted, the value is already of type interface{}.
// The returned value can be asserted to its original type when needed.
func (b *Body) Interface() any {
	// verificamos qual o tipo do valor para converter em interface
	if b.isOrderedMap() {
		orderedMap := b.OrderedMap()
		return orderedMap.Values()
	} else if b.isSliceOfOrderedMaps() {
		sliceOfOrderedMaps := b.SliceOfOrderedMaps()

		var sliceOfMaps []map[string]any
		for _, orderedMapIndex := range sliceOfOrderedMaps {
			sliceOfMaps = append(sliceOfMaps, orderedMapIndex.Values())
		}
		return sliceOfMaps
	}

	// se não tem nenhum map a ser ordenado quer dizer que ele ja é do tipo any
	return b.value
}

// IsNotEmpty checks if the Body object is not empty.
// If the Body object is representing an ordered map,
// it retrieves the keys and checks if they are not empty using the helper IsNotEmpty.
// If the Body object is representing a slice of ordered maps,
// it checks if the slice is not empty using the helper IsNotEmpty.
// If neither of the above conditions are met,
// it checks if the value of the Body object is not nil and not empty using the helper IsNotNil and IsNotEmpty.
// It returns true if the Body object is not empty, otherwise it returns false.
func (b *Body) IsNotEmpty() bool {
	if b.isOrderedMap() {
		orderedMap := b.OrderedMap()
		return helper.IsNotEmpty(orderedMap.Keys())
	} else if b.isSliceOfOrderedMaps() {
		return helper.IsNotEmpty(b.SliceOfOrderedMaps())
	}
	return helper.IsNotNil(b.value) && helper.IsNotEmpty(b.value)
}

// String returns a string representation of the current Body instance.
// It utilizes the SimpleConvertToString function from the helper package to convert the value of the Body to a string.
// The resulting string representation of the Body is returned.
func (b *Body) String() string {
	return helper.SimpleConvertToString(b.value)
}

// MarshalJSON converts the body value to JSON byte array.
// It uses the json.Marshal function to serialize the value into JSON format.
// If successful, it returns the JSON byte array representation of the value.
// Otherwise, it returns an error.
func (b *Body) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// UnmarshalJSON parses the JSON-encoded data and stores it in the Body.
// It first checks if the data is a slice of maps or a map.
// If the data is a map, it unmarshals it into an ordered map, preserving the order.
// If there are no errors, it assigns the ordered map to the Body and returns nil.
// If the data is a slice of maps, it unmarshals it into a slice of ordered maps.
// If there are no errors, it assigns the slice of ordered maps to the Body and returns nil.
// If the data is neither a map nor a slice of maps, it unmarshals it into a regular interface.
// If there are no errors, it assigns the interface to the Body and returns nil.
// If any error occurs, it assigns an empty Body and returns nil.
func (b *Body) UnmarshalJSON(bytes []byte) error {
	// convertemos pelo tipo de bytes, se for slice de mapas, e se for um mapa convertemos mantendo a ordenação
	if helper.IsMap(bytes) {
		var orderedMap orderedmap.OrderedMap
		err := json.Unmarshal(bytes, &orderedMap)
		if helper.IsNil(err) {
			*b = Body{value: orderedMap}
			return nil
		}
	} else if helper.IsSliceOfMaps(bytes) {
		// verificamos se é um slice de map
		var sliceOfOrderedMaps []orderedmap.OrderedMap
		err := json.Unmarshal(bytes, &sliceOfOrderedMaps)
		if helper.IsNil(err) {
			*b = Body{value: sliceOfOrderedMaps}
			return nil
		}
	}

	// caso os bytes não contem um map retornamos ele como interface normalmente
	var dest any
	err := json.Unmarshal(bytes, &dest)
	if helper.IsNil(err) {
		*b = Body{value: dest}
		return nil
	}

	// se tudo deu errado
	*b = Body{}
	return nil
}

// MarshalXML encodes the content of the Body struct into XML format.
// It uses the provided xml.Encoder to write the XML encoding.
// The XML element name for the Body struct is "body".
// The content of the Body is encoded by calling the encodeXML method with the provided xml.Encoder,
// the XML element name "body", and the value of the Body.
// The function returns any error encountered during the encoding process.
func (b *Body) MarshalXML(e *xml.Encoder, _ xml.StartElement) (err error) {
	return b.encodeXML(e, "body", b.value)
}

// Modify modifies a body by updating the value associated with the given key.
// If the value is a map, it is sorted and returned.
// If the value is a slice of maps, the maps are sorted and returned.
// If the value is not a map, the modified value is returned as a new Body struct.
// The modified Body is returned.
func (b *Body) Modify(key string, value any) Body {
	// verificamos se o valor é algum tipo de map para ser ordenado
	if helper.IsMapType(value) {
		return b.modifyMap(key, value.(map[string]any))
	} else if helper.IsSliceOfMapsType(value) {
		return b.modifySliceOfMaps(key, value.([]map[string]any))
	}
	// se ele não é nenhum tipo de map, retornamos o valor modificado passado
	return Body{value: value}
}

// modifyMap modifies the value of the given key in the body's map.
// It calls the modifyOrderedMap function of the ordered map, passing the modified map.
// Finally, it returns a new Body with the modified ordered map.
func (b *Body) modifyMap(key string, value map[string]any) Body {
	// chamamos o modify do mapa ordenado passando o mapa modificado
	return Body{value: b.modifyOrderedMap(b.OrderedMap(), key, value)}
}

// isOrderedMap checks if the value contained in the Body struct is of type orderedmap.OrderedMap.
// It returns true if the value is an instance of orderedmap.OrderedMap, otherwise it returns false.
func (b *Body) isOrderedMap() bool {
	_, ok := b.value.(orderedmap.OrderedMap)
	return ok
}

// isSliceOfOrderedMaps checks if the value of the Body is a slice of ordered maps.
// It attempts to type assert the value into a slice of ordered maps.
// If the type assertion succeeds, it returns true indicating that the value is a slice of ordered maps.
// Otherwise, it returns false indicating that the value is not a slice of ordered maps.
func (b *Body) isSliceOfOrderedMaps() bool {
	_, ok := b.value.([]orderedmap.OrderedMap)
	return ok
}

// modifySliceOfMaps takes a key and a slice of maps with string keys and any values.
// It initializes an empty slice to hold the modified ordered maps.
// It retrieves the current slice of ordered maps from the Body object.
// It iterates over the current slice and the slice of modified values, matching them by index.
// If the indices match, it modifies the ordered map at that index by calling modifyOrderedMap method.
// The modified ordered map is appended to the result slice.
// Finally, it returns a new Body object with the updated slice of ordered maps.
func (b *Body) modifySliceOfMaps(key string, values []map[string]any) Body {
	// inicializamos o resultado ordenado do map
	var resultSliceOfOrderedMap []orderedmap.OrderedMap

	// obtemos a fatia de mapas ordenados atual
	currentSliceOfOrderedMap := b.SliceOfOrderedMaps()

	// iteramos a fatia atual para manter os index de forma correta
	for orderedIndex, orderedMap := range currentSliceOfOrderedMap {
		// iteramos a fatia de valores modificados
		for valueIndex, valueMap := range values {
			// se for igual ao index modificamos o valor
			if helper.Equals(orderedIndex, valueIndex) {
				resultSliceOfOrderedMap = append(resultSliceOfOrderedMap, b.modifyOrderedMap(orderedMap, key, valueMap))
				break
			}
		}
	}

	// setamos o novo valor de fatia de mapa ordenado
	return Body{value: resultSliceOfOrderedMap}
}

// modifyOrderedMap is a function in the Body type. It takes an ordered map, a key string, and a map of string to any.
// The function modifies the received orderedMap based on the provided key-value pairs.
// It iteratively checks through each key in the orderedMap and updates its value if the key exists in the provided value map.
// If the key does not exist in the value map, it tries to get the value by the provided key.
// This is interpreted as a changed key name. It then updates the resultOrderedMap only if the orderedKey has not been removed.
// After iterating through the existing keys and values of orderedMap, it sets the additional key-value pairs from the
// provided value map that do not already exist in the orderedMap.
// Finally, it returns the modified resultOrderedMap.
func (b *Body) modifyOrderedMap(orderedMap orderedmap.OrderedMap, key string, value map[string]any) orderedmap.OrderedMap {
	// inicializamos o resultado ordenado do map
	resultOrderedMap := orderedmap.New()

	// alteramos oq ja tem, removendo ou alterando
	for _, orderedKey := range orderedMap.Keys() {
		// obtemos o valor a ser modificado pela chave ordenada
		valueByKey, exists := value[orderedKey]
		// caso ela não exista tentamos obter ele do map modificado
		if !exists {
			// tentamos obter o valor pela chave modificado no valor modificado
			valueByKey, exists = value[key]
			// caso ele exista, quer dizer que foi alterado o nome apenas
			if exists {
				orderedKey = key
			}
		}
		// setamos no resultado apenas se o orderedKey nao tiver sido removido
		if exists {
			resultOrderedMap.Set(orderedKey, valueByKey)
		}
	}

	// setamos os valores adicionados
	for valueKey, keyValue := range value {
		_, exists := orderedMap.Get(key)
		if !exists {
			resultOrderedMap.Set(valueKey, keyValue)
		}
	}

	// setamos o resultado do mapa ordenado modificado
	return *resultOrderedMap
}

// orderedMapXML encodes the key-value pairs of the given ordered map into XML format using the provided XML encoder.
// It iterates over each key in the ordered map and retrieves its corresponding value.
// If the value exists, it calls the `encodeXML` method to encode the key-value pair into XML.
// If any error occurs during encoding, it returns that error.
// If all key-value pairs are successfully encoded, it returns nil.
func (b *Body) orderedMapXML(e *xml.Encoder, orderedMap orderedmap.OrderedMap) error {
	for _, orderedKey := range orderedMap.Keys() {
		valueByKey, exists := orderedMap.Get(orderedKey)
		if !exists {
			continue
		}
		err := b.encodeXML(e, orderedKey, valueByKey)
		if helper.IsNotNil(err) {
			return err
		}
	}
	return nil
}

// orderedMapWithKeyXML encodes an ordered map to XML using the provided xml.Encoder and a specified key.
// It starts encoding with a StartElement containing the specified key.
// Then, it calls the orderedMapXML method to encode the actual ordered map.
// If there is an error during encoding, it returns that error.
// Finally, it encodes the End() token of the StartElement and returns nil.
func (b *Body) orderedMapWithKeyXML(e *xml.Encoder, key string, orderedMap orderedmap.OrderedMap) error {
	field := xml.StartElement{Name: xml.Name{Local: key}}
	err := e.EncodeToken(field)
	if helper.IsNotNil(err) {
		return err
	}

	err = b.orderedMapXML(e, orderedMap)
	if helper.IsNotNil(err) {
		return err
	}

	return e.EncodeToken(field.End())
}

// sliceOfOrderedMapXML encodes a slice of ordered maps to XML using the given xml.Encoder.
// It starts by encoding the start element with the specified key.
// It then iterates over each ordered map in the slice and encodes it as an "item" element.
// It uses the orderedMapXML function to encode each ordered map.
// Finally, it encodes the end element with the specified key.
// If any encoding error occurs, it returns the error.
func (b *Body) sliceOfOrderedMapXML(e *xml.Encoder, key string, slice []orderedmap.OrderedMap) error {
	field := xml.StartElement{Name: xml.Name{Local: key}}
	err := e.EncodeToken(field)
	if helper.IsNotNil(err) {
		return err
	}

	for _, value := range slice {
		indexField := xml.StartElement{Name: xml.Name{Local: "item"}}
		err = e.EncodeToken(indexField)
		if helper.IsNotNil(err) {
			return err
		}

		err = b.orderedMapXML(e, value)
		if helper.IsNotNil(err) {
			return err
		}

		err = e.EncodeToken(indexField.End())
		if helper.IsNotNil(err) {
			return err
		}
	}

	return e.EncodeToken(field.End())
}

// interfaceXML encodes the given value as XML using the provided Encoder.
// If the value is a slice, it will encode each element of the slice as a separate XML item.
// The key parameter is used as the XML element name.
// Returns any error encountered during encoding.
func (b *Body) interfaceXML(e *xml.Encoder, key string, value any) error {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Slice {
		field := xml.StartElement{Name: xml.Name{Local: key}}
		err := e.EncodeToken(field)
		if helper.IsNotNil(err) {
			return err
		}

		for i := 0; i < reflectValue.Len(); i++ {
			indexField := xml.StartElement{Name: xml.Name{Local: "item"}}
			err = e.EncodeToken(indexField)
			if helper.IsNotNil(err) {
				return err
			}

			item := reflectValue.Index(i).Interface()
			err = b.encodeWithoutKeyXML(e, item)
			if helper.IsNotNil(err) {
				return err
			}

			err = e.EncodeToken(indexField.End())
			if helper.IsNotNil(err) {
				return err
			}
		}

		return e.EncodeToken(field.End())
	}

	field := xml.StartElement{Name: xml.Name{Local: key}}
	return e.EncodeElement(value, field)
}

// encodeWithoutKeyXML encodes the given value as XML using the provided XML encoder.
// It first retrieves the type of the value and converts it to lowercase.
// Then, it calls the encodeXML method with the lowercase type as the key and the value itself.
// Finally, it returns any error that occurred during encoding.
// Note: This method does not include the key in the XML output.
func (b *Body) encodeWithoutKeyXML(e *xml.Encoder, value any) error {
	key := strings.ToLower(reflect.TypeOf(value).String())
	return b.encodeXML(e, key, value)
}

// encodeXML encodes the given value to XML using the provided XML encoder and key.
// If the value is of type orderedmap.OrderedMap, it calls the orderedMapWithKeyXML method to encode it.
// If the value is of type []orderedmap.OrderedMap, it calls the sliceOfOrderedMapXML method to encode it.
// Otherwise, it calls the interfaceXML method to encode the value as a generic interface.
// Returns any error that occurs during encoding.
func (b *Body) encodeXML(e *xml.Encoder, key string, value any) (err error) {
	if orderedMapCast, ok := value.(orderedmap.OrderedMap); ok {
		return b.orderedMapWithKeyXML(e, key, orderedMapCast)
	} else if slice, ok := value.([]orderedmap.OrderedMap); ok {
		return b.sliceOfOrderedMapXML(e, key, slice)
	} else {
		return b.interfaceXML(e, key, value)
	}
}
