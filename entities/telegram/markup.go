package telegram

import "encoding/json"

type MarkupVersion string

const (
	MarkupVersion100 MarkupVersion = "1.0.0"
)

type Markup struct {
	Version MarkupVersion `json:"version"`
	Payload string        `json:"payload"`
}

func (m *Markup) ToString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (m *Markup) FromString(s string) error {
	return json.Unmarshal([]byte(s), m)
}
