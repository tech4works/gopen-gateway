package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type Mapper struct {
	keys   []string
	values map[string]string
}

func (m *Mapper) IsEmpty() bool {
	return helper.IsEmpty(m)
}

func (m *Mapper) Exists(key string) bool {
	return helper.Contains(m.keys, key)
}

func (m *Mapper) Keys() []string {
	return m.keys
}

func (m *Mapper) Get(key string) string {
	return m.values[key]
}

func (m *Mapper) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) || helper.Equals(strings.TrimSpace(string(data)), "{}") {
		return nil
	}

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, string(data))

	m.keys = []string{}
	m.values = map[string]string{}

	for field := iter.ReadObject(); helper.IsNotEmpty(field); field = iter.ReadObject() {
		m.keys = append(m.keys, field)
		m.values[field] = iter.ReadString()
	}

	return iter.Error
}

func (m *Mapper) MarshalJSON() ([]byte, error) {
	if m.IsEmpty() {
		return []byte("null"), nil
	}

	var data []string
	for _, key := range m.Keys() {
		value := m.values[key]
		data = append(data, fmt.Sprintf("%s:%s", key, value))
	}
	obj := fmt.Sprintf("{%s}", strings.Join(data, ","))
	return []byte(obj), nil
}
