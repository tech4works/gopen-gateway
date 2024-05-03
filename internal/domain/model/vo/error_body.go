package vo

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"time"
)

// errorBody represents the structure of a httpResponse body containing error details.
// It is used to serialize the error details to JSON format.
//
// This struct is typically used in conjunction with the newBodyByError function to generate a JSON httpResponse body
// with error details based on the Endpoint and error provided.
type errorBody struct {
	// File represents the file name or path where the error occurred.
	File string `json:"file"`
	// Line represents the line number where the error occurred
	Line int `json:"line"`
	// Endpoint represents the endpoint path where the error occurred.
	Endpoint string `json:"endpoint"`
	// Message represents the error message.
	Message string `json:"message"`
	// Timestamp represents the timestamp when the error occurred.
	Timestamp time.Time `json:"timestamp"`
}

func newErrorBody(path string, err error) *errorBody {
	// obtemos o detalhe do erro usando a lib go-errors
	detailsErr := errors.Details(err)
	if helper.IsNil(detailsErr) {
		return nil
	}
	// com os detalhes, construímos o objeto de retorno padrão de erro da API Gateway
	return &errorBody{
		File:      detailsErr.GetFile(),
		Line:      detailsErr.GetLine(),
		Endpoint:  path,
		Message:   detailsErr.GetMessage(),
		Timestamp: time.Now(),
	}
}
