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

package boot

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"os"
)

type noop struct{}

func (l *noop) Error(_ string) {}

func (l *noop) Infof(_ string, _ ...interface{}) {}

type log struct {
	options logger.Options
}

func newLogger() app.Logger {
	return log{
		options: logger.Options{
			CustomAfterPrefixText: fmt.Sprintf("[%s%s%s]", logger.StyleBold, "GOPEN", logger.StyleReset),
			HideAllArgs:           true,
		},
	}
}

func (l log) PrintLogo() {
	fmt.Printf(` 
 ######    #######  ########  ######## ##    ##
##    ##  ##     ## ##     ## ##       ###   ##
##        ##     ## ##     ## ##       ####  ##
##   #### ##     ## ########  ######   ## ## ##
##    ##  ##     ## ##        ##       ##  ####
##    ##  ##     ## ##        ##       ##   ###
 ######    #######  ##        ######## ##    ##
-----------------------------------------------
Best open source API Gateway!            %s
-----------------------------------------------
2024 • Gabriel Cataldo.

`, os.Getenv("VERSION"))
}

func (l log) PrintTitle(title string) {
	l.PrintInfof("-----------------------< %s%s%s >-----------------------", logger.StyleBold, title, logger.StyleReset)
}

func (l log) PrintInfo(msg ...any) {
	logger.InfoOpts(l.options, msg...)
}

func (l log) PrintInfof(format string, msg ...any) {
	logger.InfoOptsf(format, l.options, msg...)
}

func (l log) PrintWarn(msg ...any) {
	logger.WarnOpts(l.options, msg...)
}

func (l log) PrintWarnf(format string, msg ...any) {
	logger.WarnOptsf(format, l.options, msg...)
}

func (l log) PrintError(msg ...any) {
	logger.ErrorOpts(l.options, msg...)
}
