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

package dto

import "github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"

// SettingView represents the configuration view for the application.
// It includes properties such as version, founder, code helpers, and various counts.
// It also contains a setting object of type GopenView, which represents the detailed configuration.
type SettingView struct {
	// Version represents the version of the application.
	Version string `json:"version,omitempty"`
	// VersionDate represents the date of the application version.
	VersionDate string `json:"version-date,omitempty"`
	// Founder represents the founder of a software application.
	Founder string `json:"founder,omitempty"`
	// CodeHelpers represents the code helpers configuration in the SettingView struct.
	// It is a string field that represents the code helpers for the application.
	CodeHelpers string `json:"code-helpers,omitempty"`
	// Endpoints represents the number of APIs in the Gopen application.
	Endpoints int `json:"endpoints"`
	// Middlewares represents the number of middlewares in the SettingView struct.
	// It is an integer field that specifies the count of middlewares used in the Gopen application.
	Middlewares int `json:"middlewares"`
	// Backends represents the number of backends configured in the SettingView struct.
	Backends int `json:"backends"`
	// Modifiers represents the count of modifiers in the SettingView struct.
	Modifiers int `json:"modifiers"`
	// Setting represents the detailed configuration view for the Gopen Json application.
	Setting *vo.GopenJson `json:"setting"`
}
