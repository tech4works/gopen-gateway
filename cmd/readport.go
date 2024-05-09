package main

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

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"os"
)

// JsonConfig represents the JSON configuration structure.
type JsonConfig struct {
	// Port represents the port number in a JSON configuration structure. The port number must be in
	// the range of 1 to 65535.
	Port int `json:"port"`
}

// main is the entry point of the program. It opens the file passed as an argument and reads
// the JSON configuration from it.
// If any error occurs while opening or decoding the file, it prints an error message and exits with a non-zero
// status code.
// The JSON configuration must have a valid port number in the range of 1 to 65535. If the port is out of range, it
// prints an error message and exits with a non-zero status code.
// Finally, it prints the port number to the standard output.
// The function signature is: `func main()`.
func main() {
	file, err := os.Open(os.Args[1])
	if helper.IsNotNil(err) {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	data := JsonConfig{}
	err = json.NewDecoder(file).Decode(&data)
	if helper.IsNotNil(err) {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON: %v\n", err)
		os.Exit(1)
	} else if helper.IsLessThan(data.Port, 1) || helper.IsGreaterThan(data.Port, 65535) {
		fmt.Fprintf(os.Stderr, "Port %d is out of range 1-65535", data.Port)
		os.Exit(1)
	}

	fmt.Print(data.Port)
}
