package alice

import (
	"time"

	"golang.org/x/text/language"
)

// NewAnswer создаёт основу Answer для ответа, на основе входящего Question.
func NewAnswer(question *Question) *Answer {
	return &Answer{
		Version: question.Version,
		Session: Session{
			SessionID: question.Session.SessionID,
			MessageID: question.Session.MessageID,
			UserID:    question.Session.UserID,
		},
	}
}

// NewButtons создаёт новый массив Button.
func NewButtons(buttons ...Button) []Button {
	return buttons
}

// Language декодирует поле Locale в language.Tag.
func (meta Meta) Language() language.Tag {
	return language.Make(meta.Locale)
}

// TimeLocation декодирует поле TimeZone в *time.Location.
func (meta Meta) TimeLocation() (*time.Location, error) {
	return time.LoadLocation(meta.TimeZone)
}

// IsSimpleUtterance проверяет принадлежность запроса к событию голосового ввода.
func (req Request) IsSimpleUtterance() bool {
	return req.Type == TypeSimpleUtterance
}

// IsButtonPressed проверяет принадлежность запроса к событию нажатия на кнопку.
func (req Request) IsButtonPressed() bool {
	return req.Type == TypeButtonPressed
}
