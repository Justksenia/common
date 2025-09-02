package markup

type PhotoOrVideo interface {
	_photoOrVideo()
}

type MediaType string

const (
	MediaTypePhoto     MediaType = "Photo"
	MediaTypeVideo     MediaType = "Video"
	MediaTypeAudio     MediaType = "Audio"
	MediaTypeDocument  MediaType = "Document"
	MediaTypeAnimation MediaType = "Animation"
	MediaTypeVoice     MediaType = "Voice"
	MediaTypeVideoNote MediaType = "VideoNote"
	MediaTypeSticker   MediaType = "Sticker"
	MediaTypeDice      MediaType = "Dice"
)

type (
	File struct {
		Type   MediaType `json:"type"`
		FileID string    `json:"file_id"`
		S3Path string    `json:"s3_file_path"`
	}

	MaskPosition struct {
		Point  string  `json:"point"`
		XShift float32 `json:"x_shift"`
		YShift float32 `json:"y_shift"`
		Scale  float64 `json:"scale"`
	}
)

type Photo struct {
	*File
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Caption string `json:"caption,omitempty"`
}

func (p Photo) _photoOrVideo() {}

type Video struct {
	*File
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`

	Caption   string `json:"caption,omitempty"`
	Thumbnail *Photo `json:"thumb,omitempty"`
	MIMEType  string `json:"mime_type,omitempty"`
	FileName  string `json:"file_name,omitempty"`
}

func (v Video) _photoOrVideo() {}

type Document struct {
	*File

	Thumbnail            *Photo `json:"thumb,omitempty"`
	Caption              string `json:"caption,omitempty"`
	MIME                 string `json:"mime_type"`
	FileName             string `json:"file_name,omitempty"`
	DisableTypeDetection bool   `json:"disable_content_type_detection,omitempty"`
}

type Audio struct {
	*File

	Duration int `json:"duration,omitempty"`

	// (Optional)
	Caption   string `json:"caption,omitempty"`
	Thumbnail *Photo `json:"thumb,omitempty"`
	Title     string `json:"title,omitempty"`
	Performer string `json:"performer,omitempty"`
	MIME      string `json:"mime_type,omitempty"`
	FileName  string `json:"file_name,omitempty"`
}

type Animation struct {
	*File
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`

	Caption   string `json:"caption,omitempty"`
	Thumbnail *Photo `json:"thumb,omitempty"`
	MIMEType  string `json:"mime_type,omitempty"`
	FileName  string `json:"file_name,omitempty"`
}

type Voice struct {
	*File
	Duration int `json:"duration"`

	Caption  string `json:"caption,omitempty"`
	MIMEType string `json:"mime_type,omitempty"`
}

type VideoNote struct {
	*File
	Duration  int    `json:"duration"`
	Thumbnail *Photo `json:"thumb,omitempty"`
	Length    int    `json:"length"`
}

type Sticker struct {
	*File
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	IsAnimated   bool          `json:"is_animated"`
	IsVideo      bool          `json:"is_video"`
	Thumbnail    *Photo        `json:"thumb"`
	Emoji        string        `json:"emoji"`
	SetName      string        `json:"set_name"`
	MaskPosition *MaskPosition `json:"mask_position"`
}

type DiceType string

var (
	Cube = &Dice{Type: "🎲"}
	Dart = &Dice{Type: "🎯"}
	Ball = &Dice{Type: "🏀"}
	Goal = &Dice{Type: "⚽"}
	Slot = &Dice{Type: "🎰"}
	Bowl = &Dice{Type: "🎳"}
)

type Dice struct {
	Type  DiceType `json:"emoji"`
	Value int      `json:"value"`
}
