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

package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/tech4works/gopen-gateway/internal/app"
	"os"
)

type bootLog struct {
	options logger.Options
}

func NewBoot() app.BootLog {
	return bootLog{
		options: logger.Options{
			HideAllArgs:           true,
			CustomAfterPrefixText: fmt.Sprintf("[%s%s%s]", logger.StyleBold, "BOOT", logger.StyleReset),
		},
	}
}

func (l bootLog) PrintLogo() {
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

func (l bootLog) PrintTitle(title string) {
	l.PrintInfof("-----------------------< %s%s%s >-----------------------", logger.StyleBold, title, logger.StyleReset)
}

func (l bootLog) PrintInfo(msg ...any) {
	logger.InfoOpts(l.options, msg...)
}

func (l bootLog) PrintInfof(format string, msg ...any) {
	logger.InfoOptsf(format, l.options, msg...)
}

func (l bootLog) PrintWarn(msg ...any) {
	logger.WarnOpts(l.options, msg...)
}

func (l bootLog) PrintWarnf(format string, msg ...any) {
	logger.WarnOptsf(format, l.options, msg...)
}

func (l bootLog) PrintError(msg ...any) {
	logger.ErrorOpts(l.options, msg...)
}
