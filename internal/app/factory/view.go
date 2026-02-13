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

package factory

import (
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
)

func BuildSettingView(gopen dto.Gopen) dto.SettingView {
	copied := gopen
	copied.Store = nil

	return dto.SettingView{
		Version:      "v1.0.0",
		VersionDate:  "05/10/2024",
		Founder:      "Gabriel Cataldo",
		Contributors: 1,
		Endpoints:    countEndpoints(gopen),
		Backends:     countBackends(gopen),
		Setting:      copied,
	}
}

func countEndpoints(gopen dto.Gopen) int {
	return len(gopen.Endpoints)
}

func countBackends(gopen dto.Gopen) (count int) {
	for _, endpoint := range gopen.Endpoints {
		count += len(endpoint.Beforewares) + len(endpoint.Backends) + len(endpoint.Afterwares)
	}
	return count
}
