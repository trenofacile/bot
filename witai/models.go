package witai

import "encoding/json"

// Entity witAI entity
type Entity struct {
	Confidence float64         `json:"confidence"`
	Type       string          `json:"type,omitempty"`
	Value      json.RawMessage `json:"value"`
}

// ValueToInt ...
func (e *Entity) ValueToInt() (int, error) {
	output := int(0)

	err := json.Unmarshal(e.Value, &output)
	if err != nil {
		return output, err
	}

	return output, nil
}

// ValueToString ...
func (e *Entity) ValueToString() (string, error) {
	output := ""

	err := json.Unmarshal(e.Value, &output)
	if err != nil {
		return output, err
	}

	return output, err
}

// ValueToFloat64 ...
func (e *Entity) ValueToFloat64() (float64, error) {
	output := float64(0)

	err := json.Unmarshal(e.Value, &output)
	if err != nil {
		return output, err
	}

	return output, err
}

// Meaning wiAI meaning
type Meaning struct {
	MessageID string              `json:"msg_id"`
	InputText string              `json:"_text"`
	Entities  map[string][]Entity `json:"entities"`
}
