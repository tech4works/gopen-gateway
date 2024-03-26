package infra

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	appmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"io"
	"net/http"
)

type sizeLimiterProvider struct {
	limiterDTO dto.Limiter
}

func NewSizeLimiterProvider(limiterDTO dto.Limiter) external.SizeLimiterProvider {
	return sizeLimiterProvider{
		limiterDTO: limiterDTO,
	}
}

func (s sizeLimiterProvider) Allow(request *http.Request) error {
	// checamos primeiramente o tamanho do header
	headerSize := s.GetHeadersSize(request)
	if helper.IsGreaterThan(s.limiterDTO.MaxHeaderSize, headerSize) {
		return appmapper.NewErrHeaderTooLarge(s.limiterDTO.MaxHeaderSize)
	}

	// verificamos qual Content-Type fornecido, para obter a config real
	maxBytesReader := s.limiterDTO.MaxBodySize
	if helper.ContainsIgnoreCase(request.Header.Get("Content-Type"), "multipart/form-data") {
		maxBytesReader = s.limiterDTO.MaxMultipartMemorySize
	}
	// verificamos o tamanho utilizando o maxBytesReader, e logo em seguida se der certo, voltamos o body para requisição
	read := http.MaxBytesReader(nil, request.Body, int64(maxBytesReader))
	bodyBytes, err := io.ReadAll(read)
	if helper.IsNotNil(err) {
		return appmapper.NewErrPayloadTooLarge(maxBytesReader)
	}
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// se tudo ocorrer bem retornamos nil
	return nil
}

func (s sizeLimiterProvider) GetHeadersSize(r *http.Request) int {
	size := 0
	for name, values := range r.Header {
		// o tamanho da chave mais o ': ' que o separa do valor
		size += len(name) + 2
		for _, value := range values {
			size += len(value)
			// Cada valor múltiplo está separado por ', '
			size += 2
		}
		// Subtraímos os dois últimos caracteres ', ' do último valor
		size -= 2
		// Adicionamos o '\r\n' que termina cada linha de header
		size += 2
	}
	// Adicionamos o '\r\n' final dos headers
	size += 2
	return size
}
