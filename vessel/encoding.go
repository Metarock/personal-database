package vessel

import "encoding/json"

type DataEncoder interface {
	Encode(Map) ([]byte, error)
}

type DataDecoder interface {
	Decode([]byte, any) error
}

type JSONEncoder struct{}

func (JSONEncoder) Encode(data Map) ([]byte, error) {
	return json.Marshal(data)
}

type JSONDecoder struct{}

func (JSONDecoder) Decode(bytes []byte, value any) error {
	return json.Unmarshal(bytes, &value)
}
