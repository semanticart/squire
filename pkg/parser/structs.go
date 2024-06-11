package parser

import "strings"

// Story represents a story with a title, author, and chapters.
type Story struct {
	Title    string
	Author   string
	Chapters []Chapter
}

func (s *Story) appendChapter(c Chapter) {
	c.Body = strings.TrimSpace(c.Body)

	s.Chapters = append(s.Chapters, c)
}

// Chapter represents a chapter with a title, body, and choices.
type Chapter struct {
	Title              string
	ID                 string
	Body               string
	Choices            []Choice
	StartLine          int
	EndLine            int
	IntentionalDeadEnd bool
}

// Choice represents a choice with text, a chapter ID, and a line number.
type Choice struct {
	Text      string
	ChapterID string
	Line      int
}
