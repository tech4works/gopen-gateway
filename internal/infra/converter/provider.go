package converter

import "github.com/GabrielHCataldo/gopen-gateway/internal/domain"

type provider struct {
}

func New() domain.Converter {
	return provider{}
}

func (p provider) ConvertJSONToXML(bs []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p provider) ConvertTextToXML(bs []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p provider) ConvertXMLToJSON(bs []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p provider) ConvertTextToJSON(bs []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
