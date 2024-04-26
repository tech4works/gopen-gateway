package config

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

import "github.com/GabrielHCataldo/go-logger/logger"

// loggerOptions is the configuration options for the logger package.
// It specifies the custom text to be displayed after the log prefix.
var loggerOptions = logger.Options{
	CustomAfterPrefixText: "CMD",
	HideArgCaller:         true,
}

// PrintInfoLogCmd prints an informational log message using the logger package.
func PrintInfoLogCmd(msg ...any) {
	logger.InfoSkipCallerOpts(2, loggerOptions, msg...)
}

// PrintInfoLogCmdf is a function that prints an information log message with formatting capabilities.
func PrintInfoLogCmdf(format string, msg ...any) {
	logger.InfoSkipCallerOptsf(format, 2, loggerOptions, msg...)
}

// PrintWarningLogCmd prints a warning log message using the logger package.
func PrintWarningLogCmd(msg ...any) {
	logger.WarningSkipCallerOpts(2, loggerOptions, msg...)
}

// PrintWarningLogCmdf logs a warning message with the given format and arguments using the logger package.
func PrintWarningLogCmdf(format string, msg ...any) {
	logger.WarningSkipCallerOptsf(format, 2, loggerOptions, msg...)
}
