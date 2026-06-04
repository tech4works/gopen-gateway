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

// Package log fornece logging estruturado com suporte a nível de severidade,
// controle via variáveis de ambiente e integração com OpenTelemetry para o gopen-gateway.
package log

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tech4works/checker"
)

// resolvedLevel armazena o nível mínimo de log resolvido (cache atômico).
var resolvedLevel atomic.Int32

// levelResolved garante que ResolveLevel() seja executado apenas uma vez.
var levelResolved sync.Once

const levelNotResolved int32 = -1

func init() {
	resolvedLevel.Store(levelNotResolved)
}

// Level representa o nível de severidade de uma entrada de log.
type Level int

// Níveis de log disponíveis, em ordem crescente de severidade.
const (
	// TraceLevel é o nível mais detalhado, para rastreamento de operações de infraestrutura.
	TraceLevel Level = iota
	// DebugLevel é usado para informações de depuração durante o desenvolvimento.
	DebugLevel
	// InfoLevel é o nível padrão, usado para eventos operacionais normais.
	InfoLevel
	// WarnLevel é usado para situações inesperadas que não impedem o funcionamento.
	WarnLevel
	// ErrorLevel é usado para falhas que afetam uma operação específica.
	ErrorLevel
	// FatalLevel é usado para falhas críticas irrecuperáveis (ex: panics recuperados).
	FatalLevel
)

// NewLevel converte uma string para o Level correspondente (case-insensitive).
// Retorna InfoLevel para strings não reconhecidas ou vazias.
func NewLevel(s string) Level {
	if checker.EqualsIgnoreCase(s, "trace") {
		return TraceLevel
	} else if checker.EqualsIgnoreCase(s, "debug") {
		return DebugLevel
	} else if checker.EqualsIgnoreCase(s, "warn") {
		return WarnLevel
	} else if checker.EqualsIgnoreCase(s, "error") {
		return ErrorLevel
	} else if checker.EqualsIgnoreCase(s, "fatal") {
		return FatalLevel
	}
	return InfoLevel
}

// ResolveLevel resolve o nível mínimo de log com base em LOG_LEVEL e o ambiente.
// Se LOG_LEVEL estiver definido, usa diretamente. Caso contrário:
//   - local/dev → debug
//   - stg/prd → info
//
// Deve ser chamado uma vez durante o boot. Chamadas subsequentes são no-op.
func ResolveLevel(environment string) {
	levelResolved.Do(func() {
		logLevel := os.Getenv("LOG_LEVEL")
		if checker.IsNotEmpty(logLevel) {
			resolvedLevel.Store(int32(NewLevel(logLevel)))
		} else if checker.Equals(environment, "prd") || checker.Equals(environment, "stg") {
			resolvedLevel.Store(int32(InfoLevel))
		} else {
			resolvedLevel.Store(int32(DebugLevel))
		}
	})
}

// GetResolvedLevel retorna o nível mínimo de log resolvido.
// Se ResolveLevel() ainda não foi chamado, faz fallback para InfoLevel.
func GetResolvedLevel() Level {
	v := resolvedLevel.Load()
	if checker.IsLessThan(v, 0) {
		return InfoLevel
	}
	return Level(v)
}

// Allowed retorna true se este nível for igual ou superior ao mínimo resolvido.
func (l Level) Allowed() bool {
	return GetResolvedLevel() <= l
}

// Disallowed retorna true se este nível estiver abaixo do mínimo resolvido.
func (l Level) Disallowed() bool {
	return !l.Allowed()
}

// String retorna a representação colorida do nível para exibição no console.
func (l Level) String() string {
	return fmt.Sprint(l.color(), l.abbreviation(), StyleReset)
}

// Text retorna o nome textual do nível em minúsculas.
func (l Level) Text() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

// abbreviation retorna a sigla de 3 letras do nível.
func (l Level) abbreviation() string {
	switch l {
	case TraceLevel:
		return "TRC"
	case DebugLevel:
		return "DBG"
	case InfoLevel:
		return "INF"
	case WarnLevel:
		return "WRN"
	case ErrorLevel:
		return "ERR"
	case FatalLevel:
		return "FTL"
	default:
		return "INF"
	}
}

// color retorna o código ANSI de cor associado ao nível.
func (l Level) color() string {
	switch l {
	case TraceLevel:
		return "\x1b[32m"
	case DebugLevel:
		return "\x1b[36m"
	case InfoLevel:
		return "\x1b[34m"
	case WarnLevel:
		return "\x1b[33m"
	case ErrorLevel, FatalLevel:
		return "\x1b[31m"
	default:
		return "\x1b[0m"
	}
}
