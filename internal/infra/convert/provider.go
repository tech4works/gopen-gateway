/*
 * Copyright 2024 Tech4Works
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

package convert

import (
	"bytes"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"github.com/clbanning/mxj/v2"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
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
	return converter.ToBytesWithErr(fmt.Sprintf("<root>%s</root>", string(bs)))
}

func (p provider) ConvertXMLToJSON(bs []byte) ([]byte, error) {
	jsonData, err := xj.Convert(bytes.NewBuffer(bs))
	if checker.NonNil(err) {
		return nil, err
	}

	return jsonData.Bytes(), nil
}

func (p provider) ConvertTextToJSON(bs []byte) ([]byte, error) {
	return converter.ToBytesWithErr(fmt.Sprintf("{\"text\": \"%v\"}", string(bs)))
}
