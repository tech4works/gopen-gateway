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

type JsonConfig struct {
	Port int `json:"port"`
}

func main() {
	// lemos o arquivo passado como argumento
	file, err := os.Open(os.Args[1])
	// caso tenha dado erro, ja retornamos como falha ao abrir o arquivo
	if helper.IsNotNil(err) {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	// fechamos a leitura do mesmo
	defer file.Close()

	// instanciamos o jsonConfig para ler a porta que está no json
	data := JsonConfig{}
	err = json.NewDecoder(file).Decode(&data)
	// caso ocorra um erro retornamos
	if helper.IsNotNil(err) {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON: %v\n", err)
		os.Exit(1)
	} else if helper.IsLessThan(data.Port, 1) || helper.IsGreaterThan(data.Port, 65535) {
		fmt.Fprintf(os.Stderr, "Port %d is out of range 1-65535", data.Port)
		os.Exit(1)
	}

	fmt.Print(data.Port)
}
