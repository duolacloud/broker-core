package marshaler

import (
	"encoding/json"
)

type JsonMarshaler struct{}

func (j JsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j JsonMarshaler) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j JsonMarshaler) String() string {
	return "json"
}
