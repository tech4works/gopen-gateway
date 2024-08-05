package nomenclature

import (
	"github.com/iancoleman/strcase"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type provider struct {
}

func New() domain.Nomenclature {
	return provider{}
}

func (p provider) Parse(nomenclature enum.Nomenclature, key string) string {
	switch nomenclature {
	case enum.NomenclatureCamel:
		return strcase.ToCamel(key)
	case enum.NomenclatureLowerCamel:
		return strcase.ToLowerCamel(key)
	case enum.NomenclatureSnake:
		return strcase.ToSnake(key)
	case enum.NomenclatureScreamingSnake:
		return strcase.ToScreamingSnake(key)
	case enum.NomenclatureKebab:
		return strcase.ToKebab(key)
	case enum.NomenclatureScreamingKebab:
		return strcase.ToScreamingKebab(key)
	default:
		return key
	}
}
