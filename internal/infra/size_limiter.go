package infra

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	domainmapper "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"io"
	"net/http"
)

// sizeLimiterProvider represents a provider for size limiting functionality.
// It contains fields to store the maximum sizes for the header, body, and multipart memory.
// The `Allow` method is responsible for checking the size of the request header and body,
// and returning an error if the size exceeds the maximum allowed size.
type sizeLimiterProvider struct {
	// maxHeaderSize represents the maximum size allowed for request headers.
	// It is a field of the sizeLimiterProvider struct, which is responsible for
	// size limiting functionality. The size is stored as a vo.Bytes value.
	maxHeaderSize vo.Bytes
	// maxBodySize represents the maximum size allowed for request bodies.
	// It is a field of the sizeLimiterProvider struct, which is responsible for size limiting functionality.
	// The size is stored as a vo.Bytes value.
	maxBodySize vo.Bytes
	// maxMultipartMemorySize represents the maximum size allowed for multipart
	// memory in the sizeLimiterProvider struct which is responsible for size
	// limiting functionality. The size is stored as a vo.Bytes value.
	// It is used to determine the maximum body size for requests with
	// "multipart/form-data" content type.
	maxMultipartMemorySize vo.Bytes
}

// SizeLimiterProvider is an interface that represents a provider for size limiting functionality.
type SizeLimiterProvider interface {
	// Allow checks whether the request is allowed based on the rate limit.
	// It takes a *http.Request as input parameter.
	// If the request is not allowed, it returns an error.
	Allow(request *http.Request) error
}

// NewSizeLimiterProvider returns a new SizeLimiterProvider with the specified maximum sizes for the header, body,
// and multipart memory.
func NewSizeLimiterProvider(limiterVO vo.Limiter) SizeLimiterProvider {
	return sizeLimiterProvider{
		maxHeaderSize:          limiterVO.MaxHeaderSize(),
		maxBodySize:            limiterVO.MaxBodySize(),
		maxMultipartMemorySize: limiterVO.MaxMultipartMemorySize(),
	}
}

// Allow checks the size of the header and body of the request.
// If the header size exceeds the maximum allowed size, it returns an error with
// the message "header too large error: permitted limit is {maxHeaderSize}".
// It then checks the "Content-Type" of the request. If it contains "multipart/form-data",
// it uses the maxMultipartMemorySize as the maximum body size. Otherwise, it uses the maxBodySize.
// It reads the request body using http.MaxBytesReader and checks if the read is successful.
// If not, it returns an error with the message "error payload too large: permitted limit is {maxBytesReader}".
// If everything goes well, it sets the request body to the read bytes using io.NopCloser.
// Finally, it returns nil to indicate success.
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

// GetHeadersSize calculates the size of the headers in the given http.Request.
// It iterates through each header and adds the length of the name, the separating ": ", and the values.
// Each value is separated by ", " and the last value has the trailing ", " removed.
// It also adds "\r\n" to the end of each header line and "\r\n" at the end of all headers.
// The total size of the headers is returned.
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
