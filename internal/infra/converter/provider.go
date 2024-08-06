package converter

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	xj "github.com/basgys/goxml2json"
	"github.com/clbanning/mxj/v2"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain"
)

type provider struct {
}

func New() domain.Converter {
	return provider{}
}

func (p provider) ConvertJSONToXML(bs []byte) ([]byte, error) {
	mapJson, err := mxj.NewMapJson(bs)
	if checker.NonNil(err) {
		return nil, err
	}
	return mapJson.Xml("root")
}

func (p provider) ConvertTextToXML(bs []byte) ([]byte, error) {
	return helper.ConvertToBytes(fmt.Sprintf("<root>%s</root>", string(bs)))
}

func (p provider) ConvertXMLToJSON(bs []byte) ([]byte, error) {
	jsonData, err := xj.Convert(bytes.NewBuffer(bs))
	if checker.NonNil(err) {
		return nil, err
	}

	return jsonData.Bytes(), nil
}

func (p provider) ConvertTextToJSON(bs []byte) ([]byte, error) {
	return helper.ConvertToBytes(fmt.Sprintf("{\"text\": \"%v\"}", string(bs)))
}
