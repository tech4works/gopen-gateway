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
)

// cmdLoggerOptions is the configuration options for the logger package.
// It specifies the custom text to be displayed after the log prefix.
var cmdLoggerOptions = logger.Options{
	CustomAfterPrefixText: "CMD",
	HideAllArgs:           true,
}

// PrintLogo prints the API Gateway logo along with the provided version string.
func PrintLogo(version string) {
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

`, version)
}

// PrintTitle prints the provided title with a decorated format using the cmdLoggerProvider's PrintInfof method.
func PrintTitle(title string) {
	PrintInfof("-----------------------< %s%s%s >-----------------------", logger.StyleBold, title,
		logger.StyleReset)
}

// PrintInfo prints an informational log message using the logger package.
func PrintInfo(msg ...any) {
	logger.InfoOpts(cmdLoggerOptions, msg...)
}

// PrintInfof is a function that prints an information log message with formatting capabilities.
func PrintInfof(format string, msg ...any) {
	logger.InfoOptsf(format, cmdLoggerOptions, msg...)
}

// PrintWarn prints a warning log message using the logger package.
func PrintWarn(msg ...any) {
	logger.WarnOpts(cmdLoggerOptions, msg...)
}

// PrintWarnf logs a warning message with the given format and arguments using the logger package.
func PrintWarnf(format string, msg ...any) {
	logger.WarnOptsf(format, cmdLoggerOptions, msg...)
}

func PrintError(msg ...any) {
	logger.ErrorOpts(cmdLoggerOptions, msg...)
}
