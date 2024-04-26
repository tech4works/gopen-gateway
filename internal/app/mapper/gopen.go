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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"os"
)

// BuildSettingViewDTO takes a GopenJson object and constructs a SettingView object with
// the required properties. It retrieves the version, version date, founder, and code helpers
// from environment variables, and counts the number of endpoints, middlewares, backends, and modifiers
// using methods from the GopenJson struct. It also sets the Setting property of the SettingView object
// to a copy of the GopenJson object.
//
// Parameters:
// - gopenJsonVO: A pointer to a GopenJson object that contains the configuration details for the Gopen application.
//
// Returns:
//   - A SettingView object that represents the configuration view for the application, including the version,
//     version date, founder, code helpers, counts of endpoints, middlewares, backends, modifiers,
//     and a copy of the GopenJson object.
func BuildSettingViewDTO(gopenJsonVO *vo.GopenJson) dto.SettingView {
	return dto.SettingView{
		Version:     os.Getenv("VERSION"),
		VersionDate: os.Getenv("VERSION_DATE"),
		Founder:     os.Getenv("FOUNDER"),
		CodeHelpers: os.Getenv("CODE_HELPERS"),
		Endpoints:   gopenJsonVO.CountEndpoints(),
		Middlewares: gopenJsonVO.CountMiddlewares(),
		Backends:    gopenJsonVO.CountBackends(),
		Modifiers:   gopenJsonVO.CountModifiers(),
		Setting:     gopenJsonVO.Json(),
	}
}
