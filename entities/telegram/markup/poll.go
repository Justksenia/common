package markup

type PollType string

const (
	PollAny     PollType = "any"
	PollQuiz    PollType = "quiz"
	PollRegular PollType = "regular"
)

type Poll struct {
	ID          string   `json:"id"`
	Type        PollType `json:"type"`
	Question    string   `json:"question"`
	PollOptions []string `json:"options"`
}
