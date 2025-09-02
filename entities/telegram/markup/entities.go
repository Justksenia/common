package markup

type EntityType string

const (
	EntityMention       EntityType = "mention"
	EntityTMention      EntityType = "text_mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
	EntitySpoiler       EntityType = "spoiler"
	EntityCustomEmoji   EntityType = "custom_emoji"
	EntityBlockquote    EntityType = "blockquote"
)

type Entities []Entity

type Entity struct {
	Type   EntityType `json:"type"`
	Offset int        `json:"offset"`
	Length int        `json:"length"`
	URL    string     `json:"url,omitempty"`

	// Language (Optional). Only for EntityCodeBlock
	Language    string   `json:"language,omitempty"`
	CustomEmoji *Sticker `json:"custom_emoji_id"`
}
