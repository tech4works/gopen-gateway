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

package mapper

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"os"
)

// BuildSettingViewDTO builds a SettingView DTO using the provided GopenJson and Gopen structs.
// It initializes and returns a new SettingView object with the following properties:
// - Version: value of "VERSION" environment variable.
// - VersionDate: value of "VERSION_DATE" environment variable.
// - Founder: value of "FOUNDER" environment variable.
// - Contributors: converted value of "CONTRIBUTORS" environment variable to an integer using helper.SimpleConvertToInt function.
// - Endpoints: count of endpoints in the Gopen struct.
// - Middlewares: count of middlewares in the Gopen struct.
// - Backends: count of all the beforewares, backends, and afterwares in the Gopen struct.
// - Transformations: count of all data transformations in the Gopen struct.
// - Setting: a copy of the GopenJson object.
func BuildSettingViewDTO(gopenJson *vo.GopenJson, gopen *vo.Gopen) dto.SettingView {
	return dto.SettingView{
		Version:         os.Getenv("VERSION"),
		VersionDate:     os.Getenv("VERSION_DATE"),
		Founder:         os.Getenv("FOUNDER"),
		Contributors:    helper.SimpleConvertToInt(os.Getenv("CONTRIBUTORS")),
		Endpoints:       gopen.CountEndpoints(),
		Middlewares:     gopen.CountMiddlewares(),
		Backends:        gopen.CountBackends(),
		Transformations: gopen.CountAllDataTransforms(),
		Setting:         gopenJson.Filter(),
	}
}
