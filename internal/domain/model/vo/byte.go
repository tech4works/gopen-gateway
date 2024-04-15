package vo

import "github.com/GabrielHCataldo/go-helper/helper"

// Bytes represents a type for storing byte values.
type Bytes int64

// MegaBytes represents a type for storing megabyte values.
type MegaBytes int64

// NewBytes creates a new Bytes value based on the given byte unit string.
// It converts the byte unit string to a float value using the helper.SimpleConvertByteUnitStrToFloat function.
func NewBytes(bytesUnit string) Bytes {
	return Bytes(helper.SimpleConvertByteUnitStrToFloat(bytesUnit))
}

// UnmarshalJSON unmarshals a JSON-encoded byte value into a Bytes value.
// It takes a byte slice v as input, converts it to a string s,
// and then checks if the string is not empty.
// If it's not empty, it converts the string to a float value using the helper.ConvertByteUnitStrToFloat function,
// and assigns the converted value to the pointer receiver b.
// Finally, it returns nil.
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

// MarshalJSON marshals a Bytes value into a JSON-encoded byte slice.
// It checks if the pointer receiver b is nil or empty. If it is,
// it returns nil and nil error.
// Otherwise, it converts the Bytes value to a string using the b.String() method,
// converts the string to a byte slice using []byte(),
// and returns the byte slice and nil error.
func (b *Bytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(b) || helper.IsEmpty(b) {
		return nil, nil
	}
	return []byte(b.String()), nil
}

// String returns a string representation of the Bytes value.
// It uses the helper.ConvertToByteUnitStr function to format the value as a byte unit string.
func (b *Bytes) String() string {
	return helper.ConvertToByteUnitStr(b)
}

// UnmarshalJSON unmarshals a JSON-encoded byte value into a MegaBytes value.
// It takes a byte slice v as input, converts it to a string s,
// and then checks if the string is not empty.
// If it's not empty, it converts the string to a float value using the helper.ConvertMegaByteUnitStrToFloat function,
// and assigns the converted value to the pointer receiver m.
// Finally, it returns nil.
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

// MarshalJSON marshals a MegaBytes value into a JSON-encoded byte slice.
// It checks if the pointer receiver m is nil or empty. If it is, it returns nil and nil.
// Otherwise, it converts the MegaBytes value to a string using the m.String() method,
// converts the string to a byte slice, and returns it along with nil as the error.
func (m *MegaBytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(m) || helper.IsEmpty(m) {
		return nil, nil
	}
	return []byte(m.String()), nil
}

// String returns a string representation of the MegaBytes value.
// It uses the helper.ConvertToMegaByteUnitStr function to convert the MegaBytes value to a string.
// It then returns the converted string.
func (m *MegaBytes) String() string {
	return helper.ConvertToMegaByteUnitStr(m)
}
