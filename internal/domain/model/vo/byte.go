package vo

import "github.com/GabrielHCataldo/go-helper/helper"

type Bytes int64
type MegaBytes int64

func NewBytes(bytesUnit string) Bytes {
	return Bytes(helper.SimpleConvertByteUnitStrToFloat(bytesUnit))
}

func (b *Bytes) UnmarshalJSON(v []byte) error {
	s := string(v)
	if helper.IsNotEmpty(v) {
		i, err := helper.ConvertByteUnitStrToFloat(s)
		if helper.IsNotNil(err) {
			return err
		}
		*b = Bytes(i)
	}
	return nil
}

func (b *Bytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(b) || helper.IsEmpty(b) {
		return nil, nil
	}
	return []byte(b.String()), nil
}

func (b *Bytes) String() string {
	return helper.ConvertToByteUnitStr(b)
}

func (m *MegaBytes) UnmarshalJSON(v []byte) error {
	s := string(v)
	if helper.IsNotEmpty(v) {
		i, err := helper.ConvertMegaByteUnitStrToFloat(s)
		if helper.IsNotNil(err) {
			return err
		}
		*m = MegaBytes(i)
	}
	return nil
}

func (m *MegaBytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(m) || helper.IsEmpty(m) {
		return nil, nil
	}
	return []byte(m.String()), nil
}

func (m *MegaBytes) String() string {
	return helper.ConvertToMegaByteUnitStr(m)
}
