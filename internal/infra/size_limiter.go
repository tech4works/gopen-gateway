package infra

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"io"
	"net/http"
)

type sizeLimiterProvider struct {
	maxHeaderSize          vo.Bytes
	maxBodySize            vo.Bytes
	maxMultipartMemorySize vo.Bytes
}

func NewSizeLimiterProvider(maxHeaderSize, maxBodySize, maxMultipartMemorySize vo.Bytes) interfaces.SizeLimiterProvider {
	return sizeLimiterProvider{
		maxHeaderSize:          maxHeaderSize,
		maxBodySize:            maxBodySize,
		maxMultipartMemorySize: maxMultipartMemorySize,
	}
}

func (s sizeLimiterProvider) Allow(request *http.Request) error {
	// checamos primeiramente o tamanho do header
	headerSize := s.GetHeadersSize(request)
	if helper.IsGreaterThan(headerSize, s.maxHeaderSize) {
		return domainmapper.NewErrHeaderTooLarge(s.maxHeaderSize.String())
	}

	// verificamos qual Content-Type fornecido, para obter a config real
	maxBytesReader := s.maxBodySize
	if helper.ContainsIgnoreCase(request.Header.Get("Content-Type"), "multipart/form-data") {
		maxBytesReader = s.maxMultipartMemorySize
	}
	// verificamos o tamanho utilizando o maxBytesReader, e logo em seguida se der certo, voltamos o body para requisição
	read := http.MaxBytesReader(nil, request.Body, int64(maxBytesReader))
	bodyBytes, err := io.ReadAll(read)
	if helper.IsNotNil(err) {
		return domainmapper.NewErrPayloadTooLarge(maxBytesReader.String())
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
