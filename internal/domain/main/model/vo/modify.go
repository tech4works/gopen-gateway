package vo

import (
	configEnum "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
)

type Modify struct {
	action configEnum.ModifierAction
	key    string
	value  DynamicValue
}

func NewModify(modifier *vo.Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) *Modify {
	return &Modify{
		action: modifier.Action(),
		key:    modifier.Key(),
		value:  NewDynamicValue(modifier.Value(), httpRequest, httpResponse),
	}
}

func (m Modify) Action() configEnum.ModifierAction {
	return m.action
}

func (m Modify) Key() string {
	return m.key
}

func (m Modify) ValueAsSliceOfString() []string {
	return m.value.AsSliceOfString()
}

func (m Modify) ValueAsString() string {
	return m.value.AsString()
}
