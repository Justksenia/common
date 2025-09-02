package telegram

import (
	"encoding/json"

	"github.com/Justksenia/common/entities/telegram/markup"
)

type MarkupV1 struct {
	Message    string                `json:"message"`
	Entities   []markup.Entity       `json:"entities,omitempty"`
	Media      []markup.PhotoOrVideo `json:"media,omitempty"`
	Audios     []markup.Audio        `json:"audios,omitempty"`
	Voices     []markup.Voice        `json:"voices,omitempty"`
	VideoNotes []markup.VideoNote    `json:"video_notes,omitempty"`
	Stickers   []markup.Sticker      `json:"stickers,omitempty"`
	Poll       *markup.Poll          `json:"poll,omitempty"`
	Documents  []markup.Document     `json:"documents,omitempty"`
}

func (m *MarkupV1) ToString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (m *MarkupV1) FromString(s string) error {
	return json.Unmarshal([]byte(s), m)
}
