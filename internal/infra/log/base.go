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
// Inclui atributos estruturados: tag, prefixo e trace context.
func emitOTel(ctx context.Context, lvl Level, tag, prefix, msgText string) {
	logger := global.GetLoggerProvider().Logger(instrumentationName)

	var rec otellog.Record
	rec.SetTimestamp(time.Now().UTC())
	rec.SetSeverity(levelToOTelSeverity(lvl))
	rec.SetSeverityText(lvl.Text())
	rec.SetBody(otellog.StringValue(removeAnsiCodes(msgText)))

	attrs := []otellog.KeyValue{
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

	rec.AddAttributes(attrs...)

	logger.Emit(ctx, rec)
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
