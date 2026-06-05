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

package log

import (
	"context"
	"fmt"
	"strings"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/checker"
)

const instrumentationName = "gopen-gateway"

// Print loga uma mensagem no nível especificado com tag e prefixo.
// Filtra pelo nível resolvido. Emite no console e exporta via OTel em background.
func Print(lvl Level, tag, prefix string, msg ...any) {
	if lvl.Disallowed() {
		return
	}

	msgText := escapeSpecialChars(fmt.Sprint(msg...))
	printConsole(lvl, tag, prefix, msgText)
	go emitOTel(context.Background(), lvl, tag, prefix, msgText)
}

// Printf loga uma mensagem formatada no nível especificado com tag e prefixo.
// Filtra pelo nível resolvido. Emite no console e exporta via OTel em background.
func Printf(lvl Level, tag, prefix, format string, msg ...any) {
	if lvl.Disallowed() {
		return
	}

	msgText := escapeSpecialChars(fmt.Sprintf(format, msg...))
	printConsole(lvl, tag, prefix, msgText)
	go emitOTel(context.Background(), lvl, tag, prefix, msgText)
}

// PrintCtx loga uma mensagem com contexto para propagação de trace.
// Filtra pelo nível resolvido. Emite no console e exporta via OTel em background.
func PrintCtx(ctx context.Context, lvl Level, tag, prefix string, msg ...any) {
	if lvl.Disallowed() {
		return
	}

	msgText := escapeSpecialChars(fmt.Sprint(msg...))
	printConsole(lvl, tag, prefix, msgText)
	go emitOTel(ctx, lvl, tag, prefix, msgText)
}

// PrintfCtx loga uma mensagem formatada com contexto para propagação de trace.
// Filtra pelo nível resolvido. Emite no console e exporta via OTel em background.
func PrintfCtx(ctx context.Context, lvl Level, tag, prefix, format string, msg ...any) {
	if lvl.Disallowed() {
		return
	}

	msgText := escapeSpecialChars(fmt.Sprintf(format, msg...))
	printConsole(lvl, tag, prefix, msgText)
	go emitOTel(ctx, lvl, tag, prefix, msgText)
}

// printConsole imprime a mensagem formatada no console com cores ANSI.
func printConsole(lvl Level, tag, prefix, msgText string) {
	tagText := BuildTagText(tag)
	levelText := BuildLevelText(lvl)

	if checker.IsNotEmpty(prefix) {
		fmt.Printf("[%s] %s %s %s\n", tagText, levelText, prefix, msgText)
	} else {
		fmt.Printf("[%s] %s %s\n", tagText, levelText, msgText)
	}
}

// emitOTel envia o log record ao OpenTelemetry LoggerProvider global.
// Inclui atributos estruturados: tag, prefixo, trace context e extra fields extraídos da mensagem.
func emitOTel(ctx context.Context, lvl Level, tag, prefix, msgText string) {
	logger := global.GetLoggerProvider().Logger(instrumentationName)

	cleanMsg := removeAnsiCodes(msgText)

	var rec otellog.Record
	rec.SetTimestamp(time.Now().UTC())
	rec.SetSeverity(levelToOTelSeverity(lvl))
	rec.SetSeverityText(lvl.Text())
	rec.SetBody(otellog.StringValue(cleanMsg))

	attrs := []otellog.KeyValue{
		otellog.String("level", lvl.Text()),
		otellog.String("tag", tag),
	}

	if checker.IsNotEmpty(prefix) {
		attrs = append(attrs, otellog.String("prefix", removeAnsiCodes(prefix)))
	}

	// Extrai trace_id do span context se disponível
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		attrs = append(attrs, otellog.String("trace.id", spanCtx.TraceID().String()))
	}

	// Extrai pares chave=valor da mensagem como extra fields.
	// Campos já mapeados como atributos dedicados são excluídos dos extras.
	fields := extractFields(cleanMsg)
	reservedKeys := map[string]bool{
		"level": true, "tag": true, "prefix": true,
		"header": true, "body": true,
		"method": true, "url": true, "status_code": true, "duration": true,
		"broker": true, "path": true,
	}
	for k, v := range fields {
		if reservedKeys[k] {
			continue
		}
		attrs = append(attrs, otellog.String("extra."+k, v))
	}

	rec.AddAttributes(attrs...)

	logger.Emit(ctx, rec)
}

