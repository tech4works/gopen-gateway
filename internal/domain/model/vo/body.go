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

func newBodyByErr(err error) Body {
	return Body{value: err}
}

func newBodyByAny(value any) Body {
	if orderedMap, isOrderedMap := value.(orderedmap.OrderedMap); isOrderedMap {
		return Body{value: orderedMap}
	} else if sliceOfOrderedMap, isSliceOfOrderedMap := value.([]orderedmap.OrderedMap); isSliceOfOrderedMap {
		return Body{value: sliceOfOrderedMap}
	} else if helper.IsErrorType(value) {
		return Body{value: value}
	}
	return newBody(helper.SimpleConvertToBytes(value))
}

func newBody(bytes []byte) (b Body) {
	_ = json.Unmarshal(bytes, &b)
	return b
}

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

func (b *Body) Value() any {
	if b.isOrderedMap() {
		return b.OrderedMap()
	} else if b.isSliceOfOrderedMaps() {
		return b.SliceOfOrderedMaps()
	}
	return b.value
}

func (b *Body) isOrderedMap() bool {
	_, ok := b.value.(orderedmap.OrderedMap)
	return ok
}

func (b *Body) isSliceOfOrderedMaps() bool {
	_, ok := b.value.([]orderedmap.OrderedMap)
	return ok
}

func (b *Body) OrderedMap() orderedmap.OrderedMap {
	return b.copyOrderedMap(b.value.(orderedmap.OrderedMap))
}

func (b *Body) SliceOfOrderedMaps() (copy []orderedmap.OrderedMap) {
	slicesOfOrderedMaps := b.value.([]orderedmap.OrderedMap)
	for _, orderedMap := range slicesOfOrderedMaps {
		copy = append(copy, b.copyOrderedMap(orderedMap))
	}
	return copy
}

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

func (b *Body) modifyMap(key string, value map[string]any) Body {
	// chamamos o modify do mapa ordenado passando o mapa modificado
	return Body{value: b.modifyOrderedMap(b.OrderedMap(), key, value)}
}

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

func (b *Body) IsNotEmpty() bool {
	if b.isOrderedMap() {
		orderedMap := b.OrderedMap()
		return helper.IsNotEmpty(orderedMap.Keys())
	} else if b.isSliceOfOrderedMaps() {
		return helper.IsNotEmpty(b.SliceOfOrderedMaps())
	}
	return helper.IsNotNil(b.value) && helper.IsNotEmpty(b.value)
}

func (b *Body) String() string {
	return helper.SimpleConvertToString(b.value)
}

func (b *Body) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

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

func (b *Body) MarshalXML(e *xml.Encoder, _ xml.StartElement) (err error) {
	return b.encodeXML(e, "body", b.value)
}

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

func (b *Body) encodeWithoutKeyXML(e *xml.Encoder, value any) error {
	key := strings.ToLower(reflect.TypeOf(value).String())
	return b.encodeXML(e, key, value)
}

func (b *Body) encodeXML(e *xml.Encoder, key string, value any) (err error) {
	if orderedMapCast, ok := value.(orderedmap.OrderedMap); ok {
		return b.orderedMapWithKeyXML(e, key, orderedMapCast)
	} else if slice, ok := value.([]orderedmap.OrderedMap); ok {
		return b.sliceOfOrderedMapXML(e, key, slice)
	} else {
		return b.interfaceXML(e, key, value)
	}
}
