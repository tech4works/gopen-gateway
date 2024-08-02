package domain

import (
	"context"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type Converter interface {
	ConvertJSONToXML(bs []byte) ([]byte, error)
	ConvertTextToXML(bs []byte) ([]byte, error)
	ConvertXMLToJSON(bs []byte) ([]byte, error)
	ConvertTextToJSON(bs []byte) ([]byte, error)
}

type JSONPath interface {
	Parse(raw string) JSONValue
	ForEach(raw string, iterator func(key string, value JSONValue) bool)
	Add(raw, path, value string) (string, error)
	AppendOnArray(raw, value string) (string, error)
	Set(raw, path, value string) (string, error)
	Replace(raw, path, value string) (string, error)
	Delete(raw, path string) (string, error)
	Get(raw, path string) JSONValue
}

type JSONValue interface {
	Get(path string) JSONValue
	ForEach(iterator func(key string, value JSONValue) bool)
	Exists() bool
	NotExists() bool
	IsObject() bool
	IsArray() bool
	Raw() string
	String() string
	Interface() any
}

type Nomenclature interface {
	Parse(nomenclature enum.Nomenclature, key string) string
}

type Store interface {
	Set(ctx context.Context, key string, value *vo.CacheResponse) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (*vo.CacheResponse, error)
	Close() error
}