// extractFields analisa a mensagem de log em busca de pares chave=valor e retorna um mapa.
// Suporta valores sem espaço (chave=valor) e valores entre aspas (chave="valor com espaço").
// Chaves aceitas: [a-zA-Z_][a-zA-Z0-9_.]*
func extractFields(msg string) map[string]string {
	fields := make(map[string]string)
	remaining := msg

	for len(remaining) > 0 {
		remaining = strings.TrimLeft(remaining, " ")
		if len(remaining) == 0 {
			break
		}

		key, value, rest, found := parseField(remaining)
		if found {
			fields[key] = value
			remaining = rest
		} else {
			// Avança para a próxima palavra
			idx := strings.IndexByte(remaining, ' ')
			if idx < 0 {
				break
			}
			remaining = remaining[idx+1:]
		}
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}

// parseField tenta parsear um par chave=valor no início de s.
// Valores são trimados de espaços para compensar formatação ANSI removida.
func parseField(s string) (key, value, rest string, found bool) {
	if len(s) == 0 || !isKeyStart(rune(s[0])) {
		return "", "", s, false
	}

	i := 1
	for i < len(s) && isKeyChar(rune(s[i])) {
		i++
	}

	if i >= len(s) || s[i] != '=' {
		return "", "", s, false
	}

	key = s[:i]
	i++ // pula o '='

	// Pula espaços iniciais do valor (resíduos de formatação ANSI removida)
	for i < len(s) && s[i] == ' ' {
		i++
	}

	if i >= len(s) {
		return key, "", "", true
	}

	// Valor entre aspas
	if s[i] == '"' {
		i++ // pula aspa de abertura
		end := strings.IndexByte(s[i:], '"')
		if end < 0 {
			return key, strings.TrimSpace(s[i:]), "", true
		}
		value = strings.TrimSpace(s[i : i+end])
		rest = s[i+end+1:]
		return key, value, rest, true
	}

	// Valor sem aspas: vai até o próximo espaço
	end := strings.IndexByte(s[i:], ' ')
	if end < 0 {
		return key, strings.TrimSpace(s[i:]), "", true
	}
	value = strings.TrimSpace(s[i : i+end])
	rest = s[i+end:]
	return key, value, rest, true
}

// isKeyStart verifica se o rune pode iniciar uma chave de field.
func isKeyStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

// isKeyChar verifica se o rune pode fazer parte de uma chave de field.
func isKeyChar(r rune) bool {
	return isKeyStart(r) || (r >= '0' && r <= '9') || r == '.' || r == '-'
}

// levelToOTelSeverity converte Level para otel.Severity.
func levelToOTelSeverity(lvl Level) otellog.Severity {
	switch lvl {
	case TraceLevel:
		return otellog.SeverityTrace
	case DebugLevel:
		return otellog.SeverityDebug
	case InfoLevel:
		return otellog.SeverityInfo
	case WarnLevel:
		return otellog.SeverityWarn
	case ErrorLevel:
		return otellog.SeverityError
	case FatalLevel:
		return otellog.SeverityFatal
	default:
		return otellog.SeverityInfo
	}
}

// escapeSpecialChars substitui caracteres de controle por representações literais.
func escapeSpecialChars(s string) string {
	return strings.NewReplacer("\n", "\\n", "\r", "\\r", "\t", "\\t").Replace(s)
}

// removeAnsiCodes remove códigos ANSI de escape de uma string.
func removeAnsiCodes(s string) string {
	result := make([]byte, 0, len(s))
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Avança até o fim do código ANSI
			j := i + 2
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}
			if j < len(s) {
				j++ // pula a letra final
			}
			i = j
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}
